package api

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/cardchoosen/AlphaArk_Gin/internal/config"
)

// OKXClient OKX API客户端
type OKXClient struct {
	config     *config.OKXConfig
	client     *http.Client
	timeOffset int64     // 与OKX服务器的时间偏移量（毫秒）
	lastSync   time.Time // 上次同步时间
}

// NewOKXClient 创建OKX客户端
func NewOKXClient(cfg *config.OKXConfig) *OKXClient {
	client := &OKXClient{
		config: cfg,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
		timeOffset: 0,
		lastSync:   time.Time{},
	}

	// 初始化时同步时间
	client.SyncTime()

	return client
}

// Instrument 交易对信息
type Instrument struct {
	InstType   string `json:"instType"`
	InstID     string `json:"instId"`
	InstFamily string `json:"instFamily"`
	BaseCcy    string `json:"baseCcy"`
	QuoteCcy   string `json:"quoteCcy"`
	SettleCcy  string `json:"settleCcy"`
	CtVal      string `json:"ctVal"`
	CtMult     string `json:"ctMult"`
	CtValCcy   string `json:"ctValCcy"`
	OptType    string `json:"optType"`
	Stk        string `json:"stk"`
	ListTime   string `json:"listTime"`
	ExpTime    string `json:"expTime"`
	TickSz     string `json:"tickSz"`
	LotSz      string `json:"lotSz"`
	MinSz      string `json:"minSz"`
	MaxSz      string `json:"maxSz"`
	Uly        string `json:"uly"`
	Category   string `json:"category"`
	State      string `json:"state"`
}

// InstrumentsResponse 获取交易对响应
type InstrumentsResponse struct {
	Code string       `json:"code"`
	Msg  string       `json:"msg"`
	Data []Instrument `json:"data"`
}

// Ticker 行情数据
type Ticker struct {
	InstType  string `json:"instType"`
	InstId    string `json:"instId"`
	Last      string `json:"last"`
	LastSz    string `json:"lastSz"`
	AskPx     string `json:"askPx"`
	AskSz     string `json:"askSz"`
	BidPx     string `json:"bidPx"`
	BidSz     string `json:"bidSz"`
	Open24h   string `json:"open24h"`
	High24h   string `json:"high24h"`
	Low24h    string `json:"low24h"`
	Vol24h    string `json:"vol24h"`
	VolCcy24h string `json:"volCcy24h"`
	SodUtc0   string `json:"sodUtc0"`
	SodUtc8   string `json:"sodUtc8"`
	Ts        string `json:"ts"`
}

// TickerResponse 获取行情响应
type TickerResponse struct {
	Code string   `json:"code"`
	Msg  string   `json:"msg"`
	Data []Ticker `json:"data"`
}

// SystemTimeResponse 系统时间响应
type SystemTimeResponse struct {
	Code string `json:"code"`
	Msg  string `json:"msg"`
	Data []struct {
		Ts string `json:"ts"`
	} `json:"data"`
}

// Position 持仓信息
type Position struct {
	InstType               string        `json:"instType"`
	InstId                 string        `json:"instId"`
	MgnMode                string        `json:"mgnMode"`
	PosId                  string        `json:"posId"`
	PosSide                string        `json:"posSide"`
	Pos                    string        `json:"pos"`
	BaseBal                string        `json:"baseBal"`
	QuoteBal               string        `json:"quoteBal"`
	BaseBorrowed           string        `json:"baseBorrowed"`
	BaseInterest           string        `json:"baseInterest"`
	QuoteBorrowed          string        `json:"quoteBorrowed"`
	QuoteInterest          string        `json:"quoteInterest"`
	PosCcy                 string        `json:"posCcy"`
	AvailPos               string        `json:"availPos"`
	AvgPx                  string        `json:"avgPx"`
	NonSettleAvgPx         string        `json:"nonSettleAvgPx"`
	Upl                    string        `json:"upl"`
	UplRatio               string        `json:"uplRatio"`
	UplLastPx              string        `json:"uplLastPx"`
	UplRatioLastPx         string        `json:"uplRatioLastPx"`
	Lever                  string        `json:"lever"`
	LiqPx                  string        `json:"liqPx"`
	MarkPx                 string        `json:"markPx"`
	Imr                    string        `json:"imr"`
	Margin                 string        `json:"margin"`
	MgnRatio               string        `json:"mgnRatio"`
	Mmr                    string        `json:"mmr"`
	Liab                   string        `json:"liab"`
	LiabCcy                string        `json:"liabCcy"`
	Interest               string        `json:"interest"`
	TradeId                string        `json:"tradeId"`
	OptVal                 string        `json:"optVal"`
	PendingCloseOrdLiabVal string        `json:"pendingCloseOrdLiabVal"`
	NotionalUsd            string        `json:"notionalUsd"`
	Adl                    string        `json:"adl"`
	Ccy                    string        `json:"ccy"`
	Last                   string        `json:"last"`
	IdxPx                  string        `json:"idxPx"`
	UsdPx                  string        `json:"usdPx"`
	BePx                   string        `json:"bePx"`
	DeltaBS                string        `json:"deltaBS"`
	DeltaPA                string        `json:"deltaPA"`
	GammaBS                string        `json:"gammaBS"`
	GammaPA                string        `json:"gammaPA"`
	ThetaBS                string        `json:"thetaBS"`
	ThetaPA                string        `json:"thetaPA"`
	VegaBS                 string        `json:"vegaBS"`
	VegaPA                 string        `json:"vegaPA"`
	SpotInUseAmt           string        `json:"spotInUseAmt"`
	SpotInUseCcy           string        `json:"spotInUseCcy"`
	ClSpotInUseAmt         string        `json:"clSpotInUseAmt"`
	MaxSpotInUseAmt        string        `json:"maxSpotInUseAmt"`
	RealizedPnl            string        `json:"realizedPnl"`
	SettledPnl             string        `json:"settledPnl"`
	Pnl                    string        `json:"pnl"`
	Fee                    string        `json:"fee"`
	FundingFee             string        `json:"fundingFee"`
	LiqPenalty             string        `json:"liqPenalty"`
	CloseOrderAlgo         []interface{} `json:"closeOrderAlgo"`
	CTime                  string        `json:"cTime"`
	UTime                  string        `json:"uTime"`
	BizRefId               string        `json:"bizRefId"`
	BizRefType             string        `json:"bizRefType"`
}

// PositionsResponse 持仓响应
type PositionsResponse struct {
	Code string     `json:"code"`
	Msg  string     `json:"msg"`
	Data []Position `json:"data"`
}

// GetInstruments 获取交易对信息
func (c *OKXClient) GetInstruments(instType string) (*InstrumentsResponse, error) {
	url := fmt.Sprintf("%s/api/v5/public/instruments?instType=%s", c.config.BaseURL, instType)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求失败: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}

	var result InstrumentsResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}

	return &result, nil
}

// GetTicker 获取行情数据
func (c *OKXClient) GetTicker(instId string) (*TickerResponse, error) {
	url := fmt.Sprintf("%s/api/v5/market/ticker?instId=%s", c.config.BaseURL, instId)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求失败: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}

	var result TickerResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}

	return &result, nil
}

// Sign 签名方法（用于私有API）
func (c *OKXClient) Sign(timestamp, method, requestPath, body string) string {
	message := timestamp + method + requestPath + body
	mac := hmac.New(sha256.New, []byte(c.config.SecretKey))
	mac.Write([]byte(message))
	return base64.StdEncoding.EncodeToString(mac.Sum(nil))
}

// SyncTime 同步时间
func (c *OKXClient) SyncTime() error {
	// 如果距离上次同步时间不足5分钟，跳过同步
	if !c.lastSync.IsZero() && time.Since(c.lastSync) < 5*time.Minute {
		return nil
	}

	// 获取OKX服务器时间
	serverTime, err := c.GetSystemTime()
	if err != nil {
		return fmt.Errorf("获取服务器时间失败: %w", err)
	}

	if len(serverTime.Data) == 0 {
		return fmt.Errorf("服务器时间数据为空")
	}

	// 解析服务器时间戳（毫秒）
	serverTs, err := strconv.ParseInt(serverTime.Data[0].Ts, 10, 64)
	if err != nil {
		return fmt.Errorf("解析服务器时间戳失败: %w", err)
	}

	// 计算时间偏移量（服务器时间 - 本地时间）
	localTs := time.Now().UnixMilli()
	c.timeOffset = serverTs - localTs
	c.lastSync = time.Now()

	return nil
}

// Timestamp 生成时间戳
func (c *OKXClient) Timestamp() string {
	// 确保时间同步
	c.SyncTime()

	// 使用同步后的时间生成ISO 8601格式的时间戳
	syncedTime := time.Now().Add(time.Duration(c.timeOffset) * time.Millisecond)
	return syncedTime.UTC().Format("2006-01-02T15:04:05.000Z")
}

// GenerateHeaders 生成签名头（用于私有API）
func (c *OKXClient) GenerateHeaders(method, requestPath, body string) map[string]string {
	timestamp := c.Timestamp()
	sign := c.Sign(timestamp, method, requestPath, body)

	return map[string]string{
		"OK-ACCESS-KEY":        c.config.APIKey,
		"OK-ACCESS-SIGN":       sign,
		"OK-ACCESS-TIMESTAMP":  timestamp,
		"OK-ACCESS-PASSPHRASE": c.config.Passphrase,
		"Content-Type":         "application/json",
	}
}

// GetSystemTime 获取系统时间
func (c *OKXClient) GetSystemTime() (*SystemTimeResponse, error) {
	url := fmt.Sprintf("%s/api/v5/public/time", c.config.BaseURL)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求失败: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}

	var result SystemTimeResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}

	return &result, nil
}

// GetPositions 获取当前持仓信息
func (c *OKXClient) GetPositions(instType, instId, posId, currency string) (*PositionsResponse, error) {
	url := fmt.Sprintf("%s/api/v5/account/positions", c.config.BaseURL)

	// 构建查询参数
	params := make([]string, 0)
	if instType != "" {
		params = append(params, fmt.Sprintf("instType=%s", instType))
	}
	if instId != "" {
		params = append(params, fmt.Sprintf("instId=%s", instId))
	}
	if posId != "" {
		params = append(params, fmt.Sprintf("posId=%s", posId))
	}
	if currency != "" {
		params = append(params, fmt.Sprintf("currency=%s", currency))
	}

	if len(params) > 0 {
		url += "?" + strings.Join(params, "&")
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}

	// 构建请求路径（用于签名）
	requestPath := "/api/v5/account/positions"
	if len(params) > 0 {
		requestPath += "?" + strings.Join(params, "&")
	}

	// 添加签名头（私有API需要）
	headers := c.GenerateHeaders("GET", requestPath, "")
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求失败: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}

	var result PositionsResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}

	return &result, nil
}

// GetConfig 获取配置（用于测试）
func (c *OKXClient) GetConfig() *config.OKXConfig {
	return c.config
}
