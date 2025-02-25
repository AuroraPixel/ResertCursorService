package handlers

import (
	"net/http"
	"time"

	"github.com/ResertCursorService/internal/models"
	"github.com/ResertCursorService/internal/service"
	"github.com/gin-gonic/gin"
)

type AppHandler struct {
	service *service.ActivationService
}

func NewAppHandler(service *service.ActivationService) *AppHandler {
	return &AppHandler{service: service}
}

type ActivateRequest struct {
	Code string `json:"code" binding:"required"`
}

type ActivateResponse struct {
	Token     string `json:"token"`
	ExpiresAt string `json:"expires_at"`
}

// 错误码定义
const (
	ErrCodeInvalidCode     = 1001 // 无效的激活码
	ErrCodeExpiredCode     = 1002 // 激活码已过期
	ErrCodeDisabledCode    = 1003 // 激活码已禁用
	ErrCodeMaxAccountsUsed = 1004 // 已达到最大账户数
)

// Activate 处理 app 端的激活请求
func (h *AppHandler) Activate(c *gin.Context) {
	var req ActivateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error_code": ErrCodeInvalidCode,
			"message":    "无效的请求数据",
		})
		return
	}

	code, err := h.service.GetActivationCodeByCode(req.Code)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error_code": ErrCodeInvalidCode,
			"message":    "激活码不存在",
		})
		return
	}

	// 检查激活码状态
	if code.Status != "enabled" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error_code": ErrCodeDisabledCode,
			"message":    "激活码已禁用",
		})
		return
	}

	// 检查是否过期
	if code.ExpiresAt.Before(time.Now()) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error_code": ErrCodeExpiredCode,
			"message":    "激活码已过期",
		})
		return
	}

	// 检查账户数量
	//if len(code.Accounts) >= code.MaxAccounts {
	//	c.JSON(http.StatusBadRequest, gin.H{
	//		"error_code": ErrCodeMaxAccountsUsed,
	//		"message":    "已达到最大账户数",
	//	})
	//	return
	//}

	// 生成 token
	token, err := h.service.GenerateAppToken(code.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error_code": 500,
			"message":    "生成令牌失败",
		})
		return
	}

	c.JSON(http.StatusOK, ActivateResponse{
		Token:     token,
		ExpiresAt: code.ExpiresAt.Format(time.RFC3339),
	})
}

// GetAccount 获取激活码下的所有 Cursor 账号信息
func (h *AppHandler) GetAccount(c *gin.Context) {
	// 从上下文中获取激活码ID（在中间件中已设置）
	codeID, exists := c.Get("codeID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error_code": 401,
			"message":    "未授权访问",
		})
		return
	}

	// 获取激活码信息
	code, err := h.service.GetActivationCode(codeID.(uint))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error_code": 404,
			"message":    "激活码不存在",
		})
		return
	}

	// 检查激活码状态
	if code.Status != "enabled" {
		c.JSON(http.StatusForbidden, gin.H{
			"error_code": ErrCodeDisabledCode,
			"message":    "激活码已禁用",
		})
		return
	}

	// 检查是否过期
	if code.ExpiresAt.Before(time.Now()) {
		c.JSON(http.StatusForbidden, gin.H{
			"error_code": ErrCodeExpiredCode,
			"message":    "激活码已过期",
		})
		return
	}

	// 返回该激活码下的所有账号
	c.JSON(http.StatusOK, gin.H{
		"accounts": code.Accounts,
	})
}

// CreateAccountRequest 创建Cursor账号的请求
type CreateAccountRequest struct {
	Email          string `json:"email" binding:"required,email"`
	EmailPassword  string `json:"emailPassword" binding:"required"`  // 邮箱密码
	CursorPassword string `json:"cursorPassword" binding:"required"` // Cursor密码
	AccessToken    string `json:"accessToken" binding:"required"`
	RefreshToken   string `json:"refreshToken" binding:"required"`
}

// CreateAccount 上传Cursor账号信息
func (h *AppHandler) CreateAccount(c *gin.Context) {
	// 从上下文中获取激活码ID（在中间件中已设置）
	codeID, exists := c.Get("codeID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error_code": 401,
			"message":    "未授权访问",
		})
		return
	}

	// 获取激活码信息
	code, err := h.service.GetActivationCode(codeID.(uint))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error_code": 404,
			"message":    "激活码不存在",
		})
		return
	}

	// 检查激活码状态
	if code.Status != "enabled" {
		c.JSON(http.StatusForbidden, gin.H{
			"error_code": ErrCodeDisabledCode,
			"message":    "激活码已禁用",
		})
		return
	}

	// 检查是否过期
	if code.ExpiresAt.Before(time.Now()) {
		c.JSON(http.StatusForbidden, gin.H{
			"error_code": ErrCodeExpiredCode,
			"message":    "激活码已过期",
		})
		return
	}

	//检查账户数量
	if len(code.Accounts) >= code.MaxAccounts {
		c.JSON(http.StatusBadRequest, gin.H{
			"error_code": ErrCodeMaxAccountsUsed,
			"message":    "已达到最大账户数",
		})
		return
	}

	// 解析请求数据
	var req CreateAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error_code": 400,
			"message":    "无效的请求数据",
		})
		return
	}

	// 创建账号对象
	account := &models.CursorAccount{
		Email:            req.Email,
		EmailPassword:    req.EmailPassword,  // 邮箱密码
		CursorPassword:   req.CursorPassword, // Cursor密码
		AccessToken:      req.AccessToken,
		RefreshToken:     req.RefreshToken,
		ActivationCodeID: code.ID, // 添加激活码ID
	}

	// 添加账号
	if err := h.service.AddAccount(code, account); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error_code": 500,
			"message":    "添加账号失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "账号添加成功",
		"account": account,
	})
}

// CodeInfoResponse 激活码信息响应
type CodeInfoResponse struct {
	Code         string    `json:"code"`
	ExpiresAt    time.Time `json:"expiresAt"`
	MaxAccounts  int       `json:"maxAccounts"`
	UsedAccounts int       `json:"usedAccounts"`
	Status       string    `json:"status"`
}

// GetCodeInfo 获取当前授权码信息
func (h *AppHandler) GetCodeInfo(c *gin.Context) {
	// 从上下文中获取激活码ID（在中间件中已设置）
	codeID, exists := c.Get("codeID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error_code": 401,
			"message":    "未授权访问",
		})
		return
	}

	// 获取激活码信息
	code, err := h.service.GetActivationCode(codeID.(uint))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error_code": 404,
			"message":    "激活码不存在",
		})
		return
	}

	// 检查激活码状态
	if code.Status != "enabled" {
		c.JSON(http.StatusForbidden, gin.H{
			"error_code": ErrCodeDisabledCode,
			"message":    "激活码已禁用",
		})
		return
	}

	// 检查是否过期
	if code.ExpiresAt.Before(time.Now()) {
		c.JSON(http.StatusForbidden, gin.H{
			"error_code": ErrCodeExpiredCode,
			"message":    "激活码已过期",
		})
		return
	}

	// 返回激活码信息
	c.JSON(http.StatusOK, CodeInfoResponse{
		Code:         code.Code,
		ExpiresAt:    code.ExpiresAt,
		MaxAccounts:  code.MaxAccounts,
		UsedAccounts: len(code.Accounts),
		Status:       code.Status,
	})
}
