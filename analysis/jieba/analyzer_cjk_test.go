package jieba

import (
	"testing"
)

func TestAnalyzerCjk(t *testing.T) {
	tokens := NewJiebaTokenizer().Tokenize([]byte(`星光色`))
	t.Log(tokens)
}
