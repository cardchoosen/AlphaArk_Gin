package config

import (
	"os"
	"strconv"
)

// Config 应用配置结构
type Config struct {
	Environment string
	Port        string
	DatabaseURL string
	JWTSecret   string
	OKX         OKXConfig
}

// OKXConfig OKX API配置
type OKXConfig struct {
	APIKey      string
	SecretKey   string
	Passphrase  string
	IP          string
	Remark      string
	Permissions string
	BaseURL     string
	IsTest      bool
}

// Load 加载配置
func Load() *Config {
	return &Config{
		Environment: getEnv("ENVIRONMENT", "development"),
		Port:        getEnv("PORT", "8080"),
		DatabaseURL: getEnv("DATABASE_URL", ""),
		JWTSecret:   getEnv("JWT_SECRET", "your-secret-key"),
		OKX: OKXConfig{
			APIKey:      getEnv("OKX_API_KEY", ""),
			SecretKey:   getEnv("OKX_SECRET_KEY", ""),
			Passphrase:  getEnv("OKX_PASSPHRASE", ""),
			IP:          getEnv("OKX_IP", ""),
			Remark:      getEnv("OKX_REMARK", "Gin项目"),
			Permissions: getEnv("OKX_PERMISSIONS", "读取/提现/交易"),
			BaseURL:     getEnv("OKX_BASE_URL", "https://www.okx.com"),
			IsTest:      getEnvBool("OKX_IS_TEST", false),
		},
	}
}

// getEnv 获取环境变量，如果不存在则返回默认值
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvInt 获取整数环境变量
func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

// getEnvBool 获取布尔环境变量
func getEnvBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}
