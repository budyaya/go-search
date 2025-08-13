package model

// 定义统计信息结构体
type TermFrequency struct {
	Term      string `json:"term"`
	Frequency uint64 `json:"frequency"`
}

type IndexStatistics struct {
	DocCount   uint64         // 文档数量
	IndexSize  uint64         // 索引大小(字节)
	FieldCount int            // 字段数量
	FieldFreq  map[string]int // 字段频率统计
}
