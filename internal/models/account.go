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

// PositionHistory 历史持仓信息
type PositionHistory struct {
	InstType       string    `json:"instType"`       // 产品类型
	InstId         string    `json:"instId"`         // 交易产品ID
	MgnMode        string    `json:"mgnMode"`        // 保证金模式
	Type           string    `json:"type"`           // 平仓类型
	CTime          string    `json:"cTime"`          // 仓位创建时间
	UTime          string    `json:"uTime"`          // 仓位更新时间
	OpenAvgPx      string    `json:"openAvgPx"`      // 开仓均价
	NonSettleAvgPx string    `json:"nonSettleAvgPx"` // 未结算均价
	CloseAvgPx     string    `json:"closeAvgPx"`     // 平仓均价
	PosId          string    `json:"posId"`          // 仓位ID
	OpenMaxPos     string    `json:"openMaxPos"`     // 最大持仓量
	CloseTotalPos  string    `json:"closeTotalPos"`  // 累计平仓量
	RealizedPnl    string    `json:"realizedPnl"`    // 已实现收益
	SettledPnl     string    `json:"settledPnl"`     // 已实现收益(全仓交割)
	PnlRatio       string    `json:"pnlRatio"`       // 已实现收益率
	Fee            string    `json:"fee"`            // 累计手续费金额
	FundingFee     string    `json:"fundingFee"`     // 累计资金费用
	LiqPenalty     string    `json:"liqPenalty"`     // 累计爆仓罚金
	Pnl            string    `json:"pnl"`            // 已实现收益
	PosSide        string    `json:"posSide"`        // 持仓模式方向
	Lever          string    `json:"lever"`          // 杠杆倍数
	Direction      string    `json:"direction"`      // 持仓方向
	TriggerPx      string    `json:"triggerPx"`      // 触发标记价格
	Uly            string    `json:"uly"`            // 标的指数
	Ccy            string    `json:"ccy"`            // 占用保证金的币种
	Currency       Currency  `json:"currency"`       // 显示币种
	UpdateTime     time.Time `json:"updateTime"`     // 更新时间
	CreateTime     time.Time `json:"createTime"`     // 创建时间
}

// PositionsHistoryRequest 历史持仓查询请求
type PositionsHistoryRequest struct {
	InstType string `json:"instType,omitempty" form:"instType"` // 产品类型
	InstId   string `json:"instId,omitempty" form:"instId"`     // 交易产品ID
	MgnMode  string `json:"mgnMode,omitempty" form:"mgnMode"`   // 保证金模式
	Type     string `json:"type,omitempty" form:"type"`         // 平仓类型
	PosId    string `json:"posId,omitempty" form:"posId"`       // 持仓ID
	After    string `json:"after,omitempty" form:"after"`       // 查询之前的内容
	Before   string `json:"before,omitempty" form:"before"`     // 查询之后的内容
	Limit    string `json:"limit,omitempty" form:"limit"`       // 分页数量
}

// PositionsHistoryResponse 历史持仓查询响应
type PositionsHistoryResponse struct {
	Positions []*PositionHistory `json:"positions"` // 历史持仓列表
	HasMore   bool               `json:"hasMore"`   // 是否有更多数据
	Currency  Currency           `json:"currency"`  // 显示币种
}

// Position 当前持仓信息
type Position struct {
	InstType         string    `json:"instType"`         // 产品类型
	InstId           string    `json:"instId"`           // 产品ID
	MgnMode          string    `json:"mgnMode"`          // 保证金模式
	PosId            string    `json:"posId"`            // 持仓ID
	PosSide          string    `json:"posSide"`          // 持仓方向
	Pos              string    `json:"pos"`              // 持仓数量
	BaseBal          string    `json:"baseBal"`          // 交易币余额
	QuoteBal         string    `json:"quoteBal"`         // 计价币余额
	BaseBorrowed     string    `json:"baseBorrowed"`     // 交易币已借
	BaseInterest     string    `json:"baseInterest"`     // 交易币计息
	QuoteBorrowed    string    `json:"quoteBorrowed"`    // 计价币已借
	QuoteInterest    string    `json:"quoteInterest"`    // 计价币计息
	PosCcy           string    `json:"posCcy"`           // 仓位资产币种
	AvailPos         string    `json:"availPos"`         // 可平仓数量
	AvgPx            string    `json:"avgPx"`            // 开仓均价
	NonSettleAvgPx   string    `json:"nonSettleAvgPx"`   // 未结算均价
	Upl              string    `json:"upl"`              // 未实现收益
	UplRatio         string    `json:"uplRatio"`         // 未实现收益率
	UplLastPx        string    `json:"uplLastPx"`        // 以最新成交价计算的未实现收益
	UplRatioLastPx   string    `json:"uplRatioLastPx"`   // 以最新成交价计算的未实现收益率
	Lever            string    `json:"lever"`            // 杠杆倍数
	LiqPx            string    `json:"liqPx"`            // 预估强平价
	MarkPx           string    `json:"markPx"`           // 最新标记价格
	Imr              string    `json:"imr"`              // 初始保证金
	Margin           string    `json:"margin"`           // 保证金余额
	MgnRatio         string    `json:"mgnRatio"`         // 维持保证金率
	Mmr              string    `json:"mmr"`              // 维持保证金
	Liab             string    `json:"liab"`             // 负债额
	LiabCcy          string    `json:"liabCcy"`          // 负债币种
	Interest         string    `json:"interest"`         // 利息
	TradeId          string    `json:"tradeId"`          // 最新成交ID
	OptVal           string    `json:"optVal"`           // 期权市值
	PendingCloseOrdLiabVal string `json:"pendingCloseOrdLiabVal"` // 逐仓杠杆负债对应平仓挂单的数量
	NotionalUsd      string    `json:"notionalUsd"`      // 以美金价值为单位的持仓数量
	Adl              string    `json:"adl"`              // 信号区
	Ccy              string    `json:"ccy"`              // 占用保证金的币种
	Last             string    `json:"last"`             // 最新成交价
	IdxPx            string    `json:"idxPx"`            // 最新指数价格
	UsdPx            string    `json:"usdPx"`            // 保证金币种的市场最新美金价格
	BePx             string    `json:"bePx"`             // 盈亏平衡价
	DeltaBS          string    `json:"deltaBS"`          // 美金本位持仓仓位delta
	DeltaPA          string    `json:"deltaPA"`          // 币本位持仓仓位delta
	GammaBS          string    `json:"gammaBS"`          // 美金本位持仓仓位gamma
	GammaPA          string    `json:"gammaPA"`          // 币本位持仓仓位gamma
	ThetaBS          string    `json:"thetaBS"`          // 美金本位持仓仓位theta
	ThetaPA          string    `json:"thetaPA"`          // 币本位持仓仓位theta
	VegaBS           string    `json:"vegaBS"`           // 美金本位持仓仓位vega
	VegaPA           string    `json:"vegaPA"`           // 币本位持仓仓位vega
	SpotInUseAmt     string    `json:"spotInUseAmt"`     // 现货对冲占用数量
	SpotInUseCcy     string    `json:"spotInUseCcy"`     // 现货对冲占用币种
	ClSpotInUseAmt   string    `json:"clSpotInUseAmt"`   // 用户自定义现货占用数量
	MaxSpotInUseAmt  string    `json:"maxSpotInUseAmt"`  // 系统计算得到的最大可能现货占用数量
	RealizedPnl      string    `json:"realizedPnl"`      // 已实现收益
	SettledPnl       string    `json:"settledPnl"`       // 已结算收益
	Pnl              string    `json:"pnl"`              // 平仓订单累计收益额
	Fee              string    `json:"fee"`              // 累计手续费金额
	FundingFee       string    `json:"fundingFee"`       // 累计资金费用
	LiqPenalty       string    `json:"liqPenalty"`       // 累计爆仓罚金
	CloseOrderAlgo   []CloseOrderAlgo `json:"closeOrderAlgo"` // 平仓策略委托订单
	CTime            string    `json:"cTime"`            // 持仓创建时间
	UTime            string    `json:"uTime"`            // 最近一次持仓更新时间
	BizRefId         string    `json:"bizRefId"`         // 外部业务id
	BizRefType       string    `json:"bizRefType"`       // 外部业务类型
	Currency         Currency  `json:"currency"`         // 显示币种
	CreateTime       time.Time `json:"createTime"`       // 创建时间
	UpdateTime       time.Time `json:"updateTime"`       // 更新时间
}

// CloseOrderAlgo 平仓策略委托订单
type CloseOrderAlgo struct {
	AlgoId         string `json:"algoId"`         // 策略委托单ID
	SlTriggerPx    string `json:"slTriggerPx"`    // 止损触发价
	SlTriggerPxType string `json:"slTriggerPxType"` // 止损触发价类型
	TpTriggerPx    string `json:"tpTriggerPx"`    // 止盈委托价
	TpTriggerPxType string `json:"tpTriggerPxType"` // 止盈触发价类型
	CloseFraction  string `json:"closeFraction"`  // 策略委托触发时，平仓的百分比
}

// PositionsRequest 持仓查询请求
type PositionsRequest struct {
	InstType string `json:"instType,omitempty" form:"instType"` // 产品类型
	InstId   string `json:"instId,omitempty" form:"instId"`     // 交易产品ID
	PosId    string `json:"posId,omitempty" form:"posId"`       // 持仓ID
}

// PositionsResponse 持仓查询响应
type PositionsResponse struct {
	Positions []*Position `json:"positions"` // 持仓列表
	Currency  Currency    `json:"currency"`  // 显示币种
}