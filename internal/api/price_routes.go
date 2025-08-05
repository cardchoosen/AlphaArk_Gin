package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/cardchoosen/AlphaArk_Gin/internal/config"
	"github.com/cardchoosen/AlphaArk_Gin/internal/service"
	"github.com/cardchoosen/AlphaArk_Gin/internal/utils"
)

// SetupPriceRoutes 设置价格API路由
func SetupPriceRoutes(r *gin.Engine, cfg *config.Config) {
	priceService := service.NewPriceService(&cfg.OKX)

	// 价格API路由组
	price := r.Group("/api/v1/price")
	{
		// 获取指定交易对价格
		price.GET("/:symbol", func(c *gin.Context) {
			GetPrice(c, priceService)
		})
	}
}

// GetPrice 获取价格信息
func GetPrice(c *gin.Context, priceService service.PriceService) {
	symbol := c.Param("symbol")
	
	if symbol == "" {
		utils.BadRequestResponse(c, "交易对符号不能为空")
		return
	}

	priceData, err := priceService.GetPrice(symbol)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "获取价格失败: "+err.Error())
		return
	}

	utils.SuccessResponse(c, priceData, "获取价格成功")
}