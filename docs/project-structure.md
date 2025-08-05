# 项目结构说明

## 概述

AlphaArk_Gin 是一个基于 Go Gin 框架的量化交易系统，集成了 OKX 交易所 API，提供实时价格监控、账户管理、持仓查询等功能。

## 项目结构

```
AlphaArk_Gin/
├── cmd/                    # 应用程序入口
│   └── server/
│       └── main.go        # 主程序入口
├── internal/              # 内部包
│   ├── api/              # API层
│   │   ├── account_routes.go    # 账户相关路由
│   │   ├── okx_client.go        # OKX API客户端
│   │   ├── okx_routes.go        # OKX相关路由
│   │   ├── price_routes.go      # 价格相关路由
│   │   ├── routes.go            # 主路由配置
│   │   └── websocket.go         # WebSocket服务
│   ├── config/           # 配置管理
│   │   └── config.go     # 配置结构体和加载逻辑
│   ├── middleware/       # 中间件
│   │   └── middleware.go # CORS、日志、恢复等中间件
│   ├── models/          # 数据模型
│   │   └── account.go   # 账户相关模型
│   └── service/         # 业务逻辑层
│       ├── account_service.go   # 账户服务
│       └── price_service.go     # 价格服务
├── web/                 # 前端资源
│   ├── static/         # 静态资源
│   │   ├── css/
│   │   │   └── main.css        # 主样式文件
│   │   └── js/
│   │       ├── components/     # 前端组件
│   │       │   ├── AccountCard.js    # 账户卡片组件
│   │       │   └── PriceCard.js      # 价格卡片组件
│   │       ├── services/       # 前端服务
│   │       │   ├── ApiService.js     # API服务
│   │       │   └── WebSocketService.js # WebSocket服务
│   │       └── main.js         # 主JavaScript文件
│   └── templates/      # HTML模板
│       └── index.html  # 主页面模板
├── docs/               # 文档
│   ├── okx-api.md      # OKX API集成文档
│   ├── web-architecture.md # Web前端架构文档
│   └── project-structure.md # 项目结构文档
├── tests/              # 测试文件
│   ├── okx_test.go     # OKX API测试
│   ├── positions_test.go # 持仓API测试
│   └── positions_history_test.go # 持仓历史API测试
├── scripts/            # 脚本文件
│   └── setup-env.sh    # 环境设置脚本
├── docker-compose.yml  # Docker Compose配置
├── Dockerfile         # Docker镜像配置
├── env.example        # 环境变量示例
├── go.mod             # Go模块文件
├── go.sum             # Go依赖校验文件
├── Makefile           # 构建脚本
├── README.md          # 项目说明
└── test_okx_positions.sh # OKX持仓测试脚本
```

## 核心功能模块

### 1. API层 (`internal/api/`)

- **account_routes.go**: 账户管理API
  - 账户余额查询
  - 账户汇总信息
  - 币种管理
  - 持仓查询
  - 持仓历史查询

- **okx_client.go**: OKX API客户端
  - 时间同步功能
  - 系统时间获取
  - 持仓查询
  - 交易对信息获取

- **okx_routes.go**: OKX相关路由
  - 交易对信息API
  - 系统时间API
  - 配置信息API

- **price_routes.go**: 价格相关路由
  - 实时价格获取
  - 价格历史数据

- **websocket.go**: WebSocket服务
  - 实时价格推送
  - 连接管理

### 2. 服务层 (`internal/service/`)

- **account_service.go**: 账户服务
  - 账户余额管理
  - 汇率转换
  - 时间同步
  - 持仓管理

- **price_service.go**: 价格服务
  - 价格数据获取
  - 价格格式化

### 3. 前端 (`web/`)

- **AccountCard.js**: 账户卡片组件
  - 总资产显示
  - 币种切换
  - 收益统计

- **PriceCard.js**: 价格卡片组件
  - 实时价格显示
  - 连接状态管理

- **ApiService.js**: API服务
  - HTTP请求封装
  - 错误处理

- **WebSocketService.js**: WebSocket服务
  - 实时数据连接
  - 自动重连

## 主要API端点

### 账户管理
- `GET /api/v1/account/balance` - 获取账户余额
- `GET /api/v1/account/summary/:currency` - 获取账户汇总
- `GET /api/v1/account/positions` - 获取当前持仓
- `GET /api/v1/account/positions-history` - 获取持仓历史
- `GET /api/v1/account/currencies` - 获取支持币种
- `POST /api/v1/account/currency` - 设置默认币种

### OKX API
- `GET /api/v1/okx/instruments` - 获取交易对信息
- `GET /api/v1/okx/system-time` - 获取系统时间
- `GET /api/v1/okx/config` - 获取配置信息

### 价格API
- `GET /api/v1/price/:symbol` - 获取价格信息

### WebSocket
- `WS /ws/price` - 实时价格推送

## 技术栈

### 后端
- **Go 1.21+**: 主要开发语言
- **Gin**: Web框架
- **WebSocket**: 实时通信
- **OKX API V5**: 交易所API

### 前端
- **原生JavaScript**: 无框架依赖
- **WebSocket**: 实时数据
- **CSS3**: 样式设计

### 部署
- **Docker**: 容器化部署
- **Docker Compose**: 服务编排

## 配置说明

### 环境变量
```bash
# 服务器配置
PORT=8080
ENVIRONMENT=development

# OKX API配置
OKX_API_KEY=your_api_key
OKX_SECRET_KEY=your_secret_key
OKX_PASSPHRASE=your_passphrase
OKX_BASE_URL=https://www.okx.com
OKX_IS_TEST=false
```

## 开发指南

### 本地开发
```bash
# 安装依赖
go mod download

# 设置环境变量
cp env.example .env
# 编辑 .env 文件

# 运行开发服务器
go run cmd/server/main.go
```

### 构建部署
```bash
# 构建
make build

# Docker部署
docker-compose up -d
```

## 测试

```bash
# 运行所有测试
make test

# 运行特定测试
go test ./tests -v
```

## 监控和日志

- 所有API调用都有详细日志
- WebSocket连接状态监控
- 错误处理和恢复机制
- 时间同步状态监控

## 安全特性

- CORS中间件
- 请求日志记录
- 错误恢复机制
- API密钥安全存储
- 时间戳验证

## 性能优化

- 汇率缓存机制
- WebSocket连接复用
- 静态资源压缩
- 数据库连接池（待实现）

## 扩展计划

- [ ] 数据库集成（PostgreSQL/MySQL）
- [ ] Redis缓存
- [ ] 用户认证系统
- [ ] 交易策略引擎
- [ ] 风险管理模块
- [ ] 移动端适配 