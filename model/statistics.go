package model

// 定义统计信息结构体
type TermFrequency struct {
	Term      string `json:"term"`
	Frequency uint64 `json:"frequency"`
}

type IndexStatistics struct {
	DocCount           uint64              // 文档数量
	IndexSize          uint64              // 索引大小(字节)
	FieldCount         int                 // 字段数量
	FieldFreq          map[string]int      // 字段频率统计
	RangeDistributions []RangeDistribution `json:"range_distributions,omitempty"` // 新增范围分布字段
}

// 范围分布统计结果
type RangeDistribution struct {
	FieldName string         `json:"field_name"`
	Ranges    map[string]int `json:"ranges"` // key: 范围描述(如"0-100"), value: 数量
	Min       float64        `json:"min"`
	Max       float64        `json:"max"`
	Avg       float64        `json:"avg"`
	Count     int            `json:"count"` // 非空值总数
}
