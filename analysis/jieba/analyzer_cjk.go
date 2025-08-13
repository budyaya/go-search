package jieba

import (
	"github.com/blevesearch/bleve/v2/analysis"
	"github.com/blevesearch/bleve/v2/analysis/lang/cjk"
	"github.com/blevesearch/bleve/v2/analysis/token/lowercase"
	"github.com/blevesearch/bleve/v2/registry"
)

const (
	AnalyzerName  = "jieba"
	TokenizerName = "jieba_tokenizer"
)

// 初始化函数：注册分词器和分析器
func init() {
	// 注册分词器
	registry.RegisterTokenizer(TokenizerName, func(config map[string]interface{}, cache *registry.Cache) (analysis.Tokenizer, error) {
		return NewJiebaTokenizer(), nil
	})

	// 注册分析器
	registry.RegisterAnalyzer(AnalyzerName, func(config map[string]interface{}, cache *registry.Cache) (analysis.Analyzer, error) {
		// 获取分词器
		tokenizer, err := cache.TokenizerNamed(TokenizerName)
		if err != nil {
			return nil, err
		}

		// 获取Bleve内置的CJK宽度过滤器（处理全角半角转换）
		widthFilter, err := cache.TokenFilterNamed(cjk.WidthName)
		if err != nil {
			return nil, err
		}

		// 小写过滤器
		lowerFilter, err := cache.TokenFilterNamed(lowercase.Name)
		if err != nil {
			return nil, err
		}

		// 组合分析器
		return &analysis.DefaultAnalyzer{
			Tokenizer: tokenizer,
			TokenFilters: []analysis.TokenFilter{
				widthFilter,
				lowerFilter,
			},
		}, nil
	})
}
