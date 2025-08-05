package tests

import (
	"testing"

	"github.com/yourname/my-gin-project/internal/api"
	"github.com/yourname/my-gin-project/internal/config"
)

func TestOKXConfig(t *testing.T) {
	// 测试配置加载
	cfg := &config.OKXConfig{
		APIKey:      "test-api-key",
		SecretKey:   "test-secret-key",
		Passphrase:  "test-passphrase",
		Remark:      "测试项目",
		Permissions: "读取",
		BaseURL:     "https://www.okx.com",
		IsTest:      true,
	}

	if cfg.APIKey != "test-api-key" {
		t.Errorf("Expected APIKey to be 'test-api-key', got %s", cfg.APIKey)
	}

	if cfg.Remark != "测试项目" {
		t.Errorf("Expected Remark to be '测试项目', got %s", cfg.Remark)
	}

	if !cfg.IsTest {
		t.Errorf("Expected IsTest to be true, got %v", cfg.IsTest)
	}
}

func TestOKXClientCreation(t *testing.T) {
	cfg := &config.OKXConfig{
		APIKey:     "test-api-key",
		SecretKey:  "test-secret-key",
		Passphrase: "test-passphrase",
		BaseURL:    "https://www.okx.com",
	}

	client := api.NewOKXClient(cfg)
	if client == nil {
		t.Error("Expected client to be created, got nil")
	}

	if client.GetConfig() != cfg {
		t.Error("Expected client config to match provided config")
	}
}

func TestOKXTimestamp(t *testing.T) {
	cfg := &config.OKXConfig{
		APIKey:     "test-api-key",
		SecretKey:  "test-secret-key",
		Passphrase: "test-passphrase",
	}

	client := api.NewOKXClient(cfg)
	timestamp := client.Timestamp()

	if timestamp == "" {
		t.Error("Expected timestamp to be non-empty")
	}
}

func TestOKXSign(t *testing.T) {
	cfg := &config.OKXConfig{
		APIKey:     "test-api-key",
		SecretKey:  "test-secret-key",
		Passphrase: "test-passphrase",
	}

	client := api.NewOKXClient(cfg)

	timestamp := "1234567890"
	method := "GET"
	requestPath := "/api/v5/public/instruments"
	body := ""

	sign := client.Sign(timestamp, method, requestPath, body)

	if sign == "" {
		t.Error("Expected signature to be non-empty")
	}
}

func TestOKXHeaders(t *testing.T) {
	cfg := &config.OKXConfig{
		APIKey:     "test-api-key",
		SecretKey:  "test-secret-key",
		Passphrase: "test-passphrase",
	}

	client := api.NewOKXClient(cfg)

	method := "GET"
	requestPath := "/api/v5/public/instruments"
	body := ""

	headers := client.GenerateHeaders(method, requestPath, body)

	requiredHeaders := []string{
		"OK-ACCESS-KEY",
		"OK-ACCESS-SIGN",
		"OK-ACCESS-TIMESTAMP",
		"OK-ACCESS-PASSPHRASE",
		"Content-Type",
	}

	for _, header := range requiredHeaders {
		if _, exists := headers[header]; !exists {
			t.Errorf("Expected header %s to exist", header)
		}
	}

	if headers["OK-ACCESS-KEY"] != cfg.APIKey {
		t.Errorf("Expected OK-ACCESS-KEY to be %s, got %s", cfg.APIKey, headers["OK-ACCESS-KEY"])
	}

	if headers["OK-ACCESS-PASSPHRASE"] != cfg.Passphrase {
		t.Errorf("Expected OK-ACCESS-PASSPHRASE to be %s, got %s", cfg.Passphrase, headers["OK-ACCESS-PASSPHRASE"])
	}
}
