package handler

import (
	"go-search/model"
	"go-search/service"
	"math"
	"net/http"

	"github.com/gin-gonic/gin"
)

// 创建索引请求体
type CreateIndexRequest struct {
	IndexName string            `json:"index_name" binding:"required"`
	Fields    map[string]string `json:"fields"` // 字段分词器配置
}

// 添加文档请求体
type AddDocumentRequest struct {
	IndexName string                 `json:"index_name" binding:"required"`
	ID        string                 `json:"id" binding:"required"`
	Fields    map[string]interface{} `json:"fields" binding:"required"`
}

// 搜索请求体 (新增)
type SearchRequest struct {
	IndexName string  `json:"index_name" binding:"required"`
	Type      int     `json:"type" binding:"required"` // 1: 普通搜索, 2: 范围查询
	Query     string  `json:"query" binding:"required_if=Type 1"`
	Field     string  `json:"field" binding:"required_if=Type 2"`
	Start     float64 `json:"start" binding:"required_if=Type 2"`
	End       float64 `json:"end" binding:"required_if=Type 2"`
	Page      int     `json:"page,omitempty"` // 可选分页参数
	Size      int     `json:"size,omitempty"` // 可选每页数量
}

// 创建索引
func CreateIndexHandler(c *gin.Context) {
	var req CreateIndexRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := service.InitIndex(req.IndexName, req.Fields); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "索引创建成功"})
}

// 获取索引统计信息请求体
type GetIndexStatisticsHandlerRequest struct {
	IndexName string `json:"index_name" binding:"required"`
}

// 获取索引统计信息
func GetIndexStatisticsHandler(c *gin.Context) {
	var req GetIndexStatisticsHandlerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	stat, err := service.GetIndexStatistics(req.IndexName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, stat)
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

// 获取文档统计信息请求体
type GetDocumentStatisticsHandlerRequest struct {
	IndexName string `json:"index_name" binding:"required"`
}

// 获取文档统计信息
func GetDocumentStatisticsHandler(c *gin.Context) {
	var req GetDocumentStatisticsHandlerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	stat, err := service.GetTermFrequencyRanking(req.IndexName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, stat)
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

	if req.Field == "price" && req.Start != req.End {
		result, err := service.RangeSearch(req.IndexName, req.Field, req.Start, req.End, req.Page, req.Size)
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
		return
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

// 获取统计指定数字字段的范围分布请求体
type GetNumberFieldRangeDistributionHandlerRequest struct {
	IndexName string       `json:"index_name" binding:"required"`
	FieldName string       `json:"field_name" binding:"required"`
	Ranges    [][2]float64 `json:"ranges" binding:"required"`
}

// 获取统计指定数字字段的范围分布
func GetNumberFieldRangeDistributionHandler(c *gin.Context) {
	var req GetNumberFieldRangeDistributionHandlerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ranges := req.Ranges
	// 补上一个无穷大的范围
	ranges = append(ranges, [2]float64{ranges[len(ranges)-1][1], math.Inf(1)})
	// 补上一个无穷小的范围
	ranges = append([][2]float64{{math.Inf(-1), ranges[0][0]}}, ranges...)
	dist, err := service.GetNumberFieldRangeDistribution("products", "price", ranges)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, dist)
}
