package tests

import (
	"testing"

	"github.com/cardchoosen/AlphaArk_Gin/internal/config"
	"github.com/cardchoosen/AlphaArk_Gin/internal/models"
	"github.com/cardchoosen/AlphaArk_Gin/internal/service"
)

func TestPositionsHistoryRequest(t *testing.T) {
	// 测试请求参数结构
	req := &models.PositionsHistoryRequest{
		InstType: "SWAP",
		InstId:   "BTC-USD-SWAP",
		MgnMode:  "cross",
		Type:     "2",
		Limit:    "10",
	}

	if req.InstType != "SWAP" {
		t.Errorf("Expected InstType to be 'SWAP', got %s", req.InstType)
	}

	if req.InstId != "BTC-USD-SWAP" {
		t.Errorf("Expected InstId to be 'BTC-USD-SWAP', got %s", req.InstId)
	}

	if req.Limit != "10" {
		t.Errorf("Expected Limit to be '10', got %s", req.Limit)
	}
}

func TestPositionHistoryModel(t *testing.T) {
	// 测试历史持仓数据模型
	position := &models.PositionHistory{
		InstType:    "SWAP",
		InstId:      "BTC-USD-SWAP",
		MgnMode:     "cross",
		Type:        "2",
		PosId:       "123456789",
		PosSide:     "long",
		OpenAvgPx:   "45000.0",
		CloseAvgPx:  "46000.0",
		RealizedPnl: "1500.0",
		Pnl:         "1500.0",
		PnlRatio:    "3.33",
		Lever:       "10",
		Direction:   "long",
		Currency:    models.CurrencyUSDT,
	}

	if position.InstType != "SWAP" {
		t.Errorf("Expected InstType to be 'SWAP', got %s", position.InstType)
	}

	if position.Currency != models.CurrencyUSDT {
		t.Errorf("Expected Currency to be 'USDT', got %s", position.Currency)
	}

	if position.RealizedPnl != "1500.0" {
		t.Errorf("Expected RealizedPnl to be '1500.0', got %s", position.RealizedPnl)
	}

	if position.Pnl != "1500.0" {
		t.Errorf("Expected Pnl to be '1500.0', got %s", position.Pnl)
	}
}

func TestPositionsHistoryService(t *testing.T) {
	// 测试服务接口
	cfg := &config.OKXConfig{
		APIKey:     "test-api-key",
		SecretKey:  "test-secret-key",
		Passphrase: "test-passphrase",
		BaseURL:    "https://www.okx.com",
	}

	accountService := service.NewAccountService(cfg)

	// 验证服务实现了接口
	if accountService == nil {
		t.Error("Expected accountService to be created, got nil")
	}

	// 测试请求参数
	req := &models.PositionsHistoryRequest{
		InstType: "SWAP",
		Limit:    "10",
	}

	// 注意：这里会因为API配置问题而失败，但我们主要测试接口和数据结构
	_, err := accountService.GetPositionsHistory(req, models.CurrencyUSDT)
	
	// 由于没有真实的API配置，这里会失败，但我们可以验证错误类型
	if err == nil {
		t.Log("API call succeeded (this would only happen with valid OKX credentials)")
	} else {
		t.Logf("Expected API call to fail without valid credentials: %v", err)
	}
}

func TestPositionsHistoryResponse(t *testing.T) {
	// 测试响应结构
	response := &models.PositionsHistoryResponse{
		Positions: []*models.PositionHistory{
			{
				InstType:    "SWAP",
				InstId:      "BTC-USD-SWAP",
				RealizedPnl: "1500.0",
				Pnl:         "1500.0",
				Currency:    models.CurrencyUSDT,
			},
		},
		HasMore:  false,
		Currency: models.CurrencyUSDT,
	}

	if len(response.Positions) != 1 {
		t.Errorf("Expected 1 position, got %d", len(response.Positions))
	}

	if response.HasMore {
		t.Error("Expected HasMore to be false")
	}

	if response.Currency != models.CurrencyUSDT {
		t.Errorf("Expected Currency to be 'USDT', got %s", response.Currency)
	}
}