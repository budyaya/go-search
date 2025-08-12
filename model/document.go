package model

// Document 定义搜索文档结构
type Document struct {
	ID     string                 `json:"id"`
	Fields map[string]interface{} `json:"fields"`
}
