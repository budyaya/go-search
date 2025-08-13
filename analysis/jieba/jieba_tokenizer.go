package jieba

import (
	"github.com/blevesearch/bleve/v2/analysis"
	"github.com/yanyiwu/gojieba"
)

type JiebaTokenizer struct {
	jieba *gojieba.Jieba
}

// 确保JiebaTokenizer实现bleve的Tokenizer接口
var _ analysis.Tokenizer = &JiebaTokenizer{}

func NewJiebaTokenizer() *JiebaTokenizer {
	// 初始化Jieba（使用默认词典）
	return &JiebaTokenizer{
		jieba: gojieba.NewJieba(),
	}
}

// Tokenize 实现分词逻辑
func (t *JiebaTokenizer) Tokenize(input []byte) analysis.TokenStream {
	// 使用Jieba精确模式分词
	words := t.jieba.Cut(string(input), true)
	tokens := make(analysis.TokenStream, 0, len(words))
	pos := 0

	for _, word := range words {
		term := []byte(word)
		start := pos
		end := pos + len(word)
		tokens = append(tokens, &analysis.Token{
			Term:     term,
			Start:    start,
			End:      end,
			Position: pos + 1,
			Type:     analysis.Ideographic,
		})
		pos = end
	}
	return tokens
}

// Close 释放Jieba资源
func (t *JiebaTokenizer) Close() error {
	if t.jieba != nil {
		t.jieba.Free()
	}
	return nil
}
