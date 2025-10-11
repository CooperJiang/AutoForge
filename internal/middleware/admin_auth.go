package middleware

import (
	"net/http"
	"strings"

	"auto-forge/pkg/common"

	"github.com/gin-gonic/gin"
)

// AdminAuth 管理员认证中间件（基于用户JWT Token + Role检查）
func AdminAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取 Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "未提供认证信息"})
			c.Abort()
			return
		}

		// 检查格式：Bearer token
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "认证格式错误"})
			c.Abort()
			return
		}

		token := parts[1]

		// 验证 JWT token 并解析claims
		claims, err := common.ParseToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "无效或过期的token"})
			c.Abort()
			return
		}

		// 检查用户角色：必须是管理员或超级管理员
		if claims.Role != common.UserRoleAdmin && claims.Role != common.UserRoleSuperAdmin {
			c.JSON(http.StatusForbidden, gin.H{"error": "无权访问管理后台"})
			c.Abort()
			return
		}

		// 将用户信息存入上下文
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("role", claims.Role)

		// 继续处理请求
		c.Next()
	}
}
