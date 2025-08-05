package service

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/cardchoosen/AlphaArk_Gin/internal/config"
)

// PriceData 价格数据结构
type PriceData struct {
	Symbol          string  `json:"symbol"`
	Price           string  `json:"price"`
	Change24h       string  `json:"change24h,omitempty"`
	ChangePercent24h string  `json:"changePercent24h,omitempty"`
	Volume24h       string  `json:"volume24h,omitempty"`
	High24h         string  `json:"high24h,omitempty"`
	Low24h          string  `json:"low24h,omitempty"`
	Timestamp       int64   `json:"timestamp"`
}

// PriceService 价格服务接口
type PriceService interface {
	GetPrice(symbol string) (*PriceData, error)
	StartPriceStream(symbol string, callback func(*PriceData))
	StopPriceStream()
}

// priceService 价格服务实现
type priceService struct {
	config   *config.OKXConfig
	stopChan chan struct{}
}

// NewPriceService 创建价格服务实例
func NewPriceService(cfg *config.OKXConfig) PriceService {
	return &priceService{
		config:   cfg,
		stopChan: make(chan struct{}),
	}
}

// Ticker OKX行情数据结构
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

// TickerResponse OKX行情响应
type TickerResponse struct {
	Code string   `json:"code"`
	Msg  string   `json:"msg"`
	Data []Ticker `json:"data"`
}

// GetPrice 获取指定交易对的当前价格
func (s *priceService) GetPrice(symbol string) (*PriceData, error) {
	// 调用OKX API获取ticker数据
	url := fmt.Sprintf("%s/api/v5/market/ticker?instId=%s", s.config.BaseURL, symbol)
	
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

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

	var result TickerResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}

	if len(result.Data) == 0 {
		return nil, fmt.Errorf("未找到交易对 %s 的价格数据", symbol)
	}

	ticker := result.Data[0]
	
	// 计算24小时变化
	price, _ := strconv.ParseFloat(ticker.Last, 64)
	open24h, _ := strconv.ParseFloat(ticker.Open24h, 64)
	change24h := price - open24h
	changePercent := 0.0
	if open24h != 0 {
		changePercent = (change24h / open24h) * 100
	}

	return &PriceData{
		Symbol:           symbol,
		Price:            ticker.Last,
		Change24h:        fmt.Sprintf("%.2f", change24h),
		ChangePercent24h: fmt.Sprintf("%.2f", changePercent),
		Volume24h:        ticker.Vol24h,
		High24h:          ticker.High24h,
		Low24h:           ticker.Low24h,
		Timestamp:        time.Now().Unix(),
	}, nil
}

// StartPriceStream 启动价格数据流
func (s *priceService) StartPriceStream(symbol string, callback func(*PriceData)) {
	go func() {
		ticker := time.NewTicker(5 * time.Second) // 每5秒更新一次
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				priceData, err := s.GetPrice(symbol)
				if err != nil {
					log.Printf("获取价格数据失败: %v", err)
					continue
				}
				callback(priceData)
			case <-s.stopChan:
				return
			}
		}
	}()
}

// StopPriceStream 停止价格数据流
func (s *priceService) StopPriceStream() {
	select {
	case s.stopChan <- struct{}{}:
	default:
	}
}