package service

import (
	"fmt"
	"go-search/analysis/jieba"
	"go-search/model"
	"log"
	"os"
	"regexp"
	"sort"
	"sync"

	"github.com/blevesearch/bleve/v2"
	"github.com/blevesearch/bleve/v2/mapping"
)

var (
	indexes        = make(map[string]bleve.Index)
	mu             sync.RWMutex
	indexNameRegex = regexp.MustCompile(`^[a-zA-Z_.]+$`)
)

// 验证索引名称是否合法
func IsValidIndexName(name string) bool {
	// 验证索引名称只能包含字母、下划线和点
	return indexNameRegex.MatchString(name)
}

// 初始化索引 - 支持字段分词器配置
func InitIndex(indexName string, fields map[string]string) error {

	// 验证索引名称是否合法
	if !IsValidIndexName(indexName) {
		return fmt.Errorf("索引名称不合法")
	}

	mu.Lock()
	defer mu.Unlock()

	if _, exists := indexes[indexName]; exists {
		return fmt.Errorf("索引 %s 已存在", indexName)
	}

	// 尝试打开已存在的索引
	index, err := bleve.Open("./data/" + indexName)
	if err == nil {
		indexes[indexName] = index
		return nil
	}

	// 如果索引不存在，则创建新索引
	if err == bleve.ErrorIndexPathDoesNotExist {
		indexMapping := bleve.NewIndexMapping()

		// 配置字段分词器
		for fieldName, analyzer := range fields {
			var fieldMapping *mapping.FieldMapping

			// 根据配置设置分析器
			switch analyzer {
			case "jieba":
				fieldMapping = bleve.NewTextFieldMapping()
				fieldMapping.Analyzer = jieba.AnalyzerName
			case "keyword":
				fieldMapping = bleve.NewKeywordFieldMapping()
			case "number":
				fieldMapping = bleve.NewNumericFieldMapping()
				log.Printf("number field: %s", fieldName)
			default:
				fieldMapping = bleve.NewTextFieldMapping()
			}

			indexMapping.DefaultMapping.AddFieldMappingsAt(fieldName, fieldMapping)
		}

		index, err = bleve.New("./data/"+indexName, indexMapping)
		if err != nil {
			return fmt.Errorf("创建索引失败: %v", err)
		}
		indexes[indexName] = index
		return nil
	}

	return fmt.Errorf("打开索引失败: %v", err)
}

// 加载所有已存在的索引
func LoadAllIndexes() error {

	// 读取当前目录下的所有项目
	entries, err := os.ReadDir("./data")
	if err != nil {
		return fmt.Errorf("读取目录失败: %v", err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			// 尝试打开目录作为索引
			err = InitIndex(entry.Name(), nil)
			if err == nil {
				log.Printf("成功加载索引: %s", entry.Name())
			} else {
				// 仅记录非不存在错误的警告
				log.Printf("警告: 无法打开目录 %s 作为索引: %v", entry.Name(), err)
			}
		}
	}

	return nil
}

// 添加文档到指定索引
func AddDocument(indexName string, doc model.Document) error {
	mu.RLock()
	defer mu.RUnlock()

	index, exists := indexes[indexName]
	if !exists {
		return fmt.Errorf("索引 %s 不存在", indexName)
	}

	return index.Index(doc.ID, doc.Fields)
}

// 更新文档
func UpdateDocument(indexName string, doc model.Document) error {
	mu.RLock()
	defer mu.RUnlock()

	index, exists := indexes[indexName]
	if !exists {
		return fmt.Errorf("索引 %s 不存在", indexName)
	}

	// 使用bleve的Index方法实现更新（已存在的ID会被覆盖）
	return index.Index(doc.ID, doc.Fields)
}

// 删除文档
func DeleteDocument(indexName string, docID string) error {
	mu.RLock()
	defer mu.RUnlock()

	index, exists := indexes[indexName]
	if !exists {
		return fmt.Errorf("索引 %s 不存在", indexName)
	}

	return index.Delete(docID)
}

// 搜索文档 (增加分页参数)
func Search(indexName string, query string, page, size int) (*bleve.SearchResult, error) {
	mu.RLock()
	defer mu.RUnlock()

	index, exists := indexes[indexName]
	if !exists {
		return nil, fmt.Errorf("索引 %s 不存在", indexName)
	}

	searchQuery := bleve.NewQueryStringQuery(query) // NewMatchQuery
	searchRequest := bleve.NewSearchRequest(searchQuery)

	// 返回所有字段
	searchRequest.Fields = []string{"*"}

	// 设置分页
	from := (page - 1) * size
	searchRequest.From = from
	searchRequest.Size = size

	return index.Search(searchRequest)
}

// 获取索引统计信息
func GetIndexStatistics(indexName string) (*model.IndexStatistics, error) {
	mu.RLock()
	defer mu.RUnlock()

	index, exists := indexes[indexName]
	if !exists {
		return nil, fmt.Errorf("索引 %s 不存在", indexName)
	}

	stats := &model.IndexStatistics{
		FieldFreq: make(map[string]int),
	}

	var err error
	// 1. 获取文档数量
	stats.DocCount, err = index.DocCount()
	if err != nil {
		return nil, fmt.Errorf("获取文档数量失败: %v", err)
	}

	// 2. 获取索引大小
	statsMap := index.StatsMap()
	if indexStats, ok := statsMap["index"].(map[string]interface{}); ok {
		if size, ok := indexStats["size_in_bytes"].(uint64); ok {
			stats.IndexSize = size
		}
	}

	// 3. 获取字段数量及频率统计
	fields, err := index.Fields()
	if err != nil {
		return nil, fmt.Errorf("获取字段列表失败: %v", err)
	}
	stats.FieldCount = len(fields)

	// 4. 统计每个字段的频率
	for _, field := range fields {
		dict, err := index.FieldDict(field)
		if err != nil {
			return nil, fmt.Errorf("获取字段词典失败: %v", err)
		}
		defer dict.Close()

		// 累加字段基数(不同词条数量)
		stats.FieldFreq[field] = dict.Cardinality()
	}

	return stats, nil
}

// 返回按频率排序的词条列表（降序）
func GetTermFrequencyRanking(indexName string) ([]model.TermFrequency, error) {
	mu.RLock()
	defer mu.RUnlock()

	idx, exists := indexes[indexName]
	if !exists {
		return nil, fmt.Errorf("索引 %s 不存在", indexName)
	}

	// 获取所有字段
	fields, err := idx.Fields()
	if err != nil {
		return nil, err
	}

	// 收集所有词条频率
	termFreq := make(map[string]uint64)
	for _, field := range fields {
		if field == "_all" {
			continue
		}
		log.Printf("--------------field: %s", field)
		dict, err := idx.FieldDict(field)
		if err != nil {
			return nil, err
		}
		defer dict.Close()

		// 遍历词条
		for {
			term, err := dict.Next()
			if err != nil {
				return nil, err
			}
			if term == nil {
				break
			}
			log.Printf("词条: %s, 频率: %d", term.Term, term.Count)
			termFreq[term.Term] += term.Count
		}
	}

	// 转换为切片并排序
	rankings := make([]model.TermFrequency, 0, len(termFreq))

	for term, count := range termFreq {
		rankings = append(rankings, model.TermFrequency{Term: term, Frequency: count})
	}

	// 按频率降序排序
	sort.Slice(rankings, func(i, j int) bool {
		return rankings[i].Frequency > rankings[j].Frequency
	})

	return rankings, nil
}
