package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response 标准响应结构
type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// SuccessResponse 成功响应
func SuccessResponse(c *gin.Context, data interface{}, message string) {
	c.JSON(http.StatusOK, Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// ErrorResponse 错误响应
func ErrorResponse(c *gin.Context, statusCode int, message string) {
	c.JSON(statusCode, Response{
		Success: false,
		Error:   message,
	})
}

// BadRequestResponse 400错误响应
func BadRequestResponse(c *gin.Context, message string) {
	ErrorResponse(c, http.StatusBadRequest, message)
}

// NotFoundResponse 404错误响应
func NotFoundResponse(c *gin.Context, message string) {
	ErrorResponse(c, http.StatusNotFound, message)
}

// InternalServerErrorResponse 500错误响应
func InternalServerErrorResponse(c *gin.Context, message string) {
	ErrorResponse(c, http.StatusInternalServerError, message)
}

// UnauthorizedResponse 401错误响应
func UnauthorizedResponse(c *gin.Context, message string) {
	ErrorResponse(c, http.StatusUnauthorized, message)
} 