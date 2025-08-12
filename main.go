package main

import (
	"go-search/config"
	"go-search/handler"

	"github.com/gin-gonic/gin"
)

func main() {
	// 初始化配置
	config.Init()

	// 创建Gin路由
	router := gin.Default()

	// 注册API路由
	api := router.Group("/api")
	{
		api.POST("/index", handler.CreateIndexHandler)
		api.POST("/document", handler.AddDocumentHandler)
		api.POST("/search", handler.SearchHandler) // 修改为POST方法
	}

	// 启动服务器
	router.Run(":8080")
}
