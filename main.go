package main

import (
	"go-search/config"
	"go-search/handler"
	"log/slog"

	"github.com/gin-gonic/gin"
)

func main() {
	// 初始化配置
	config.Init()
	slog.Info("Server started on port 8080")

	// 创建Gin路由
	router := gin.Default()

	// 注册API路由
	api := router.Group("/api")
	{
		api.POST("/index", handler.CreateIndexHandler)
		api.POST("/index/stats", handler.GetIndexStatisticsHandler) // 获取索引统计信息
		api.POST("/document", handler.AddDocumentHandler)
		api.POST("/document/stats", handler.GetDocumentStatisticsHandler)
		api.PUT("/document", handler.UpdateDocumentHandler)
		api.DELETE("/document", handler.DeleteDocumentHandler)
		api.POST("/search", handler.SearchHandler) // 修改为POST方法
	}

	// 启动服务器
	slog.Info("Server started on port 8080")
	router.Run(":8080")
}
