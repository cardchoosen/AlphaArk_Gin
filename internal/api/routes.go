package api

import (
	"github.com/gin-gonic/gin"
)

// SetupRoutes 设置API路由
func SetupRoutes(r *gin.Engine) {
	// API v1 路由组
	v1 := r.Group("/api/v1")
	{
		// 用户相关路由
		users := v1.Group("/users")
		{
			users.GET("", GetUsers)
			users.GET("/:id", GetUser)
			users.POST("", CreateUser)
			users.PUT("/:id", UpdateUser)
			users.DELETE("/:id", DeleteUser)
		}

		// 其他API路由可以在这里添加
		v1.GET("/ping", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "pong",
			})
		})
	}
}

// GetUsers 获取用户列表
func GetUsers(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "Get users list",
		"data":    []string{},
	})
}

// GetUser 获取单个用户
func GetUser(c *gin.Context) {
	id := c.Param("id")
	c.JSON(200, gin.H{
		"message": "Get user",
		"id":      id,
	})
}

// CreateUser 创建用户
func CreateUser(c *gin.Context) {
	c.JSON(201, gin.H{
		"message": "User created",
	})
}

// UpdateUser 更新用户
func UpdateUser(c *gin.Context) {
	id := c.Param("id")
	c.JSON(200, gin.H{
		"message": "User updated",
		"id":      id,
	})
}

// DeleteUser 删除用户
func DeleteUser(c *gin.Context) {
	id := c.Param("id")
	c.JSON(200, gin.H{
		"message": "User deleted",
		"id":      id,
	})
} 