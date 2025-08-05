package api

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/cardchoosen/AlphaArk_Gin/internal/config"
	"github.com/cardchoosen/AlphaArk_Gin/internal/utils"
)

// SetupOKXRoutes 设置OKX API路由
func SetupOKXRoutes(r *gin.Engine, cfg *config.Config) {
	okxClient := NewOKXClient(&cfg.OKX)

	// OKX API路由组
	okx := r.Group("/api/v1/okx")
	{
		// 获取交易对信息
		okx.GET("/instruments", func(c *gin.Context) {
			GetInstruments(c, okxClient)
		})

		// 获取特定类型的交易对
		okx.GET("/instruments/:type", func(c *gin.Context) {
			GetInstrumentsByType(c, okxClient)
		})

		// 获取API配置信息（仅显示非敏感信息）
		okx.GET("/config", func(c *gin.Context) {
			GetOKXConfig(c, cfg)
		})
	}
}

// GetInstruments 获取所有交易对信息
func GetInstruments(c *gin.Context, client *OKXClient) {
	// 默认获取SPOT类型的交易对
	instType := c.DefaultQuery("instType", "SPOT")

	result, err := client.GetInstruments(instType)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "获取交易对信息失败: "+err.Error())
		return
	}

	utils.SuccessResponse(c, result, "获取交易对信息成功")
}

// GetInstrumentsByType 根据类型获取交易对信息
func GetInstrumentsByType(c *gin.Context, client *OKXClient) {
	instType := c.Param("type")

	// 验证交易对类型
	validTypes := map[string]bool{
		"SPOT":    true,
		"MARGIN":  true,
		"SWAP":    true,
		"FUTURES": true,
		"OPTION":  true,
	}

	if !validTypes[instType] {
		utils.BadRequestResponse(c, "无效的交易对类型，支持的类型: SPOT, MARGIN, SWAP, FUTURES, OPTION")
		return
	}

	result, err := client.GetInstruments(instType)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "获取交易对信息失败: "+err.Error())
		return
	}

	utils.SuccessResponse(c, result, "获取"+instType+"交易对信息成功")
}

// GetOKXConfig 获取OKX配置信息（仅显示非敏感信息）
func GetOKXConfig(c *gin.Context, cfg *config.Config) {
	// 直接检查环境变量
	apiKey := os.Getenv("OKX_API_KEY")
	secretKey := os.Getenv("OKX_SECRET_KEY")
	passphrase := os.Getenv("OKX_PASSPHRASE")

	// 只返回非敏感信息
	safeConfig := map[string]interface{}{
		"remark":        cfg.OKX.Remark,
		"permissions":   cfg.OKX.Permissions,
		"baseUrl":       cfg.OKX.BaseURL,
		"isTest":        cfg.OKX.IsTest,
		"hasApiKey":     apiKey != "",
		"hasSecretKey":  secretKey != "",
		"hasPassphrase": passphrase != "",
	}

	utils.SuccessResponse(c, safeConfig, "获取OKX配置信息成功")
}
