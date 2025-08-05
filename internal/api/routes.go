package api

import (
	"github.com/cardchoosen/AlphaArk_Gin/internal/config"
	"github.com/gin-gonic/gin"
)

// SetupRoutes 设置API路由
func SetupRoutes(r *gin.Engine, cfg *config.Config) {
	// API v1 路由组
	v1 := r.Group("/api/v1")
	{
		// 健康检查
		v1.GET("/ping", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "pong",
			})
		})
	}

	// 设置OKX API路由
	SetupOKXRoutes(r, cfg)

	// 设置价格API路由
	SetupPriceRoutes(r, cfg)

	// 设置WebSocket路由
	SetupWebSocketRoutes(r, cfg)

	// 设置账户API路由
	SetupAccountRoutes(r, cfg)
}
