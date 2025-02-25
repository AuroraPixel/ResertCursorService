package middleware

import (
	"net/http"
	"strings"

	"github.com/ResertCursorService/pkg/jwt"
	"github.com/gin-gonic/gin"
)

// AdminAuthMiddleware 管理员认证中间件
func AdminAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "未提供认证信息"})
			c.Abort()
			return
		}

		tokenString := strings.Replace(authHeader, "Bearer ", "", 1)
		claims, err := jwt.ParseToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "无效的管理员认证信息"})
			c.Abort()
			return
		}

		c.Set("adminID", claims.AdminID)
		c.Next()
	}
}

// AppAuthMiddleware app认证中间件
func AppAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "未提供认证信息"})
			c.Abort()
			return
		}

		tokenString := strings.Replace(authHeader, "Bearer ", "", 1)
		claims, err := jwt.ParseAppToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "无效的app认证信息"})
			c.Abort()
			return
		}

		c.Set("codeID", claims.CodeID)
		c.Next()
	}
}
