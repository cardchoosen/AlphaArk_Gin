package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yourname/my-gin-project/internal/api"
	"github.com/yourname/my-gin-project/internal/config"
	"github.com/yourname/my-gin-project/internal/middleware"
)

func main() {
	// 加载配置
	cfg := config.Load()

	// 设置Gin模式
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// 创建Gin引擎
	r := gin.Default()

	// 添加中间件
	r.Use(middleware.CORS())
	r.Use(middleware.Logger())
	r.Use(middleware.Recovery())

	// 设置静态文件路由
	r.Static("/static", "./web/static")
	r.LoadHTMLGlob("web/templates/*")

	// 设置API路由
	api.SetupRoutes(r)

	// 根路由
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"title": "AlphaArk Gin Project",
		})
	})

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"message": "Server is running",
		})
	})

	// 启动服务器
	log.Printf("Server starting on port %s", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
