package handlers

import (
	"net/http"
	"strconv"

	"github.com/ResertCursorService/internal/service"
	"github.com/gin-gonic/gin"
)

type ActivationHandler struct {
	service *service.ActivationService
}

func NewActivationHandler(service *service.ActivationService) *ActivationHandler {
	return &ActivationHandler{service: service}
}

type CreateActivationCodeRequest struct {
	Duration    int `json:"duration" binding:"required,min=1"`
	MaxAccounts int `json:"maxAccounts" binding:"required,min=1,max=100"`
}

func (h *ActivationHandler) CreateActivationCode(c *gin.Context) {
	var req CreateActivationCodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	code, err := h.service.CreateActivationCode(req.Duration, req.MaxAccounts)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建激活码失败"})
		return
	}

	c.JSON(http.StatusOK, code)
}

func (h *ActivationHandler) ListActivationCodes(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	result, err := h.service.ListActivationCodes(page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取激活码列表失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"items":      result.Items,
		"total":      result.Total,
		"page":       result.Page,
		"pageSize":   result.PageSize,
		"totalPages": result.TotalPages,
	})
}

func (h *ActivationHandler) GetActivationCode(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的ID"})
		return
	}

	code, err := h.service.GetActivationCode(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "激活码不存在"})
		return
	}

	c.JSON(http.StatusOK, code)
}

type UpdateStatusRequest struct {
	Status string `json:"status" binding:"required,oneof=enabled disabled"`
}

func (h *ActivationHandler) UpdateActivationCodeStatus(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的ID"})
		return
	}

	var req UpdateStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.UpdateActivationCodeStatus(uint(id), req.Status); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新状态失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "状态更新成功"})
}
