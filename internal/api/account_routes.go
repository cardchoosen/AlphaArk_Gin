package api

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/cardchoosen/AlphaArk_Gin/internal/config"
	"github.com/cardchoosen/AlphaArk_Gin/internal/models"
	"github.com/cardchoosen/AlphaArk_Gin/internal/service"
	"github.com/cardchoosen/AlphaArk_Gin/internal/utils"
)

// SetupAccountRoutes 设置账户API路由
func SetupAccountRoutes(r *gin.Engine, cfg *config.Config) {
	accountService := service.NewAccountService(&cfg.OKX)

	// 账户API路由组
	account := r.Group("/api/v1/account")
	{
		// 获取账户余额
		account.GET("/balance", func(c *gin.Context) {
			GetAccountBalance(c, accountService)
		})

		// 获取账户余额（指定币种）
		account.GET("/balance/:currency", func(c *gin.Context) {
			GetAccountBalanceWithCurrency(c, accountService)
		})

		// 获取盈亏信息
		account.GET("/profit-loss", func(c *gin.Context) {
			GetProfitLoss(c, accountService)
		})

		// 获取账户汇总
		account.GET("/summary", func(c *gin.Context) {
			GetAccountSummary(c, accountService)
		})

		// 获取账户汇总（指定币种）
		account.GET("/summary/:currency", func(c *gin.Context) {
			GetAccountSummaryWithCurrency(c, accountService)
		})

		// 设置默认币种
		account.POST("/currency", func(c *gin.Context) {
			SetDefaultCurrency(c, accountService)
		})

		// 获取默认币种
		account.GET("/currency", func(c *gin.Context) {
			GetDefaultCurrency(c, accountService)
		})

		// 获取支持的币种列表
		account.GET("/currencies", func(c *gin.Context) {
			GetSupportedCurrencies(c)
		})

		// 获取汇率信息
		account.GET("/exchange-rates", func(c *gin.Context) {
			GetExchangeRates(c, accountService)
		})

		// 获取当前持仓信息
		account.GET("/positions", func(c *gin.Context) {
			GetPositions(c, accountService)
		})

		// 获取历史持仓信息
		account.GET("/positions-history", func(c *gin.Context) {
			GetPositionsHistory(c, accountService)
		})

		// 获取持仓完整历史（基于当前持仓的更新时间）
		account.GET("/positions/:posId/history", func(c *gin.Context) {
			GetPositionHistoryByPosId(c, accountService)
		})
	}
}

// GetAccountBalance 获取账户余额（使用默认币种）
func GetAccountBalance(c *gin.Context, accountService service.AccountService) {
	currency := accountService.GetDefaultCurrency()
	
	balance, err := accountService.GetAccountBalance(currency)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "获取账户余额失败: "+err.Error())
		return
	}

	utils.SuccessResponse(c, balance, "获取账户余额成功")
}

// GetAccountBalanceWithCurrency 获取账户余额（指定币种）
func GetAccountBalanceWithCurrency(c *gin.Context, accountService service.AccountService) {
	currencyStr := strings.ToUpper(c.Param("currency"))
	currency := models.Currency(currencyStr)

	// 验证币种
	if !isValidCurrency(currency) {
		utils.BadRequestResponse(c, "不支持的币种: "+currencyStr)
		return
	}

	balance, err := accountService.GetAccountBalance(currency)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "获取账户余额失败: "+err.Error())
		return
	}

	utils.SuccessResponse(c, balance, "获取账户余额成功")
}

// GetProfitLoss 获取盈亏信息
func GetProfitLoss(c *gin.Context, accountService service.AccountService) {
	// 获取币种参数
	currencyStr := strings.ToUpper(c.DefaultQuery("currency", string(accountService.GetDefaultCurrency())))
	currency := models.Currency(currencyStr)

	if !isValidCurrency(currency) {
		utils.BadRequestResponse(c, "不支持的币种: "+currencyStr)
		return
	}

	// 获取时间周期参数
	periodsStr := c.DefaultQuery("periods", "1d,1w,1m,6m")
	periods := parsePeriods(periodsStr)

	profitLoss, err := accountService.GetProfitLoss(currency, periods)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "获取盈亏信息失败: "+err.Error())
		return
	}

	utils.SuccessResponse(c, profitLoss, "获取盈亏信息成功")
}

// GetAccountSummary 获取账户汇总（使用默认币种）
func GetAccountSummary(c *gin.Context, accountService service.AccountService) {
	currency := accountService.GetDefaultCurrency()
	
	summary, err := accountService.GetAccountSummary(currency)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "获取账户汇总失败: "+err.Error())
		return
	}

	utils.SuccessResponse(c, summary, "获取账户汇总成功")
}

// GetAccountSummaryWithCurrency 获取账户汇总（指定币种）
func GetAccountSummaryWithCurrency(c *gin.Context, accountService service.AccountService) {
	currencyStr := strings.ToUpper(c.Param("currency"))
	currency := models.Currency(currencyStr)

	if !isValidCurrency(currency) {
		utils.BadRequestResponse(c, "不支持的币种: "+currencyStr)
		return
	}

	summary, err := accountService.GetAccountSummary(currency)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "获取账户汇总失败: "+err.Error())
		return
	}

	utils.SuccessResponse(c, summary, "获取账户汇总成功")
}

// SetDefaultCurrency 设置默认币种
func SetDefaultCurrency(c *gin.Context, accountService service.AccountService) {
	var req struct {
		Currency string `json:"currency" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(c, "请求参数错误: "+err.Error())
		return
	}

	currency := models.Currency(strings.ToUpper(req.Currency))
	if !isValidCurrency(currency) {
		utils.BadRequestResponse(c, "不支持的币种: "+req.Currency)
		return
	}

	if err := accountService.SetDefaultCurrency(currency); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "设置默认币种失败: "+err.Error())
		return
	}

	utils.SuccessResponse(c, gin.H{
		"currency": currency,
		"symbol":   currency.GetCurrencySymbol(),
	}, "设置默认币种成功")
}

// GetDefaultCurrency 获取默认币种
func GetDefaultCurrency(c *gin.Context, accountService service.AccountService) {
	currency := accountService.GetDefaultCurrency()
	
	utils.SuccessResponse(c, gin.H{
		"currency": currency,
		"symbol":   currency.GetCurrencySymbol(),
	}, "获取默认币种成功")
}

// GetSupportedCurrencies 获取支持的币种列表
func GetSupportedCurrencies(c *gin.Context) {
	currencies := models.SupportedCurrencies()
	
	var currencyList []gin.H
	for _, currency := range currencies {
		currencyList = append(currencyList, gin.H{
			"currency": currency,
			"symbol":   currency.GetCurrencySymbol(),
			"name":     getCurrencyName(currency),
		})
	}

	utils.SuccessResponse(c, currencyList, "获取支持币种列表成功")
}

// GetExchangeRates 获取汇率信息
func GetExchangeRates(c *gin.Context, accountService service.AccountService) {
	rates, err := accountService.GetExchangeRates()
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "获取汇率信息失败: "+err.Error())
		return
	}

	utils.SuccessResponse(c, rates, "获取汇率信息成功")
}

// 辅助函数

// isValidCurrency 验证币种是否有效
func isValidCurrency(currency models.Currency) bool {
	supportedCurrencies := models.SupportedCurrencies()
	for _, supported := range supportedCurrencies {
		if currency == supported {
			return true
		}
	}
	return false
}

// parsePeriods 解析时间周期参数
func parsePeriods(periodsStr string) []models.TimePeriod {
	periodStrs := strings.Split(periodsStr, ",")
	var periods []models.TimePeriod

	for _, periodStr := range periodStrs {
		periodStr = strings.TrimSpace(periodStr)
		period := models.TimePeriod(periodStr)
		
		// 验证时间周期是否有效
		supportedPeriods := models.SupportedPeriods()
		for _, supported := range supportedPeriods {
			if period == supported {
				periods = append(periods, period)
				break
			}
		}
	}

	// 如果没有有效的时间周期，返回默认值
	if len(periods) == 0 {
		periods = models.SupportedPeriods()
	}

	return periods
}

// GetPositions 获取当前持仓信息
func GetPositions(c *gin.Context, accountService service.AccountService) {
	// 获取查询参数
	var req models.PositionsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		utils.BadRequestResponse(c, "请求参数错误: "+err.Error())
		return
	}

	// 获取币种参数
	currencyStr := strings.ToUpper(c.DefaultQuery("currency", string(accountService.GetDefaultCurrency())))
	currency := models.Currency(currencyStr)

	if !isValidCurrency(currency) {
		utils.BadRequestResponse(c, "不支持的币种: "+currencyStr)
		return
	}

	// 获取当前持仓信息
	response, err := accountService.GetPositions(&req, currency)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "获取当前持仓信息失败: "+err.Error())
		return
	}

	utils.SuccessResponse(c, response, "获取当前持仓信息成功")
}

// GetPositionsHistory 获取历史持仓信息
func GetPositionsHistory(c *gin.Context, accountService service.AccountService) {
	// 获取查询参数
	var req models.PositionsHistoryRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		utils.BadRequestResponse(c, "请求参数错误: "+err.Error())
		return
	}

	// 获取币种参数
	currencyStr := strings.ToUpper(c.DefaultQuery("currency", string(accountService.GetDefaultCurrency())))
	currency := models.Currency(currencyStr)

	if !isValidCurrency(currency) {
		utils.BadRequestResponse(c, "不支持的币种: "+currencyStr)
		return
	}

	// 如果设置了fromCurrentPositions参数，自动获取当前持仓的时间戳作为参考
	if c.Query("fromCurrentPositions") == "true" {
		// 获取当前持仓信息
		currentPosReq := &models.PositionsRequest{
			InstType: req.InstType,
			InstId:   req.InstId,
		}
		
		currentPositions, err := accountService.GetPositions(currentPosReq, currency)
		if err == nil && len(currentPositions.Positions) > 0 {
			// 使用最新的持仓更新时间作为before参数
			if req.Before == "" {
				req.Before = currentPositions.Positions[0].UTime
			}
		}
	}

	// 获取历史持仓信息
	response, err := accountService.GetPositionsHistory(&req, currency)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "获取历史持仓信息失败: "+err.Error())
		return
	}

	utils.SuccessResponse(c, response, "获取历史持仓信息成功")
}

// GetPositionHistoryByPosId 根据持仓ID获取完整历史
func GetPositionHistoryByPosId(c *gin.Context, accountService service.AccountService) {
	posId := c.Param("posId")
	if posId == "" {
		utils.BadRequestResponse(c, "持仓ID不能为空")
		return
	}

	// 获取币种参数
	currencyStr := strings.ToUpper(c.DefaultQuery("currency", string(accountService.GetDefaultCurrency())))
	currency := models.Currency(currencyStr)

	if !isValidCurrency(currency) {
		utils.BadRequestResponse(c, "不支持的币种: "+currencyStr)
		return
	}

	// 首先获取当前持仓信息以获取uTime
	currentPosReq := &models.PositionsRequest{
		PosId: posId,
	}
	
	currentPositions, err := accountService.GetPositions(currentPosReq, currency)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "获取当前持仓信息失败: "+err.Error())
		return
	}

	// 构建历史持仓查询请求
	historyReq := &models.PositionsHistoryRequest{
		PosId: posId,
		Limit: c.DefaultQuery("limit", "100"),
	}

	// 如果找到了当前持仓，使用其uTime作为before参数来查询历史
	var currentUTime string
	if len(currentPositions.Positions) > 0 {
		currentUTime = currentPositions.Positions[0].UTime
		// 可以选择性地设置before参数，查询该时间点之前的历史
		if c.Query("includeCurrent") != "true" {
			historyReq.Before = currentUTime
		}
	}

	// 获取历史持仓信息
	historyResponse, err := accountService.GetPositionsHistory(historyReq, currency)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "获取历史持仓信息失败: "+err.Error())
		return
	}

	// 构建完整响应
	response := gin.H{
		"posId":           posId,
		"currentPosition": nil,
		"history":         historyResponse.Positions,
		"hasMore":         historyResponse.HasMore,
		"currency":        currency,
	}

	// 如果有当前持仓，添加到响应中
	if len(currentPositions.Positions) > 0 {
		response["currentPosition"] = currentPositions.Positions[0]
		response["currentUTime"] = currentUTime
	}

	utils.SuccessResponse(c, response, "获取持仓完整历史成功")
}

// getCurrencyName 获取币种中文名称
func getCurrencyName(currency models.Currency) string {
	switch currency {
	case models.CurrencyCNY:
		return "人民币"
	case models.CurrencyUSD:
		return "美元"
	case models.CurrencyUSDT:
		return "泰达币"
	case models.CurrencyBTC:
		return "比特币"
	default:
		return string(currency)
	}
}