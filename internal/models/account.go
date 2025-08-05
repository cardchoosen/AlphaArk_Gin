package models

import "time"

// Currency 支持的币种单位
type Currency string

const (
	CurrencyCNY  Currency = "CNY"
	CurrencyUSD  Currency = "USD"
	CurrencyUSDT Currency = "USDT"
	CurrencyBTC  Currency = "BTC"
)

// TimePeriod 时间周期
type TimePeriod string

const (
	Period1Day   TimePeriod = "1d"
	Period1Week  TimePeriod = "1w"
	Period1Month TimePeriod = "1m"
	Period6Month TimePeriod = "6m"
)

// AccountBalance 账户余额信息
type AccountBalance struct {
	TotalEquity    string    `json:"totalEquity"`    // 总资产
	Currency       Currency  `json:"currency"`       // 当前显示币种
	LastUpdateTime time.Time `json:"lastUpdateTime"` // 最后更新时间
	Details        []Balance `json:"details"`        // 详细余额
}

// Balance 单币种余额
type Balance struct {
	Currency  string `json:"currency"`  // 币种
	Balance   string `json:"balance"`   // 余额
	Available string `json:"available"` // 可用余额
	Frozen    string `json:"frozen"`    // 冻结余额
	Equity    string `json:"equity"`    // 权益（按当前显示币种计算）
}

// ProfitLoss 盈亏信息
type ProfitLoss struct {
	Period         TimePeriod `json:"period"`         // 时间周期
	ProfitAmount   string     `json:"profitAmount"`   // 盈亏金额
	ProfitPercent  string     `json:"profitPercent"`  // 盈亏百分比
	Currency       Currency   `json:"currency"`       // 币种单位
	IsProfit       bool       `json:"isProfit"`       // 是否盈利
	StartTime      time.Time  `json:"startTime"`      // 开始时间
	EndTime        time.Time  `json:"endTime"`        // 结束时间
}

// AccountSummary 账户汇总信息
type AccountSummary struct {
	Balance     *AccountBalance `json:"balance"`     // 余额信息
	ProfitLoss  []*ProfitLoss   `json:"profitLoss"`  // 盈亏信息
	Currency    Currency        `json:"currency"`    // 当前显示币种
	UpdateTime  time.Time       `json:"updateTime"`  // 更新时间
}

// CurrencySettings 币种设置
type CurrencySettings struct {
	DefaultCurrency Currency           `json:"defaultCurrency"` // 默认币种
	ExchangeRates   map[string]float64 `json:"exchangeRates"`   // 汇率信息
	LastUpdate      time.Time          `json:"lastUpdate"`      // 汇率更新时间
}

// SupportedCurrencies 获取支持的币种列表
func SupportedCurrencies() []Currency {
	return []Currency{CurrencyCNY, CurrencyUSD, CurrencyUSDT, CurrencyBTC}
}

// SupportedPeriods 获取支持的时间周期列表
func SupportedPeriods() []TimePeriod {
	return []TimePeriod{Period1Day, Period1Week, Period1Month, Period6Month}
}

// GetPeriodDuration 获取时间周期对应的持续时间
func (p TimePeriod) GetPeriodDuration() time.Duration {
	switch p {
	case Period1Day:
		return 24 * time.Hour
	case Period1Week:
		return 7 * 24 * time.Hour
	case Period1Month:
		return 30 * 24 * time.Hour
	case Period6Month:
		return 180 * 24 * time.Hour
	default:
		return 24 * time.Hour
	}
}

// GetPeriodName 获取时间周期的中文名称
func (p TimePeriod) GetPeriodName() string {
	switch p {
	case Period1Day:
		return "1日"
	case Period1Week:
		return "1周"
	case Period1Month:
		return "1月"
	case Period6Month:
		return "半年"
	default:
		return "1日"
	}
}

// GetCurrencySymbol 获取币种符号
func (c Currency) GetCurrencySymbol() string {
	switch c {
	case CurrencyCNY:
		return "¥"
	case CurrencyUSD:
		return "$"
	case CurrencyUSDT:
		return "₮"
	case CurrencyBTC:
		return "₿"
	default:
		return "$"
	}
}