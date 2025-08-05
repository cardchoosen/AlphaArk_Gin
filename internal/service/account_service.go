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
}

// accountService 账户服务实现
type accountService struct {
	config          *config.OKXConfig
	defaultCurrency models.Currency
	exchangeRates   map[string]float64
	ratesMutex      sync.RWMutex
	lastRatesUpdate time.Time
}

// NewAccountService 创建账户服务实例
func NewAccountService(cfg *config.OKXConfig) AccountService {
	return &accountService{
		config:          cfg,
		defaultCurrency: models.CurrencyUSDT, // 默认使用USDT
		exchangeRates:   make(map[string]float64),
		lastRatesUpdate: time.Time{},
	}
}

// OKXAccountBalance OKX账户余额响应
type OKXAccountBalance struct {
	Code string `json:"code"`
	Msg  string `json:"msg"`
	Data []struct {
		Details []struct {
			AvailBal string `json:"availBal"`
			Bal      string `json:"bal"`
			FrozenBal string `json:"frozenBal"`
			Ccy      string `json:"ccy"`
		} `json:"details"`
		TotalEq string `json:"totalEq"`
		UTime   string `json:"uTime"`
	} `json:"data"`
}

// GetAccountBalance 获取账户余额
func (s *accountService) GetAccountBalance(currency models.Currency) (*models.AccountBalance, error) {
	// 更新汇率
	if err := s.updateExchangeRates(); err != nil {
		log.Printf("更新汇率失败: %v", err)
	}

	// 尝试获取真实OKX账户余额，失败时使用模拟数据
	okxBalance, err := s.fetchOKXBalance()
	if err != nil {
		log.Printf("获取真实OKX余额失败，使用模拟数据: %v", err)
		return s.getMockAccountBalance(currency)
	}

	if len(okxBalance.Data) == 0 {
		log.Printf("未获取到账户数据，使用模拟数据")
		return s.getMockAccountBalance(currency)
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
		// 计算历史余额（这里使用模拟数据，实际应该从历史数据获取）
		historicalEquity := s.getHistoricalEquity(currency, period)
		
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
	url := fmt.Sprintf("%s/api/v5/account/balance", s.config.BaseURL)
	
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}

	// 添加认证头
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
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

	// 这里应该从实际的汇率API获取，现在使用模拟数据
	s.exchangeRates = map[string]float64{
		"USDT_CNY": 7.25,   // 1 USDT = 7.25 CNY
		"USD_CNY":  7.30,   // 1 USD = 7.30 CNY
		"BTC_USDT": 43500.0, // 1 BTC = 43500 USDT
		"BTC_USD":  43200.0, // 1 BTC = 43200 USD
		"USD_USDT": 1.0,    // 1 USD = 1.0 USDT
	}

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

	return fmt.Sprintf("%.2f", convertedAmount), nil
}

// getHistoricalEquity 获取历史权益（模拟数据）
func (s *accountService) getHistoricalEquity(currency models.Currency, period models.TimePeriod) float64 {
	// 这里使用模拟的历史数据，实际应该从数据库或历史API获取
	baseAmount := 10000.0 // 基础金额

	switch period {
	case models.Period1Day:
		return baseAmount * 0.98 // 模拟1日前少2%
	case models.Period1Week:
		return baseAmount * 0.95 // 模拟1周前少5%
	case models.Period1Month:
		return baseAmount * 0.90 // 模拟1月前少10%
	case models.Period6Month:
		return baseAmount * 0.80 // 模拟半年前少20%
	default:
		return baseAmount
	}
}

// getMockAccountBalance 获取模拟账户余额数据
func (s *accountService) getMockAccountBalance(currency models.Currency) (*models.AccountBalance, error) {
	// 模拟基础USDT余额
	baseUSDTBalance := 10000.0
	
	// 转换到目标币种
	totalEquity, err := s.convertCurrency(fmt.Sprintf("%.2f", baseUSDTBalance), models.CurrencyUSDT, currency)
	if err != nil {
		return nil, fmt.Errorf("币种转换失败: %w", err)
	}

	// 模拟详细余额
	details := []models.Balance{
		{
			Currency:  "USDT",
			Balance:   "8500.00",
			Available: "8000.00",
			Frozen:    "500.00",
			Equity:    totalEquity,
		},
		{
			Currency:  "BTC",
			Balance:   "0.05",
			Available: "0.05",
			Frozen:    "0.00",
			Equity:    "2175.00", // 假设BTC价格43500
		},
	}

	// 如果目标币种不是USDT，需要转换详细余额的权益
	if currency != models.CurrencyUSDT {
		for i := range details {
			if details[i].Currency != "USDT" {
				convertedEquity, _ := s.convertCurrency(details[i].Equity, models.CurrencyUSDT, currency)
				details[i].Equity = convertedEquity
			}
		}
	}

	return &models.AccountBalance{
		TotalEquity:    totalEquity,
		Currency:       currency,
		LastUpdateTime: time.Now(),
		Details:        details,
	}, nil
}

// generateSignature 生成OKX API签名
func (s *accountService) generateSignature(timestamp, method, requestPath, body string) string {
	message := timestamp + method + requestPath + body
	mac := hmac.New(sha256.New, []byte(s.config.SecretKey))
	mac.Write([]byte(message))
	return base64.StdEncoding.EncodeToString(mac.Sum(nil))
}