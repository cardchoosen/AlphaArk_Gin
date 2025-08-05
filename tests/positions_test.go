package tests

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/cardchoosen/AlphaArk_Gin/internal/api"
	"github.com/cardchoosen/AlphaArk_Gin/internal/config"
	"github.com/cardchoosen/AlphaArk_Gin/internal/models"
	"github.com/cardchoosen/AlphaArk_Gin/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestGetPositions 测试获取当前持仓信息
func TestGetPositions(t *testing.T) {
	// 设置测试环境
	gin.SetMode(gin.TestMode)
	
	// 创建测试配置
	cfg := &config.Config{
		OKX: config.OKXConfig{
			APIKey:     "test-api-key",
			SecretKey:  "test-secret-key",
			Passphrase: "test-passphrase",
			BaseURL:    "https://www.okx.com",
			IsTest:     true,
		},
	}

	// 创建路由
	r := gin.New()
	api.SetupAccountRoutes(r, cfg)

	tests := []struct {
		name           string
		queryParams    map[string]string
		expectedStatus int
		checkResponse  bool
	}{
		{
			name:           "获取所有持仓",
			queryParams:    map[string]string{},
			expectedStatus: http.StatusOK,
			checkResponse:  true,
		},
		{
			name: "按产品类型筛选",
			queryParams: map[string]string{
				"instType": "SWAP",
			},
			expectedStatus: http.StatusOK,
			checkResponse:  true,
		},
		{
			name: "按产品ID筛选",
			queryParams: map[string]string{
				"instId": "BTC-USDT-SWAP",
			},
			expectedStatus: http.StatusOK,
			checkResponse:  true,
		},
		{
			name: "指定币种",
			queryParams: map[string]string{
				"currency": "USD",
			},
			expectedStatus: http.StatusOK,
			checkResponse:  true,
		},
		{
			name: "不支持的币种",
			queryParams: map[string]string{
				"currency": "INVALID",
			},
			expectedStatus: http.StatusBadRequest,
			checkResponse:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 构建查询参数
			queryString := ""
			if len(tt.queryParams) > 0 {
				params := url.Values{}
				for key, value := range tt.queryParams {
					params.Add(key, value)
				}
				queryString = "?" + params.Encode()
			}

			// 创建请求
			req, err := http.NewRequest("GET", "/api/v1/account/positions"+queryString, nil)
			require.NoError(t, err)

			// 创建响应记录器
			w := httptest.NewRecorder()

			// 执行请求
			r.ServeHTTP(w, req)

			// 验证状态码
			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.checkResponse && w.Code == http.StatusOK {
				// 验证响应格式
				var response struct {
					Success bool                        `json:"success"`
					Message string                      `json:"message"`
					Data    *models.PositionsResponse   `json:"data"`
				}

				err := json.Unmarshal(w.Body.Bytes(), &response)
				require.NoError(t, err)

				assert.True(t, response.Success)
				assert.Equal(t, "获取当前持仓信息成功", response.Message)
				assert.NotNil(t, response.Data)
				assert.NotNil(t, response.Data.Positions)
			}
		})
	}
}

// TestGetPositionsHistory 测试获取历史持仓信息
func TestGetPositionsHistory(t *testing.T) {
	// 设置测试环境
	gin.SetMode(gin.TestMode)
	
	// 创建测试配置
	cfg := &config.Config{
		OKX: config.OKXConfig{
			APIKey:     "test-api-key",
			SecretKey:  "test-secret-key",
			Passphrase: "test-passphrase",
			BaseURL:    "https://www.okx.com",
			IsTest:     true,
		},
	}

	// 创建路由
	r := gin.New()
	api.SetupAccountRoutes(r, cfg)

	tests := []struct {
		name           string
		queryParams    map[string]string
		expectedStatus int
		checkResponse  bool
	}{
		{
			name:           "获取所有历史持仓",
			queryParams:    map[string]string{},
			expectedStatus: http.StatusOK,
			checkResponse:  true,
		},
		{
			name: "按产品类型筛选",
			queryParams: map[string]string{
				"instType": "SWAP",
			},
			expectedStatus: http.StatusOK,
			checkResponse:  true,
		},
		{
			name: "按平仓类型筛选",
			queryParams: map[string]string{
				"type": "2",
			},
			expectedStatus: http.StatusOK,
			checkResponse:  true,
		},
		{
			name: "分页限制",
			queryParams: map[string]string{
				"limit": "10",
			},
			expectedStatus: http.StatusOK,
			checkResponse:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 构建查询参数
			queryString := ""
			if len(tt.queryParams) > 0 {
				params := url.Values{}
				for key, value := range tt.queryParams {
					params.Add(key, value)
				}
				queryString = "?" + params.Encode()
			}

			// 创建请求
			req, err := http.NewRequest("GET", "/api/v1/account/positions-history"+queryString, nil)
			require.NoError(t, err)

			// 创建响应记录器
			w := httptest.NewRecorder()

			// 执行请求
			r.ServeHTTP(w, req)

			// 验证状态码
			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.checkResponse && w.Code == http.StatusOK {
				// 验证响应格式
				var response struct {
					Success bool                                `json:"success"`
					Message string                              `json:"message"`
					Data    *models.PositionsHistoryResponse    `json:"data"`
				}

				err := json.Unmarshal(w.Body.Bytes(), &response)
				require.NoError(t, err)

				assert.True(t, response.Success)
				assert.Equal(t, "获取历史持仓信息成功", response.Message)
				assert.NotNil(t, response.Data)
				assert.NotNil(t, response.Data.Positions)
			}
		})
	}
}

// TestPositionModels 测试持仓相关模型
func TestPositionModels(t *testing.T) {
	// 测试当前持仓模型
	position := &models.Position{
		InstType: "SWAP",
		InstId:   "BTC-USDT-SWAP",
		MgnMode:  "cross",
		PosId:    "123456789",
		PosSide:  "long",
		Pos:      "1.5",
		Currency: models.CurrencyUSDT,
	}

	assert.Equal(t, "SWAP", position.InstType)
	assert.Equal(t, "BTC-USDT-SWAP", position.InstId)
	assert.Equal(t, models.CurrencyUSDT, position.Currency)

	// 测试历史持仓模型
	historyPosition := &models.PositionHistory{
		InstType:      "SWAP",
		InstId:        "BTC-USDT-SWAP",
		MgnMode:       "cross",
		Type:          "2",
		PosId:         "123456789",
		OpenAvgPx:     "45000.0",
		CloseAvgPx:    "46000.0",
		RealizedPnl:   "1500.0",
		Currency:      models.CurrencyUSDT,
	}

	assert.Equal(t, "SWAP", historyPosition.InstType)
	assert.Equal(t, "2", historyPosition.Type)
	assert.Equal(t, "1500.0", historyPosition.RealizedPnl)

	// 测试请求模型
	positionsReq := &models.PositionsRequest{
		InstType: "SWAP",
		InstId:   "BTC-USDT-SWAP",
	}

	assert.Equal(t, "SWAP", positionsReq.InstType)
	assert.Equal(t, "BTC-USDT-SWAP", positionsReq.InstId)

	historyReq := &models.PositionsHistoryRequest{
		InstType: "SWAP",
		Type:     "2",
		Limit:    "10",
	}

	assert.Equal(t, "SWAP", historyReq.InstType)
	assert.Equal(t, "2", historyReq.Type)
	assert.Equal(t, "10", historyReq.Limit)
}

// TestAccountServiceInterface 测试账户服务接口
func TestAccountServiceInterface(t *testing.T) {
	cfg := &config.OKXConfig{
		APIKey:     "test-api-key",
		SecretKey:  "test-secret-key",
		Passphrase: "test-passphrase",
		BaseURL:    "https://www.okx.com",
		IsTest:     true,
	}

	accountService := service.NewAccountService(cfg)

	// 验证接口实现
	assert.NotNil(t, accountService)

	// 测试默认币种
	defaultCurrency := accountService.GetDefaultCurrency()
	assert.Equal(t, models.CurrencyUSDT, defaultCurrency)

	// 测试设置币种
	err := accountService.SetDefaultCurrency(models.CurrencyUSD)
	assert.NoError(t, err)

	newDefaultCurrency := accountService.GetDefaultCurrency()
	assert.Equal(t, models.CurrencyUSD, newDefaultCurrency)
}

// BenchmarkGetPositions 性能测试
func BenchmarkGetPositions(b *testing.B) {
	gin.SetMode(gin.TestMode)
	
	cfg := &config.Config{
		OKX: config.OKXConfig{
			APIKey:     "test-api-key",
			SecretKey:  "test-secret-key",
			Passphrase: "test-passphrase",
			BaseURL:    "https://www.okx.com",
			IsTest:     true,
		},
	}

	r := gin.New()
	api.SetupAccountRoutes(r, cfg)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req, _ := http.NewRequest("GET", "/api/v1/account/positions", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
	}
}

// TestPositionResponseSerialization 测试响应序列化
func TestPositionResponseSerialization(t *testing.T) {
	response := &models.PositionsResponse{
		Positions: []*models.Position{
			{
				InstType: "SWAP",
				InstId:   "BTC-USDT-SWAP",
				PosId:    "123456789",
				Pos:      "1.5",
				Currency: models.CurrencyUSDT,
			},
		},
		Currency: models.CurrencyUSDT,
	}

	// 序列化测试
	data, err := json.Marshal(response)
	require.NoError(t, err)
	assert.Contains(t, string(data), "BTC-USDT-SWAP")

	// 反序列化测试
	var decoded models.PositionsResponse
	err = json.Unmarshal(data, &decoded)
	require.NoError(t, err)
	assert.Equal(t, models.CurrencyUSDT, decoded.Currency)
	assert.Len(t, decoded.Positions, 1)
	assert.Equal(t, "BTC-USDT-SWAP", decoded.Positions[0].InstId)
}