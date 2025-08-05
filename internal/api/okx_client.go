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
	"time"

	"github.com/yourname/my-gin-project/internal/config"
)

// OKXClient OKX API客户端
type OKXClient struct {
	config *config.OKXConfig
	client *http.Client
}

// NewOKXClient 创建OKX客户端
func NewOKXClient(cfg *config.OKXConfig) *OKXClient {
	return &OKXClient{
		config: cfg,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
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

// Sign 签名方法（用于私有API）
func (c *OKXClient) Sign(timestamp, method, requestPath, body string) string {
	message := timestamp + method + requestPath + body
	mac := hmac.New(sha256.New, []byte(c.config.SecretKey))
	mac.Write([]byte(message))
	return base64.StdEncoding.EncodeToString(mac.Sum(nil))
}

// Timestamp 生成时间戳
func (c *OKXClient) Timestamp() string {
	return strconv.FormatInt(time.Now().Unix(), 10)
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

// GetConfig 获取配置（用于测试）
func (c *OKXClient) GetConfig() *config.OKXConfig {
	return c.config
}
