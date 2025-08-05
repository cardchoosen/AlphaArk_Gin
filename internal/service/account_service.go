package service

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/cardchoosen/AlphaArk_Gin/internal/config"
	"github.com/cardchoosen/AlphaArk_Gin/internal/models"
)

// AccountService 账户服务接口
type AccountService interface {
	GetAccountBalance(currency models.Currency) (*models.AccountBalance, error)
	GetProfitLoss(currency models.Currency, periods []models.TimePeriod) ([]*models.ProfitLoss, error)
	GetAccountSummary(currency models.Currency) (*models.AccountSummary, error)
	SetDefaultCurrency(currency models.Currency) error
	GetDefaultCurrency() models.Currency
	GetExchangeRates() (map[string]float64, error)
	GetPositions(req *models.PositionsRequest, currency models.Currency) (*models.PositionsResponse, error)
	GetPositionsHistory(req *models.PositionsHistoryRequest, currency models.Currency) (*models.PositionsHistoryResponse, error)
}

// accountService 账户服务实现
type accountService struct {
	config          *config.OKXConfig
	defaultCurrency models.Currency
	exchangeRates   map[string]float64
	ratesMutex      sync.RWMutex
	lastRatesUpdate time.Time
	timeOffset      int64     // 与OKX服务器的时间偏移量（毫秒）
	lastSync        time.Time // 上次同步时间
}

// NewAccountService 创建账户服务实例
func NewAccountService(cfg *config.OKXConfig) AccountService {
	service := &accountService{
		config:          cfg,
		defaultCurrency: models.CurrencyUSDT, // 默认使用USDT
		exchangeRates:   make(map[string]float64),
		lastRatesUpdate: time.Time{},
		timeOffset:      0,
		lastSync:        time.Time{},
	}

	// 初始化时同步时间
	service.syncTimeWithOKX()

	return service
}

// OKXAccountBalance OKX账户余额响应
type OKXAccountBalance struct {
	Code string `json:"code"`
	Msg  string `json:"msg"`
	Data []struct {
		Details []struct {
			AvailBal  string `json:"availBal"`
			Bal       string `json:"bal"`
			FrozenBal string `json:"frozenBal"`
			Ccy       string `json:"ccy"`
		} `json:"details"`
		TotalEq string `json:"totalEq"`
		UTime   string `json:"uTime"`
	} `json:"data"`
}

// GetAccountBalance 获取账户余额
func (s *accountService) GetAccountBalance(currency models.Currency) (*models.AccountBalance, error) {
	// 检查API配置
	if s.config.APIKey == "" || s.config.SecretKey == "" || s.config.Passphrase == "" {
		return nil, fmt.Errorf("OKX API配置不完整，请检查环境变量 OKX_API_KEY, OKX_SECRET_KEY, OKX_PASSPHRASE")
	}

	// 更新汇率
	if err := s.updateExchangeRates(); err != nil {
		log.Printf("更新汇率失败: %v", err)
	}

	// 获取真实OKX账户余额
	okxBalance, err := s.fetchOKXBalance()
	if err != nil {
		return nil, fmt.Errorf("获取OKX账户余额失败: %w", err)
	}

	if len(okxBalance.Data) == 0 {
		return nil, fmt.Errorf("未获取到账户数据，请检查OKX API权限或账户状态")
	}

	accountData := okxBalance.Data[0]

	// 转换币种
	totalEquity, err := s.convertCurrency(accountData.TotalEq, models.CurrencyUSDT, currency)
	if err != nil {
		return nil, fmt.Errorf("币种转换失败: %w", err)
	}

	// 解析更新时间
	updateTime := time.Now()
	if accountData.UTime != "" {
		if timestamp, err := strconv.ParseInt(accountData.UTime, 10, 64); err == nil {
			updateTime = time.Unix(timestamp/1000, 0)
		}
	}

	// 构建余额详情
	var details []models.Balance
	for _, detail := range accountData.Details {
		if detail.Bal == "0" && detail.AvailBal == "0" {
			continue // 跳过零余额
		}

		equity, _ := s.convertCurrency(detail.Bal, models.Currency(detail.Ccy), currency)

		details = append(details, models.Balance{
			Currency:  detail.Ccy,
			Balance:   detail.Bal,
			Available: detail.AvailBal,
			Frozen:    detail.FrozenBal,
			Equity:    equity,
		})
	}

	return &models.AccountBalance{
		TotalEquity:    totalEquity,
		Currency:       currency,
		LastUpdateTime: updateTime,
		Details:        details,
	}, nil
}

// GetProfitLoss 获取盈亏信息
func (s *accountService) GetProfitLoss(currency models.Currency, periods []models.TimePeriod) ([]*models.ProfitLoss, error) {
	var profitLossList []*models.ProfitLoss

	// 获取当前余额作为基准
	currentBalance, err := s.GetAccountBalance(currency)
	if err != nil {
		return nil, fmt.Errorf("获取当前余额失败: %w", err)
	}

	currentEquity, _ := strconv.ParseFloat(currentBalance.TotalEquity, 64)

	for _, period := range periods {
		// 获取历史余额
		historicalEquity, err := s.getHistoricalEquity(currency, period)
		if err != nil {
			// 历史数据功能未实现，跳过该周期
			log.Printf("获取历史数据失败: %v", err)
			continue
		}

		// 计算盈亏
		profitAmount := currentEquity - historicalEquity
		profitPercent := 0.0
		if historicalEquity != 0 {
			profitPercent = (profitAmount / historicalEquity) * 100
		}

		endTime := time.Now()
		startTime := endTime.Add(-period.GetPeriodDuration())

		profitLoss := &models.ProfitLoss{
			Period:        period,
			ProfitAmount:  fmt.Sprintf("%.2f", profitAmount),
			ProfitPercent: fmt.Sprintf("%.2f", profitPercent),
			Currency:      currency,
			IsProfit:      profitAmount >= 0,
			StartTime:     startTime,
			EndTime:       endTime,
		}

		profitLossList = append(profitLossList, profitLoss)
	}

	return profitLossList, nil
}

// GetAccountSummary 获取账户汇总信息
func (s *accountService) GetAccountSummary(currency models.Currency) (*models.AccountSummary, error) {
	// 获取余额信息
	balance, err := s.GetAccountBalance(currency)
	if err != nil {
		return nil, fmt.Errorf("获取余额信息失败: %w", err)
	}

	// 获取盈亏信息
	periods := models.SupportedPeriods()
	profitLoss, err := s.GetProfitLoss(currency, periods)
	if err != nil {
		return nil, fmt.Errorf("获取盈亏信息失败: %w", err)
	}

	return &models.AccountSummary{
		Balance:    balance,
		ProfitLoss: profitLoss,
		Currency:   currency,
		UpdateTime: time.Now(),
	}, nil
}

// SetDefaultCurrency 设置默认币种
func (s *accountService) SetDefaultCurrency(currency models.Currency) error {
	s.defaultCurrency = currency
	return nil
}

// GetDefaultCurrency 获取默认币种
func (s *accountService) GetDefaultCurrency() models.Currency {
	return s.defaultCurrency
}

// GetExchangeRates 获取汇率信息
func (s *accountService) GetExchangeRates() (map[string]float64, error) {
	if err := s.updateExchangeRates(); err != nil {
		return nil, err
	}

	s.ratesMutex.RLock()
	defer s.ratesMutex.RUnlock()

	rates := make(map[string]float64)
	for k, v := range s.exchangeRates {
		rates[k] = v
	}

	return rates, nil
}

// fetchOKXBalance 获取OKX账户余额
func (s *accountService) fetchOKXBalance() (*OKXAccountBalance, error) {
	return s.fetchOKXBalanceWithRetry(3) // 最多重试3次
}

// fetchOKXBalanceWithRetry 带重试机制的账户余额获取
func (s *accountService) fetchOKXBalanceWithRetry(maxRetries int) (*OKXAccountBalance, error) {
	var lastErr error

	for i := 0; i < maxRetries; i++ {
		// 每次重试都使用新的时间戳，减少延迟
		if i > 0 {
			time.Sleep(time.Duration(i*50) * time.Millisecond) // 减少延迟到50ms
		}

		result, err := s.fetchOKXBalanceOnce()
		if err == nil {
			return result, nil
		}

		lastErr = err
		// 如果是时间戳过期错误，继续重试
		if strings.Contains(err.Error(), "Timestamp request expired") ||
			strings.Contains(err.Error(), "Invalid OK-ACCESS-TIMESTAMP") {
			continue
		}
		// 其他错误直接返回
		break
	}

	return nil, fmt.Errorf("重试%d次后仍然失败: %w", maxRetries, lastErr)
}

// fetchOKXBalanceOnce 单次获取OKX账户余额
func (s *accountService) fetchOKXBalanceOnce() (*OKXAccountBalance, error) {
	url := fmt.Sprintf("%s/api/v5/account/balance", s.config.BaseURL)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}

	// 添加认证头 - 使用精确时间戳
	timestamp := s.getCurrentTimestamp()
	sign := s.generateSignature(timestamp, "GET", "/api/v5/account/balance", "")

	req.Header.Set("OK-ACCESS-KEY", s.config.APIKey)
	req.Header.Set("OK-ACCESS-SIGN", sign)
	req.Header.Set("OK-ACCESS-TIMESTAMP", timestamp)
	req.Header.Set("OK-ACCESS-PASSPHRASE", s.config.Passphrase)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求失败: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}

	var result OKXAccountBalance
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}

	if result.Code != "0" {
		return nil, fmt.Errorf("OKX API错误: %s", result.Msg)
	}

	return &result, nil
}

// updateExchangeRates 更新汇率信息
func (s *accountService) updateExchangeRates() error {
	s.ratesMutex.Lock()
	defer s.ratesMutex.Unlock()

	// 如果汇率更新时间在5分钟内，跳过更新
	if time.Since(s.lastRatesUpdate) < 5*time.Minute {
		return nil
	}

	// 从OKX API获取汇率信息
	rates, err := s.fetchExchangeRatesFromOKX()
	if err != nil {
		log.Printf("从OKX获取汇率失败: %v", err)
		// 如果获取失败，保持现有汇率不变
		if len(s.exchangeRates) == 0 {
			return fmt.Errorf("无法获取汇率信息: %w", err)
		}
		return nil
	}

	s.exchangeRates = rates
	s.lastRatesUpdate = time.Now()
	return nil
}

// convertCurrency 币种转换
func (s *accountService) convertCurrency(amount string, fromCurrency, toCurrency models.Currency) (string, error) {
	if fromCurrency == toCurrency {
		return amount, nil
	}

	amountFloat, err := strconv.ParseFloat(amount, 64)
	if err != nil {
		return "0", err
	}

	s.ratesMutex.RLock()
	defer s.ratesMutex.RUnlock()

	// 转换逻辑
	var convertedAmount float64

	switch {
	case fromCurrency == models.CurrencyUSDT && toCurrency == models.CurrencyCNY:
		rate := s.exchangeRates["USDT_CNY"]
		convertedAmount = amountFloat * rate
	case fromCurrency == models.CurrencyUSD && toCurrency == models.CurrencyCNY:
		rate := s.exchangeRates["USD_CNY"]
		convertedAmount = amountFloat * rate
	case fromCurrency == models.CurrencyBTC && toCurrency == models.CurrencyUSDT:
		rate := s.exchangeRates["BTC_USDT"]
		convertedAmount = amountFloat * rate
	case fromCurrency == models.CurrencyBTC && toCurrency == models.CurrencyUSD:
		rate := s.exchangeRates["BTC_USD"]
		convertedAmount = amountFloat * rate
	case fromCurrency == models.CurrencyUSDT && toCurrency == models.CurrencyBTC:
		rate := s.exchangeRates["BTC_USDT"]
		convertedAmount = amountFloat / rate // USDT到BTC是除法
	case fromCurrency == models.CurrencyUSD && toCurrency == models.CurrencyUSDT:
		rate := s.exchangeRates["USD_USDT"]
		convertedAmount = amountFloat * rate
	default:
		// 通过USDT作为中间币种转换
		if fromCurrency != models.CurrencyUSDT {
			usdtAmount, err := s.convertCurrency(amount, fromCurrency, models.CurrencyUSDT)
			if err != nil {
				return "0", err
			}
			return s.convertCurrency(usdtAmount, models.CurrencyUSDT, toCurrency)
		}
		convertedAmount = amountFloat
	}

	// 根据目标货币调整精度
	var format string
	switch toCurrency {
	case models.CurrencyBTC:
		format = "%.5f" // BTC显示5位小数
	default:
		format = "%.2f" // 其他货币显示2位小数
	}

	return fmt.Sprintf(format, convertedAmount), nil
}

// getHistoricalEquity 获取历史权益
func (s *accountService) getHistoricalEquity(currency models.Currency, period models.TimePeriod) (float64, error) {
	// TODO: 实现从数据库或OKX历史API获取真实历史数据
	// 这里暂时返回错误，提示需要实现历史数据功能
	return 0, fmt.Errorf("历史数据功能暂未实现，需要集成数据库或OKX历史API")
}

// fetchExchangeRatesFromOKX 从OKX获取汇率信息
func (s *accountService) fetchExchangeRatesFromOKX() (map[string]float64, error) {
	rates := make(map[string]float64)

	// 需要获取的交易对列表（使用OKX实际支持的交易对）
	pairs := []string{
		"BTC-USDT",
		"BTC-USD",
		"ETH-USDT",
		"ETH-USD",
	}

	for _, pair := range pairs {
		url := fmt.Sprintf("%s/api/v5/market/ticker?instId=%s", s.config.BaseURL, pair)

		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return nil, fmt.Errorf("创建汇率请求失败: %w", err)
		}

		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{Timeout: 10 * time.Second}
		resp, err := client.Do(req)
		if err != nil {
			return nil, fmt.Errorf("获取汇率请求失败: %w", err)
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("读取汇率响应失败: %w", err)
		}

		var tickerResp struct {
			Code string `json:"code"`
			Msg  string `json:"msg"`
			Data []struct {
				InstId string `json:"instId"`
				Last   string `json:"last"`
			} `json:"data"`
		}

		if err := json.Unmarshal(body, &tickerResp); err != nil {
			log.Printf("解析%s汇率响应失败: %v", pair, err)
			continue
		}

		if tickerResp.Code != "0" || len(tickerResp.Data) == 0 {
			log.Printf("获取%s汇率失败: %s", pair, tickerResp.Msg)
			continue
		}

		if price, err := strconv.ParseFloat(tickerResp.Data[0].Last, 64); err == nil {
			rateKey := pair
			rateKey = rateKey[0:3] + "_" + rateKey[4:7] // 转换格式：BTC-USDT -> BTC_USDT
			rates[rateKey] = price

			// 如果是BTC-USD，同时计算BTC_USDT
			if pair == "BTC-USD" {
				rates["BTC_USDT"] = price // BTC_USDT ≈ BTC_USD (因为USD_USDT ≈ 1)
			}
		}
	}

	// 设置USD_USDT固定汇率
	rates["USD_USDT"] = 1.0

	// 使用第三方API获取CNY汇率
	if err := s.fetchCNYExchangeRates(rates); err != nil {
		log.Printf("获取CNY汇率失败: %v", err)
		// 使用固定汇率作为备用
		rates["USDT_CNY"] = 7.2 // 固定汇率
		rates["USD_CNY"] = 7.2  // 固定汇率
	}

	return rates, nil
}

// fetchCNYExchangeRates 从第三方API获取CNY汇率
func (s *accountService) fetchCNYExchangeRates(rates map[string]float64) error {
	// 尝试从多个汇率API获取CNY汇率
	apis := []string{
		"https://api.exchangerate-api.com/v4/latest/USD",
		"https://open.er-api.com/v6/latest/USD",
	}

	for _, apiURL := range apis {
		req, err := http.NewRequest("GET", apiURL, nil)
		if err != nil {
			continue
		}

		client := &http.Client{Timeout: 5 * time.Second}
		resp, err := client.Do(req)
		if err != nil {
			continue
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			continue
		}

		var rateResp struct {
			Rates map[string]float64 `json:"rates"`
		}

		if err := json.Unmarshal(body, &rateResp); err != nil {
			continue
		}

		if cnyRate, exists := rateResp.Rates["CNY"]; exists {
			rates["USD_CNY"] = cnyRate
			// USDT通常与USD接近1:1
			rates["USDT_CNY"] = cnyRate
			return nil
		}
	}

	return fmt.Errorf("无法从任何API获取CNY汇率")
}

// GetPositionsHistory 获取历史持仓信息
func (s *accountService) GetPositionsHistory(req *models.PositionsHistoryRequest, currency models.Currency) (*models.PositionsHistoryResponse, error) {
	// 检查API配置
	if s.config.APIKey == "" || s.config.SecretKey == "" || s.config.Passphrase == "" {
		return nil, fmt.Errorf("OKX API配置不完整，请检查环境变量 OKX_API_KEY, OKX_SECRET_KEY, OKX_PASSPHRASE")
	}

	// 先同步OKX服务器时间
	if err := s.syncTimeWithOKX(); err != nil {
		log.Printf("时间同步失败: %v", err)
		// 继续执行，不因为时间同步失败而中断
	}

	// 获取真实OKX历史持仓数据
	okxPositions, err := s.fetchOKXPositionsHistory(req)
	if err != nil {
		return nil, fmt.Errorf("获取OKX历史持仓信息失败: %w", err)
	}

	if len(okxPositions.Data) == 0 {
		return &models.PositionsHistoryResponse{
			Positions: []*models.PositionHistory{},
			HasMore:   false,
			Currency:  currency,
		}, nil
	}

	var positions []*models.PositionHistory
	for _, pos := range okxPositions.Data {
		// 解析时间戳
		var updateTime, createTime time.Time
		if pos.UTime != "" {
			if timestamp, err := strconv.ParseInt(pos.UTime, 10, 64); err == nil {
				updateTime = time.Unix(timestamp/1000, 0)
			}
		}
		if pos.CTime != "" {
			if timestamp, err := strconv.ParseInt(pos.CTime, 10, 64); err == nil {
				createTime = time.Unix(timestamp/1000, 0)
			}
		}

		position := &models.PositionHistory{
			InstType:       pos.InstType,
			InstId:         pos.InstId,
			MgnMode:        pos.MgnMode,
			Type:           pos.Type,
			CTime:          pos.CTime,
			UTime:          pos.UTime,
			OpenAvgPx:      pos.OpenAvgPx,
			NonSettleAvgPx: pos.NonSettleAvgPx,
			CloseAvgPx:     pos.CloseAvgPx,
			PosId:          pos.PosId,
			OpenMaxPos:     pos.OpenMaxPos,
			CloseTotalPos:  pos.CloseTotalPos,
			RealizedPnl:    pos.RealizedPnl,
			SettledPnl:     pos.SettledPnl,
			PnlRatio:       pos.PnlRatio,
			Fee:            pos.Fee,
			FundingFee:     pos.FundingFee,
			LiqPenalty:     pos.LiqPenalty,
			Pnl:            pos.Pnl,
			PosSide:        pos.PosSide,
			Lever:          pos.Lever,
			Direction:      pos.Direction,
			TriggerPx:      pos.TriggerPx,
			Uly:            pos.Uly,
			Ccy:            pos.Ccy,
			Currency:       currency,
			UpdateTime:     updateTime,
			CreateTime:     createTime,
		}

		positions = append(positions, position)
	}

	return &models.PositionsHistoryResponse{
		Positions: positions,
		HasMore:   len(positions) >= 100, // OKX最大返回100条，如果等于100可能还有更多
		Currency:  currency,
	}, nil
}

// OKXPositionsHistoryResponse OKX历史持仓响应
type OKXPositionsHistoryResponse struct {
	Code string `json:"code"`
	Msg  string `json:"msg"`
	Data []struct {
		InstType       string `json:"instType"`       // 产品类型
		InstId         string `json:"instId"`         // 交易产品ID
		MgnMode        string `json:"mgnMode"`        // 保证金模式
		Type           string `json:"type"`           // 平仓类型
		CTime          string `json:"cTime"`          // 仓位创建时间
		UTime          string `json:"uTime"`          // 仓位更新时间
		OpenAvgPx      string `json:"openAvgPx"`      // 开仓均价
		NonSettleAvgPx string `json:"nonSettleAvgPx"` // 未结算均价
		CloseAvgPx     string `json:"closeAvgPx"`     // 平仓均价
		PosId          string `json:"posId"`          // 仓位ID
		OpenMaxPos     string `json:"openMaxPos"`     // 最大持仓量
		CloseTotalPos  string `json:"closeTotalPos"`  // 累计平仓量
		RealizedPnl    string `json:"realizedPnl"`    // 已实现收益
		SettledPnl     string `json:"settledPnl"`     // 已实现收益(全仓交割)
		PnlRatio       string `json:"pnlRatio"`       // 已实现收益率
		Fee            string `json:"fee"`            // 累计手续费金额
		FundingFee     string `json:"fundingFee"`     // 累计资金费用
		LiqPenalty     string `json:"liqPenalty"`     // 累计爆仓罚金
		Pnl            string `json:"pnl"`            // 已实现收益
		PosSide        string `json:"posSide"`        // 持仓模式方向
		Lever          string `json:"lever"`          // 杠杆倍数
		Direction      string `json:"direction"`      // 持仓方向
		TriggerPx      string `json:"triggerPx"`      // 触发标记价格
		Uly            string `json:"uly"`            // 标的指数
		Ccy            string `json:"ccy"`            // 占用保证金的币种
	} `json:"data"`
}

// fetchOKXPositionsHistory 获取OKX历史持仓数据
func (s *accountService) fetchOKXPositionsHistory(req *models.PositionsHistoryRequest) (*OKXPositionsHistoryResponse, error) {
	return s.fetchOKXPositionsHistoryWithRetry(req, 3) // 最多重试3次
}

// fetchOKXPositionsHistoryWithRetry 带重试机制的历史持仓数据获取
func (s *accountService) fetchOKXPositionsHistoryWithRetry(req *models.PositionsHistoryRequest, maxRetries int) (*OKXPositionsHistoryResponse, error) {
	var lastErr error

	for i := 0; i < maxRetries; i++ {
		// 每次重试都使用新的时间戳，减少延迟
		if i > 0 {
			time.Sleep(time.Duration(i*50) * time.Millisecond) // 减少延迟到50ms
		}

		result, err := s.fetchOKXPositionsHistoryOnce(req)
		if err == nil {
			return result, nil
		}

		lastErr = err
		// 如果是时间戳过期错误，继续重试
		if strings.Contains(err.Error(), "Timestamp request expired") ||
			strings.Contains(err.Error(), "Invalid OK-ACCESS-TIMESTAMP") {
			continue
		}
		// 其他错误直接返回
		break
	}

	return nil, fmt.Errorf("重试%d次后仍然失败: %w", maxRetries, lastErr)
}

// fetchOKXPositionsHistoryOnce 单次获取OKX历史持仓数据
func (s *accountService) fetchOKXPositionsHistoryOnce(req *models.PositionsHistoryRequest) (*OKXPositionsHistoryResponse, error) {
	// 构建查询参数
	params := make(map[string]string)
	if req.InstType != "" {
		params["instType"] = req.InstType
	}
	if req.InstId != "" {
		params["instId"] = req.InstId
	}
	if req.MgnMode != "" {
		params["mgnMode"] = req.MgnMode
	}
	if req.Type != "" {
		params["type"] = req.Type
	}
	if req.PosId != "" {
		params["posId"] = req.PosId
	}
	if req.After != "" {
		params["after"] = req.After
	}
	if req.Before != "" {
		params["before"] = req.Before
	}
	if req.Limit != "" {
		params["limit"] = req.Limit
	} else {
		params["limit"] = "100" // 默认100条
	}

	// 构建URL
	url := fmt.Sprintf("%s/api/v5/account/positions-history", s.config.BaseURL)
	if len(params) > 0 {
		var queryParams []string
		for key, value := range params {
			queryParams = append(queryParams, fmt.Sprintf("%s=%s", key, value))
		}
		url += "?" + strings.Join(queryParams, "&")
	}

	req_obj, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}

	// 添加认证头
	timestamp := s.getCurrentTimestamp()
	requestPath := "/api/v5/account/positions-history"
	if len(params) > 0 {
		var queryParams []string
		for key, value := range params {
			queryParams = append(queryParams, fmt.Sprintf("%s=%s", key, value))
		}
		requestPath += "?" + strings.Join(queryParams, "&")
	}
	sign := s.generateSignature(timestamp, "GET", requestPath, "")

	req_obj.Header.Set("OK-ACCESS-KEY", s.config.APIKey)
	req_obj.Header.Set("OK-ACCESS-SIGN", sign)
	req_obj.Header.Set("OK-ACCESS-TIMESTAMP", timestamp)
	req_obj.Header.Set("OK-ACCESS-PASSPHRASE", s.config.Passphrase)
	req_obj.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req_obj)
	if err != nil {
		return nil, fmt.Errorf("请求失败: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}

	var result OKXPositionsHistoryResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}

	if result.Code != "0" {
		return nil, fmt.Errorf("OKX API错误: %s", result.Msg)
	}

	return &result, nil
}

// GetPositions 获取当前持仓信息
func (s *accountService) GetPositions(req *models.PositionsRequest, currency models.Currency) (*models.PositionsResponse, error) {
	// 检查API配置
	if s.config.APIKey == "" || s.config.SecretKey == "" || s.config.Passphrase == "" {
		return nil, fmt.Errorf("OKX API配置不完整，请检查环境变量 OKX_API_KEY, OKX_SECRET_KEY, OKX_PASSPHRASE")
	}

	// 先同步OKX服务器时间
	if err := s.syncTimeWithOKX(); err != nil {
		log.Printf("时间同步失败: %v", err)
		// 继续执行，不因为时间同步失败而中断
	}

	// 获取真实OKX当前持仓数据
	okxPositions, err := s.fetchOKXPositions(req)
	if err != nil {
		return nil, fmt.Errorf("获取OKX当前持仓信息失败: %w", err)
	}

	if len(okxPositions.Data) == 0 {
		return &models.PositionsResponse{
			Positions: []*models.Position{},
			Currency:  currency,
		}, nil
	}

	var positions []*models.Position
	for _, pos := range okxPositions.Data {
		// 解析时间戳
		var updateTime, createTime time.Time
		if pos.UTime != "" {
			if timestamp, err := strconv.ParseInt(pos.UTime, 10, 64); err == nil {
				updateTime = time.Unix(timestamp/1000, 0)
			}
		}
		if pos.CTime != "" {
			if timestamp, err := strconv.ParseInt(pos.CTime, 10, 64); err == nil {
				createTime = time.Unix(timestamp/1000, 0)
			}
		}

		// 解析平仓策略委托订单
		var closeOrderAlgo []models.CloseOrderAlgo
		for _, algo := range pos.CloseOrderAlgo {
			closeOrderAlgo = append(closeOrderAlgo, models.CloseOrderAlgo{
				AlgoId:          algo.AlgoId,
				SlTriggerPx:     algo.SlTriggerPx,
				SlTriggerPxType: algo.SlTriggerPxType,
				TpTriggerPx:     algo.TpTriggerPx,
				TpTriggerPxType: algo.TpTriggerPxType,
				CloseFraction:   algo.CloseFraction,
			})
		}

		position := &models.Position{
			InstType:               pos.InstType,
			InstId:                 pos.InstId,
			MgnMode:                pos.MgnMode,
			PosId:                  pos.PosId,
			PosSide:                pos.PosSide,
			Pos:                    pos.Pos,
			BaseBal:                pos.BaseBal,
			QuoteBal:               pos.QuoteBal,
			BaseBorrowed:           pos.BaseBorrowed,
			BaseInterest:           pos.BaseInterest,
			QuoteBorrowed:          pos.QuoteBorrowed,
			QuoteInterest:          pos.QuoteInterest,
			PosCcy:                 pos.PosCcy,
			AvailPos:               pos.AvailPos,
			AvgPx:                  pos.AvgPx,
			NonSettleAvgPx:         pos.NonSettleAvgPx,
			Upl:                    pos.Upl,
			UplRatio:               pos.UplRatio,
			UplLastPx:              pos.UplLastPx,
			UplRatioLastPx:         pos.UplRatioLastPx,
			Lever:                  pos.Lever,
			LiqPx:                  pos.LiqPx,
			MarkPx:                 pos.MarkPx,
			Imr:                    pos.Imr,
			Margin:                 pos.Margin,
			MgnRatio:               pos.MgnRatio,
			Mmr:                    pos.Mmr,
			Liab:                   pos.Liab,
			LiabCcy:                pos.LiabCcy,
			Interest:               pos.Interest,
			TradeId:                pos.TradeId,
			OptVal:                 pos.OptVal,
			PendingCloseOrdLiabVal: pos.PendingCloseOrdLiabVal,
			NotionalUsd:            pos.NotionalUsd,
			Adl:                    pos.Adl,
			Ccy:                    pos.Ccy,
			Last:                   pos.Last,
			IdxPx:                  pos.IdxPx,
			UsdPx:                  pos.UsdPx,
			BePx:                   pos.BePx,
			DeltaBS:                pos.DeltaBS,
			DeltaPA:                pos.DeltaPA,
			GammaBS:                pos.GammaBS,
			GammaPA:                pos.GammaPA,
			ThetaBS:                pos.ThetaBS,
			ThetaPA:                pos.ThetaPA,
			VegaBS:                 pos.VegaBS,
			VegaPA:                 pos.VegaPA,
			SpotInUseAmt:           pos.SpotInUseAmt,
			SpotInUseCcy:           pos.SpotInUseCcy,
			ClSpotInUseAmt:         pos.ClSpotInUseAmt,
			MaxSpotInUseAmt:        pos.MaxSpotInUseAmt,
			RealizedPnl:            pos.RealizedPnl,
			SettledPnl:             pos.SettledPnl,
			Pnl:                    pos.Pnl,
			Fee:                    pos.Fee,
			FundingFee:             pos.FundingFee,
			LiqPenalty:             pos.LiqPenalty,
			CloseOrderAlgo:         closeOrderAlgo,
			CTime:                  pos.CTime,
			UTime:                  pos.UTime,
			BizRefId:               pos.BizRefId,
			BizRefType:             pos.BizRefType,
			Currency:               currency,
			CreateTime:             createTime,
			UpdateTime:             updateTime,
		}

		positions = append(positions, position)
	}

	return &models.PositionsResponse{
		Positions: positions,
		Currency:  currency,
	}, nil
}

// OKXPositionsResponse OKX当前持仓响应
type OKXPositionsResponse struct {
	Code string `json:"code"`
	Msg  string `json:"msg"`
	Data []struct {
		InstType               string `json:"instType"`               // 产品类型
		InstId                 string `json:"instId"`                 // 产品ID
		MgnMode                string `json:"mgnMode"`                // 保证金模式
		PosId                  string `json:"posId"`                  // 持仓ID
		PosSide                string `json:"posSide"`                // 持仓方向
		Pos                    string `json:"pos"`                    // 持仓数量
		BaseBal                string `json:"baseBal"`                // 交易币余额
		QuoteBal               string `json:"quoteBal"`               // 计价币余额
		BaseBorrowed           string `json:"baseBorrowed"`           // 交易币已借
		BaseInterest           string `json:"baseInterest"`           // 交易币计息
		QuoteBorrowed          string `json:"quoteBorrowed"`          // 计价币已借
		QuoteInterest          string `json:"quoteInterest"`          // 计价币计息
		PosCcy                 string `json:"posCcy"`                 // 仓位资产币种
		AvailPos               string `json:"availPos"`               // 可平仓数量
		AvgPx                  string `json:"avgPx"`                  // 开仓均价
		NonSettleAvgPx         string `json:"nonSettleAvgPx"`         // 未结算均价
		Upl                    string `json:"upl"`                    // 未实现收益
		UplRatio               string `json:"uplRatio"`               // 未实现收益率
		UplLastPx              string `json:"uplLastPx"`              // 以最新成交价计算的未实现收益
		UplRatioLastPx         string `json:"uplRatioLastPx"`         // 以最新成交价计算的未实现收益率
		Lever                  string `json:"lever"`                  // 杠杆倍数
		LiqPx                  string `json:"liqPx"`                  // 预估强平价
		MarkPx                 string `json:"markPx"`                 // 最新标记价格
		Imr                    string `json:"imr"`                    // 初始保证金
		Margin                 string `json:"margin"`                 // 保证金余额
		MgnRatio               string `json:"mgnRatio"`               // 维持保证金率
		Mmr                    string `json:"mmr"`                    // 维持保证金
		Liab                   string `json:"liab"`                   // 负债额
		LiabCcy                string `json:"liabCcy"`                // 负债币种
		Interest               string `json:"interest"`               // 利息
		TradeId                string `json:"tradeId"`                // 最新成交ID
		OptVal                 string `json:"optVal"`                 // 期权市值
		PendingCloseOrdLiabVal string `json:"pendingCloseOrdLiabVal"` // 逐仓杠杆负债对应平仓挂单的数量
		NotionalUsd            string `json:"notionalUsd"`            // 以美金价值为单位的持仓数量
		Adl                    string `json:"adl"`                    // 信号区
		Ccy                    string `json:"ccy"`                    // 占用保证金的币种
		Last                   string `json:"last"`                   // 最新成交价
		IdxPx                  string `json:"idxPx"`                  // 最新指数价格
		UsdPx                  string `json:"usdPx"`                  // 保证金币种的市场最新美金价格
		BePx                   string `json:"bePx"`                   // 盈亏平衡价
		DeltaBS                string `json:"deltaBS"`                // 美金本位持仓仓位delta
		DeltaPA                string `json:"deltaPA"`                // 币本位持仓仓位delta
		GammaBS                string `json:"gammaBS"`                // 美金本位持仓仓位gamma
		GammaPA                string `json:"gammaPA"`                // 币本位持仓仓位gamma
		ThetaBS                string `json:"thetaBS"`                // 美金本位持仓仓位theta
		ThetaPA                string `json:"thetaPA"`                // 币本位持仓仓位theta
		VegaBS                 string `json:"vegaBS"`                 // 美金本位持仓仓位vega
		VegaPA                 string `json:"vegaPA"`                 // 币本位持仓仓位vega
		SpotInUseAmt           string `json:"spotInUseAmt"`           // 现货对冲占用数量
		SpotInUseCcy           string `json:"spotInUseCcy"`           // 现货对冲占用币种
		ClSpotInUseAmt         string `json:"clSpotInUseAmt"`         // 用户自定义现货占用数量
		MaxSpotInUseAmt        string `json:"maxSpotInUseAmt"`        // 系统计算得到的最大可能现货占用数量
		RealizedPnl            string `json:"realizedPnl"`            // 已实现收益
		SettledPnl             string `json:"settledPnl"`             // 已结算收益
		Pnl                    string `json:"pnl"`                    // 平仓订单累计收益额
		Fee                    string `json:"fee"`                    // 累计手续费金额
		FundingFee             string `json:"fundingFee"`             // 累计资金费用
		LiqPenalty             string `json:"liqPenalty"`             // 累计爆仓罚金
		CloseOrderAlgo         []struct {
			AlgoId          string `json:"algoId"`          // 策略委托单ID
			SlTriggerPx     string `json:"slTriggerPx"`     // 止损触发价
			SlTriggerPxType string `json:"slTriggerPxType"` // 止损触发价类型
			TpTriggerPx     string `json:"tpTriggerPx"`     // 止盈委托价
			TpTriggerPxType string `json:"tpTriggerPxType"` // 止盈触发价类型
			CloseFraction   string `json:"closeFraction"`   // 策略委托触发时，平仓的百分比
		} `json:"closeOrderAlgo"` // 平仓策略委托订单
		CTime      string `json:"cTime"`      // 持仓创建时间
		UTime      string `json:"uTime"`      // 最近一次持仓更新时间
		BizRefId   string `json:"bizRefId"`   // 外部业务id
		BizRefType string `json:"bizRefType"` // 外部业务类型
	} `json:"data"`
}

// fetchOKXPositions 获取OKX当前持仓数据
func (s *accountService) fetchOKXPositions(req *models.PositionsRequest) (*OKXPositionsResponse, error) {
	return s.fetchOKXPositionsWithRetry(req, 3) // 最多重试3次
}

// fetchOKXPositionsWithRetry 带重试机制的持仓数据获取
func (s *accountService) fetchOKXPositionsWithRetry(req *models.PositionsRequest, maxRetries int) (*OKXPositionsResponse, error) {
	var lastErr error

	for i := 0; i < maxRetries; i++ {
		// 每次重试都使用新的时间戳，减少延迟
		if i > 0 {
			time.Sleep(time.Duration(i*50) * time.Millisecond) // 减少延迟到50ms
		}

		result, err := s.fetchOKXPositionsOnce(req)
		if err == nil {
			return result, nil
		}

		lastErr = err
		// 如果是时间戳过期错误，继续重试
		if strings.Contains(err.Error(), "Timestamp request expired") ||
			strings.Contains(err.Error(), "Invalid OK-ACCESS-TIMESTAMP") {
			continue
		}
		// 其他错误直接返回
		break
	}

	return nil, fmt.Errorf("重试%d次后仍然失败: %w", maxRetries, lastErr)
}

// fetchOKXPositionsOnce 单次获取OKX当前持仓数据
func (s *accountService) fetchOKXPositionsOnce(req *models.PositionsRequest) (*OKXPositionsResponse, error) {
	// 构建查询参数
	params := make(map[string]string)
	if req.InstType != "" {
		params["instType"] = req.InstType
	}
	if req.InstId != "" {
		params["instId"] = req.InstId
	}
	if req.PosId != "" {
		params["posId"] = req.PosId
	}

	// 构建URL
	url := fmt.Sprintf("%s/api/v5/account/positions", s.config.BaseURL)
	if len(params) > 0 {
		var queryParams []string
		for key, value := range params {
			queryParams = append(queryParams, fmt.Sprintf("%s=%s", key, value))
		}
		url += "?" + strings.Join(queryParams, "&")
	}

	req_obj, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}

	// 添加认证头
	timestamp := s.getCurrentTimestamp()
	requestPath := "/api/v5/account/positions"
	if len(params) > 0 {
		var queryParams []string
		for key, value := range params {
			queryParams = append(queryParams, fmt.Sprintf("%s=%s", key, value))
		}
		requestPath += "?" + strings.Join(queryParams, "&")
	}
	sign := s.generateSignature(timestamp, "GET", requestPath, "")

	req_obj.Header.Set("OK-ACCESS-KEY", s.config.APIKey)
	req_obj.Header.Set("OK-ACCESS-SIGN", sign)
	req_obj.Header.Set("OK-ACCESS-TIMESTAMP", timestamp)
	req_obj.Header.Set("OK-ACCESS-PASSPHRASE", s.config.Passphrase)
	req_obj.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req_obj)
	if err != nil {
		return nil, fmt.Errorf("请求失败: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}

	var result OKXPositionsResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}

	if result.Code != "0" {
		return nil, fmt.Errorf("OKX API错误: %s", result.Msg)
	}

	return &result, nil
}

// generateSignature 生成OKX API签名
func (s *accountService) generateSignature(timestamp, method, requestPath, body string) string {
	message := timestamp + method + requestPath + body
	mac := hmac.New(sha256.New, []byte(s.config.SecretKey))
	mac.Write([]byte(message))
	return base64.StdEncoding.EncodeToString(mac.Sum(nil))
}

// getCurrentTimestamp 获取当前时间戳（ISO 8601格式）
func (s *accountService) getCurrentTimestamp() string {
	// 确保时间同步
	s.syncTimeWithOKX()

	// 使用同步后的时间生成ISO 8601格式的时间戳
	syncedTime := time.Now().Add(time.Duration(s.timeOffset) * time.Millisecond)
	return syncedTime.UTC().Format("2006-01-02T15:04:05.000Z")
}

// getOKXServerTime 获取OKX服务器时间戳
func (s *accountService) getOKXServerTime() (string, error) {
	url := fmt.Sprintf("%s/api/v5/public/time", s.config.BaseURL)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("创建时间请求失败: %w", err)
	}

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("时间请求失败: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("读取时间响应失败: %w", err)
	}

	var timeResp struct {
		Code string `json:"code"`
		Msg  string `json:"msg"`
		Data []struct {
			TS string `json:"ts"`
		} `json:"data"`
	}

	if err := json.Unmarshal(body, &timeResp); err != nil {
		return "", fmt.Errorf("解析时间响应失败: %w", err)
	}

	if timeResp.Code != "0" || len(timeResp.Data) == 0 {
		return "", fmt.Errorf("获取OKX服务器时间失败: %s", timeResp.Msg)
	}

	return timeResp.Data[0].TS, nil
}

// syncTimeWithOKX 与OKX服务器时间同步
func (s *accountService) syncTimeWithOKX() error {
	// 如果距离上次同步时间不足5分钟，跳过同步
	if !s.lastSync.IsZero() && time.Since(s.lastSync) < 5*time.Minute {
		return nil
	}

	// 获取OKX服务器时间
	serverTime, err := s.getOKXServerTime()
	if err != nil {
		return fmt.Errorf("获取服务器时间失败: %w", err)
	}

	// 解析服务器时间戳（毫秒）
	serverTs, err := strconv.ParseInt(serverTime, 10, 64)
	if err != nil {
		return fmt.Errorf("解析服务器时间戳失败: %w", err)
	}

	// 计算时间偏移量（服务器时间 - 本地时间）
	localTs := time.Now().UnixMilli()
	s.timeOffset = serverTs - localTs
	s.lastSync = time.Now()

	log.Printf("时间同步完成: 偏移量=%dms, 服务器时间=%s", s.timeOffset, serverTime)

	return nil
}
