# OKX API 集成文档

## 概述

本项目已集成OKX交易所的V5 API，支持获取交易对信息、市场数据等功能。

## API配置

### 环境变量配置

在 `.env` 文件中配置以下环境变量：

```bash
# OKX API配置
OKX_API_KEY=XXXXX
OKX_SECRET_KEY=XXXXX
OKX_PASSPHRASE=XXXXX
OKX_IP=
OKX_REMARK=Gin项目
OKX_PERMISSIONS=读取/提现/交易
OKX_BASE_URL=https://www.okx.com
OKX_IS_TEST=false
```

### 配置说明

- `OKX_API_KEY`: API密钥
- `OKX_SECRET_KEY`: 密钥签名
- `OKX_PASSPHRASE`: API密码
- `OKX_IP`: IP白名单（可选）
- `OKX_REMARK`: 备注名称
- `OKX_PERMISSIONS`: API权限
- `OKX_BASE_URL`: API基础URL
- `OKX_IS_TEST`: 是否使用测试环境

## API端点

### 1. 获取交易对信息

**GET** `/api/v1/okx/instruments`

获取所有交易对信息，默认返回SPOT类型。

**查询参数：**
- `instType` (可选): 交易对类型，支持 SPOT, MARGIN, SWAP, FUTURES, OPTION

**示例：**
```bash
# 获取所有SPOT交易对
curl http://localhost:8080/api/v1/okx/instruments

# 获取FUTURES交易对
curl http://localhost:8080/api/v1/okx/instruments?instType=FUTURES
```

**响应示例：**
```json
{
  "success": true,
  "message": "获取交易对信息成功",
  "data": {
    "code": "0",
    "msg": "",
    "data": [
      {
        "instType": "SPOT",
        "instId": "BTC-USDT",
        "baseCcy": "BTC",
        "quoteCcy": "USDT",
        "tickSz": "0.1",
        "lotSz": "0.00000001",
        "minSz": "0.00000001",
        "maxSz": "100000000",
        "state": "live"
      }
    ]
  }
}
```

### 2. 根据类型获取交易对

**GET** `/api/v1/okx/instruments/:type`

获取指定类型的交易对信息。

**路径参数：**
- `type`: 交易对类型 (SPOT, MARGIN, SWAP, FUTURES, OPTION)

**示例：**
```bash
# 获取SPOT交易对
curl http://localhost:8080/api/v1/okx/instruments/SPOT

# 获取FUTURES交易对
curl http://localhost:8080/api/v1/okx/instruments/FUTURES
```

### 3. 获取API配置信息

**GET** `/api/v1/okx/config`

获取OKX API配置信息（仅显示非敏感信息）。

**示例：**
```bash
curl http://localhost:8080/api/v1/okx/config
```

**响应示例：**
```json
{
  "success": true,
  "message": "获取OKX配置信息成功",
  "data": {
    "remark": "Gin项目",
    "permissions": "读取/提现/交易",
    "baseUrl": "https://www.okx.com",
    "isTest": false,
    "hasApiKey": true,
    "hasSecretKey": true,
    "hasPassphrase": true
  }
}
```

## 交易对类型说明

- **SPOT**: 现货交易
- **MARGIN**: 杠杆交易
- **SWAP**: 永续合约
- **FUTURES**: 交割合约
- **OPTION**: 期权交易

## 错误处理

API会返回标准的错误响应格式：

```json
{
  "success": false,
  "error": "错误信息"
}
```

常见错误：
- `400`: 请求参数错误
- `500`: 服务器内部错误或API调用失败

## 安全注意事项

1. **API密钥安全**: 请妥善保管API密钥，不要提交到版本控制系统
2. **IP白名单**: 建议设置IP白名单以提高安全性
3. **权限控制**: 根据实际需要设置API权限
4. **测试环境**: 开发时建议使用测试环境

## 扩展功能

### 添加更多API功能

可以在 `internal/api/okx_client.go` 中添加更多API方法：

```go
// 获取K线数据
func (c *OKXClient) GetKlineData(instId, period string) (*KlineResponse, error) {
    // 实现代码
}

// 获取市场深度
func (c *OKXClient) GetOrderBook(instId string) (*OrderBookResponse, error) {
    // 实现代码
}

// 下单（需要私有API）
func (c *OKXClient) PlaceOrder(order OrderRequest) (*OrderResponse, error) {
    // 实现代码
}
```

### 添加缓存机制

为了提高性能，可以添加Redis缓存：

```go
// 缓存交易对信息
func (c *OKXClient) GetInstrumentsWithCache(instType string) (*InstrumentsResponse, error) {
    // 先检查缓存
    // 如果缓存不存在，调用API并缓存结果
}
```

## 测试

运行测试：

```bash
# 运行所有测试
make test

# 运行OKX API测试
go test ./internal/api -v -run TestOKX
```

## 监控和日志

建议添加以下监控指标：

- API调用次数
- API响应时间
- API错误率
- 缓存命中率

## 相关链接

- [OKX API V5 官方文档](https://www.okx.com/docs-v5/zh/)
- [OKX API 状态页面](https://status.okx.com/) 