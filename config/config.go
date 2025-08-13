package config

import (
	"go-search/service"
	"log"
)

// 初始化配置
func Init() {
	// 默认初始化一个名为"default"的索引
	if err := service.InitIndex("default", nil); err != nil {
		log.Printf("默认索引初始化失败: %v", err)
	}
	// 加载所有已存在的索引
	if err := service.LoadAllIndexes(); err != nil {
		log.Printf("加载索引失败: %v", err)
	}
}
