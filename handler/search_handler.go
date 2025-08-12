package handler

import (
	"go-search/model"
	"go-search/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

// 创建索引请求体
type CreateIndexRequest struct {
	IndexName string `json:"index_name" binding:"required"`
}

// 添加文档请求体
type AddDocumentRequest struct {
	IndexName string                 `json:"index_name" binding:"required"`
	ID        string                 `json:"id" binding:"required"`
	Fields    map[string]interface{} `json:"fields" binding:"required"`
}

// 搜索请求体 (新增)
type SearchRequest struct {
	IndexName string `json:"index_name" binding:"required"`
	Query     string `json:"query" binding:"required"`
	Page      int    `json:"page,omitempty"` // 可选分页参数
	Size      int    `json:"size,omitempty"` // 可选每页数量
}

// 创建索引
func CreateIndexHandler(c *gin.Context) {
	var req CreateIndexRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := service.InitIndex(req.IndexName); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "索引创建成功"})
}

// 添加文档
func AddDocumentHandler(c *gin.Context) {
	var req AddDocumentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	doc := model.Document{
		ID:     req.ID,
		Fields: req.Fields,
	}

	if err := service.AddDocument(req.IndexName, doc); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "文档添加成功"})
}

// 更新文档请求体
type UpdateDocumentRequest struct {
	IndexName string                 `json:"index_name" binding:"required"`
	ID        string                 `json:"id" binding:"required"`
	Fields    map[string]interface{} `json:"fields" binding:"required"`
}

// 更新文档
func UpdateDocumentHandler(c *gin.Context) {
	var req UpdateDocumentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	doc := model.Document{
		ID:     req.ID,
		Fields: req.Fields,
	}

	if err := service.UpdateDocument(req.IndexName, doc); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "文档更新成功"})
}

// 删除文档请求体
type DeleteDocumentRequest struct {
	IndexName string `json:"index_name" binding:"required"`
	ID        string `json:"id" binding:"required"`
}

// 删除文档
func DeleteDocumentHandler(c *gin.Context) {
	var req DeleteDocumentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := service.DeleteDocument(req.IndexName, req.ID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "文档删除成功"})
}

// 搜索文档 (修改为JSON请求)
func SearchHandler(c *gin.Context) {
	var req SearchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 设置默认分页参数
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Size <= 0 {
		req.Size = 10
	}

	result, err := service.Search(req.IndexName, req.Query, req.Page, req.Size)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"total": result.Total,
		"page":  req.Page,
		"size":  req.Size,
		"hits":  result.Hits,
	})
}
