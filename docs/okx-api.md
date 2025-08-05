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

### 时间戳同步

如果遇到 `50102 Timestamp request expired` 错误，说明本地时间与OKX服务器时间不同步。需要调用系统时间接口进行时间同步：

**GET** `/api/v1/okx/system-time`

获取OKX服务器时间，用于本地时间同步。

**示例：**
```bash
curl http://localhost:8080/api/v1/okx/system-time
```

**响应示例：**
```json
{
  "success": true,
  "message": "获取系统时间成功",
  "data": {
    "ts": "1699689600000"
  }
}
```

**时间同步说明：**
- 请求头中的时间戳必须是UTC+0时区的时间
- 系统时间接口返回的是UTC+8时区的时间戳
- 时间差值需要控制在30秒内才能有效规避时间戳过期问题
- 建议在应用启动时和定期调用此接口进行时间同步

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

## 历史持仓信息 API

### 端点

```bash
GET /api/v1/account/positions-history
```

### 查询参数

| 参数名 | 类型 | 必须 | 描述 |
|--------|------|------|------|
| instType | String | 否 | 产品类型：MARGIN(币币杠杆), SWAP(永续合约), FUTURES(交割合约), OPTION(期权) |
| instId | String | 否 | 交易产品ID，如：BTC-USD-SWAP |
| mgnMode | String | 否 | 保证金模式：cross(全仓), isolated(逐仓) |
| type | String | 否 | 平仓类型：1(部分平仓), 2(完全平仓), 3(强平), 4(强减), 5(ADL自动减仓) |
| posId | String | 否 | 持仓ID |
| after | String | 否 | 查询仓位更新之前的内容，Unix时间戳(毫秒) |
| before | String | 否 | 查询仓位更新之后的内容，Unix时间戳(毫秒) |
| limit | String | 否 | 分页数量，最大100，默认100 |
| currency | String | 否 | 显示币种：CNY, USD, USDT, BTC |

### 请求示例

```bash
# 获取所有历史持仓
curl "http://localhost:8080/api/v1/account/positions-history"

# 获取永续合约的历史持仓
curl "http://localhost:8080/api/v1/account/positions-history?instType=SWAP"

# 获取特定产品的历史持仓，限制10条
curl "http://localhost:8080/api/v1/account/positions-history?instId=BTC-USD-SWAP&limit=10"

# 获取强平记录
curl "http://localhost:8080/api/v1/account/positions-history?type=3"
```

### 响应示例

```json
{
  "success": true,
  "message": "获取历史持仓信息成功",
  "data": {
    "positions": [
      {
        "instType": "SWAP",
        "instId": "BTC-USD-SWAP",
        "mgnMode": "cross",
        "type": "2",
        "cTime": "1699689600000",
        "uTime": "1699776000000",
        "openAvgPx": "45000.0",
        "nonSettleAvgPx": "45000.0",
        "closeAvgPx": "46000.0",
        "posId": "123456789",
        "openMaxPos": "2.0",
        "closeTotalPos": "1.5",
        "realizedPnl": "1500.0",
        "settledPnl": "",
        "pnlRatio": "3.33",
        "fee": "-15.0",
        "fundingFee": "5.2",
        "liqPenalty": "",
        "pnl": "1500.0",
        "posSide": "long",
        "lever": "10",
        "direction": "long",
        "triggerPx": "",
        "uly": "BTC-USD",
        "ccy": "USD",
        "currency": "USDT",
        "updateTime": "2024-11-12T12:00:00Z",
        "createTime": "2024-11-11T12:00:00Z"
      }
    ],
    "hasMore": false,
    "currency": "USDT"
  }
}
```

### 字段说明

| 字段名 | 描述 |
|--------|------|
| instType | 产品类型 |
| instId | 交易产品ID |
| mgnMode | 保证金模式：cross(全仓), isolated(逐仓) |
| type | 平仓类型：1(部分平仓), 2(完全平仓), 3(强平), 4(强减), 5(ADL自动减仓) |
| cTime | 仓位创建时间 |
| uTime | 仓位更新时间 |
| openAvgPx | 开仓均价 |
| nonSettleAvgPx | 未结算均价 |
| closeAvgPx | 平仓均价 |
| posId | 仓位ID |
| openMaxPos | 最大持仓量 |
| closeTotalPos | 累计平仓量 |
| realizedPnl | 已实现收益 |
| settledPnl | 已实现收益(全仓交割) |
| pnlRatio | 已实现收益率 |
| fee | 累计手续费金额 |
| fundingFee | 累计资金费用 |
| liqPenalty | 累计爆仓罚金 |
| pnl | 已实现收益 |
| posSide | 持仓模式方向：long/short/net |
| lever | 杠杆倍数 |
| direction | 持仓方向：long(多), short(空) |
| triggerPx | 触发标记价格 |
| uly | 标的指数 |
| ccy | 占用保证金的币种 |
| hasMore | 是否有更多数据 |

## 当前持仓信息 API

### 端点

```bash
GET /api/v1/account/positions
```

### 查询参数

| 参数名 | 类型 | 必须 | 描述 |
|--------|------|------|------|
| instType | String | 否 | 产品类型：MARGIN(币币杠杆), SWAP(永续合约), FUTURES(交割合约), OPTION(期权) |
| instId | String | 否 | 交易产品ID，如：BTC-USDT-SWAP。支持多个instId查询（不超过10个），半角逗号分隔 |
| posId | String | 否 | 持仓ID。支持多个posId查询（不超过20个） |
| currency | String | 否 | 显示币种：CNY, USD, USDT, BTC |

### 请求示例

```bash
# 获取所有当前持仓
curl "http://localhost:8080/api/v1/account/positions"

# 获取BTC-USDT的持仓信息
curl "http://localhost:8080/api/v1/account/positions?instId=BTC-USDT"

# 获取永续合约的持仓信息
curl "http://localhost:8080/api/v1/account/positions?instType=SWAP"

# 获取多个产品的持仓信息
curl "http://localhost:8080/api/v1/account/positions?instId=BTC-USDT,ETH-USDT"
```

### 响应示例

```json
{
  "success": true,
  "message": "获取当前持仓信息成功",
  "data": {
    "positions": [
      {
        "instType": "MARGIN",
        "instId": "BTC-USDT",
        "mgnMode": "isolated",
        "posId": "1752810569801498626",
        "posSide": "net",
        "pos": "0.00190433573",
        "baseBal": "",
        "quoteBal": "",
        "baseBorrowed": "",
        "baseInterest": "",
        "quoteBorrowed": "",
        "quoteInterest": "",
        "posCcy": "BTC",
        "availPos": "0.00190433573",
        "avgPx": "62961.4",
        "nonSettleAvgPx": "",
        "upl": "-0.0000033452492717",
        "uplRatio": "-0.0105311101755551",
        "uplLastPx": "-0.0000033199677697",
        "uplRatioLastPx": "-0.0104515220008934",
        "lever": "5",
        "liqPx": "53615.448336593756",
        "markPx": "62891.9",
        "imr": "",
        "margin": "0.000317654",
        "mgnRatio": "9.404143929947395",
        "mmr": "0.0000318005395854",
        "liab": "-99.9998177776581948",
        "liabCcy": "USDT",
        "interest": "0",
        "tradeId": "785524470",
        "optVal": "",
        "pendingCloseOrdLiabVal": "0",
        "notionalUsd": "119.756628017499",
        "adl": "1",
        "ccy": "BTC",
        "last": "62892.9",
        "idxPx": "62890.5",
        "usdPx": "",
        "bePx": "",
        "deltaBS": "",
        "deltaPA": "",
        "gammaBS": "",
        "gammaPA": "",
        "thetaBS": "",
        "thetaPA": "",
        "vegaBS": "",
        "vegaPA": "",
        "spotInUseAmt": "",
        "spotInUseCcy": "",
        "clSpotInUseAmt": "",
        "maxSpotInUseAmt": "",
        "realizedPnl": "",
        "settledPnl": "",
        "pnl": "",
        "fee": "",
        "fundingFee": "",
        "liqPenalty": "",
        "closeOrderAlgo": [],
        "cTime": "1724740225685",
        "uTime": "1724742632153",
        "bizRefId": "",
        "bizRefType": "",
        "currency": "USDT",
        "createTime": "2024-08-27T03:10:25Z",
        "updateTime": "2024-08-27T03:50:32Z"
      }
    ],
    "currency": "USDT"
  }
}
```

### 字段说明

| 字段名 | 描述 |
|--------|------|
| instType | 产品类型：MARGIN(币币杠杆), SWAP(永续合约), FUTURES(交割合约), OPTION(期权) |
| instId | 产品ID，如 BTC-USDT-SWAP |
| mgnMode | 保证金模式：cross(全仓), isolated(逐仓) |
| posId | 持仓ID |
| posSide | 持仓方向：long(开平仓模式开多), short(开平仓模式开空), net(买卖模式) |
| pos | 持仓数量 |
| baseBal | 交易币余额，适用于币币杠杆 |
| quoteBal | 计价币余额，适用于币币杠杆 |
| baseBorrowed | 交易币已借，适用于币币杠杆 |
| baseInterest | 交易币计息，适用于币币杠杆 |
| quoteBorrowed | 计价币已借，适用于币币杠杆 |
| quoteInterest | 计价币计息，适用于币币杠杆 |
| posCcy | 仓位资产币种，仅适用于币币杠杆仓位 |
| availPos | 可平仓数量 |
| avgPx | 开仓均价 |
| nonSettleAvgPx | 未结算均价，不受结算影响的加权开仓价格 |
| upl | 未实现收益（以标记价格计算） |
| uplRatio | 未实现收益率（以标记价格计算） |
| uplLastPx | 以最新成交价格计算的未实现收益 |
| uplRatioLastPx | 以最新成交价格计算的未实现收益率 |
| lever | 杠杆倍数 |
| liqPx | 预估强平价 |
| markPx | 最新标记价格 |
| imr | 初始保证金，仅适用于全仓 |
| margin | 保证金余额，仅适用于逐仓 |
| mgnRatio | 维持保证金率 |
| mmr | 维持保证金 |
| liab | 负债额，仅适用于币币杠杆 |
| liabCcy | 负债币种，仅适用于币币杠杆 |
| interest | 利息 |
| tradeId | 最新成交ID |
| optVal | 期权市值，仅适用于期权 |
| pendingCloseOrdLiabVal | 逐仓杠杆负债对应平仓挂单的数量 |
| notionalUsd | 以美金价值为单位的持仓数量 |
| adl | 信号区，分为5档，从1到5，数字越小代表adl强度越弱 |
| ccy | 占用保证金的币种 |
| last | 最新成交价 |
| idxPx | 最新指数价格 |
| usdPx | 保证金币种的市场最新美金价格，仅适用于期权 |
| bePx | 盈亏平衡价 |
| deltaBS | 美金本位持仓仓位delta，仅适用于期权 |
| deltaPA | 币本位持仓仓位delta，仅适用于期权 |
| gammaBS | 美金本位持仓仓位gamma，仅适用于期权 |
| gammaPA | 币本位持仓仓位gamma，仅适用于期权 |
| thetaBS | 美金本位持仓仓位theta，仅适用于期权 |
| thetaPA | 币本位持仓仓位theta，仅适用于期权 |
| vegaBS | 美金本位持仓仓位vega，仅适用于期权 |
| vegaPA | 币本位持仓仓位vega，仅适用于期权 |
| spotInUseAmt | 现货对冲占用数量，适用于组合保证金模式 |
| spotInUseCcy | 现货对冲占用币种，适用于组合保证金模式 |
| clSpotInUseAmt | 用户自定义现货占用数量，适用于组合保证金模式 |
| maxSpotInUseAmt | 系统计算得到的最大可能现货占用数量，适用于组合保证金模式 |
| realizedPnl | 已实现收益 |
| settledPnl | 已结算收益，仅适用于全仓交割 |
| pnl | 平仓订单累计收益额 |
| fee | 累计手续费金额 |
| fundingFee | 累计资金费用 |
| liqPenalty | 累计爆仓罚金 |
| closeOrderAlgo | 平仓策略委托订单 |
| cTime | 持仓创建时间，Unix时间戳的毫秒数格式 |
| uTime | 最近一次持仓更新时间，Unix时间戳的毫秒数格式 |
| bizRefId | 外部业务id |
| bizRefType | 外部业务类型 |

## 高级功能：基于时间戳的持仓历史查询

### 1. 根据持仓ID获取完整历史

**GET** `/api/v1/account/positions/{posId}/history`

根据持仓ID获取该持仓的完整历史记录，包括当前状态和历史变更。

**路径参数：**
- `posId`: 持仓ID

**查询参数：**
- `currency` (可选): 显示币种
- `limit` (可选): 返回记录数量，默认100
- `includeCurrent` (可选): 是否包含当前持仓时间点，默认false

**请求示例：**
```bash
# 获取持仓ID为123456789的完整历史
curl "http://localhost:8080/api/v1/account/positions/123456789/history"

# 包含当前持仓状态的历史
curl "http://localhost:8080/api/v1/account/positions/123456789/history?includeCurrent=true"
```

**响应示例：**
```json
{
  "success": true,
  "message": "获取持仓完整历史成功",
  "data": {
    "posId": "123456789",
    "currentPosition": {
      "instId": "BTC-USDT",
      "pos": "0.5",
      "uTime": "1699776000000",
      // ... 其他当前持仓字段
    },
    "currentUTime": "1699776000000",
    "history": [
      {
        "instId": "BTC-USDT",
        "type": "1",
        "uTime": "1699689600000",
        // ... 历史持仓字段
      }
    ],
    "hasMore": false,
    "currency": "USDT"
  }
}
```

### 2. 基于当前持仓的历史查询

在历史持仓API中添加 `fromCurrentPositions=true` 参数，系统会自动获取当前持仓的更新时间作为查询基准。

**请求示例：**
```bash
# 获取BTC-USDT从当前持仓时间点往前的历史
curl "http://localhost:8080/api/v1/account/positions-history?instId=BTC-USDT&fromCurrentPositions=true"

# 获取永续合约从当前持仓时间点往前的历史
curl "http://localhost:8080/api/v1/account/positions-history?instType=SWAP&fromCurrentPositions=true"
```

### 3. 时间戳使用说明

- **uTime**: 持仓最近一次更新时间，Unix时间戳毫秒格式
- **before**: 查询此时间戳之前的记录
- **after**: 查询此时间戳之后的记录
- **时间顺序**: 历史持仓按uTime倒序返回（最新的在前）

**实用场景：**
1. **持仓演变追踪**: 查看某个持仓从开仓到当前的完整变化过程
2. **分页查询**: 使用历史记录的uTime作为下一页的before参数
3. **时间点分析**: 基于特定时间点查询之前或之后的持仓状态

## 相关链接

- [OKX API V5 官方文档](https://www.okx.com/docs-v5/zh/)
- [OKX API 状态页面](https://status.okx.com/) 