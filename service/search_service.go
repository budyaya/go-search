package service

import (
	"fmt"
	"go-search/model"
	"log"
	"os"
	"sync"

	"github.com/blevesearch/bleve/v2"
)

var (
	indexes = make(map[string]bleve.Index)
	mu      sync.RWMutex
)

// 初始化索引 - 先尝试加载已存在索引，不存在则创建新索引
func InitIndex(indexName string) error {
	mu.Lock()
	defer mu.Unlock()

	if _, exists := indexes[indexName]; exists {
		return fmt.Errorf("索引 %s 已存在", indexName)
	}

	// 尝试打开已存在的索引
	index, err := bleve.Open(indexName)
	if err == nil {
		indexes[indexName] = index
		return nil
	}

	// 如果索引不存在，则创建新索引
	if err == bleve.ErrorIndexPathDoesNotExist {
		mapping := bleve.NewIndexMapping()
		index, err = bleve.New(indexName, mapping)
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
	mu.Lock()
	defer mu.Unlock()

	// 读取当前目录下的所有项目
	entries, err := os.ReadDir(".")
	if err != nil {
		return fmt.Errorf("读取目录失败: %v", err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			// 尝试打开目录作为索引
			index, err := bleve.Open(entry.Name())
			if err == nil {
				indexes[entry.Name()] = index
				log.Printf("成功加载索引: %s", entry.Name())
			} else if err != bleve.ErrorIndexPathDoesNotExist {
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

	searchQuery := bleve.NewMatchQuery(query)
	searchRequest := bleve.NewSearchRequest(searchQuery)

	// 设置分页
	from := (page - 1) * size
	searchRequest.From = from
	searchRequest.Size = size

	return index.Search(searchRequest)
}
