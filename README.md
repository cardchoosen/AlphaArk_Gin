# AlphaArk Gin Project

这是一个基于Gin框架的Go Web项目。

## 项目结构

```
AlphaArk_Gin/
├── cmd/                    # 应用程序入口点
│   └── server/            # 服务器启动程序
│       └── main.go
├── internal/              # 私有应用程序和库代码
│   ├── api/              # API处理器
│   ├── config/           # 配置管理
│   ├── database/         # 数据库相关
│   ├── middleware/       # 中间件
│   ├── models/           # 数据模型
│   ├── repository/       # 数据访问层
│   ├── service/          # 业务逻辑层
│   └── utils/            # 工具函数
├── pkg/                  # 可以被外部应用程序使用的库代码
├── web/                  # Web静态文件
│   ├── static/           # 静态资源
│   └── templates/        # 模板文件
├── docs/                 # 文档
├── scripts/              # 脚本文件
├── tests/                # 测试文件
├── .env.example          # 环境变量示例
├── .gitignore           # Git忽略文件
├── go.mod               # Go模块文件
├── go.sum               # Go依赖校验文件
└── README.md            # 项目说明
```

## 快速开始

1. 克隆项目
2. 复制 `env.example` 为 `.env` 并配置环境变量
3. 配置OKX API信息（见下方说明）
4. 运行 `go mod tidy` 安装依赖
5. 运行 `go run cmd/server/main.go` 启动服务器

### OKX API配置

在 `.env` 文件中配置以下OKX API信息：

```bash
# OKX API配置
OKX_API_KEY=your_api_key_here
OKX_SECRET_KEY=your_secret_key_here
OKX_PASSPHRASE=your_passphrase_here
OKX_IP=
OKX_REMARK=Gin项目
OKX_PERMISSIONS=读取/提现/交易
OKX_BASE_URL=https://www.okx.com
OKX_IS_TEST=false
```

## 功能特性

- ✅ Gin Web框架
- ✅ 结构化项目布局
- ✅ 配置管理
- ✅ 中间件支持
- ✅ 静态文件服务
- ✅ 模板渲染
- ✅ 用户管理API
- ✅ 价格查询API
- ✅ WebSocket支持
- ✅ OKX API集成
- ✅ 账户余额查询
- ✅ 盈亏分析
- ✅ 多币种支持
- ✅ 当前持仓信息查询
- ✅ 历史持仓信息查询

## API端点

### 账户相关API

- `GET /api/v1/account/balance` - 获取账户余额
- `GET /api/v1/account/positions` - 获取当前持仓信息
- `GET /api/v1/account/positions-history` - 获取历史持仓信息
- `GET /api/v1/account/positions/{posId}/history` - 获取指定持仓的完整历史
- `GET /api/v1/account/profit-loss` - 获取盈亏信息
- `GET /api/v1/account/summary` - 获取账户汇总

### 价格相关API

- `GET /api/v1/price/:symbol` - 获取指定币种价格
- `WebSocket /ws` - 实时价格推送

### OKX API

- `GET /api/v1/okx/instruments` - 获取交易对信息
- `GET /api/v1/okx/config` - 获取API配置信息

## 开发指南

- 遵循Go官方代码规范
- 使用依赖注入管理依赖关系
- 编写单元测试和集成测试
- 使用日志记录关键操作

## 测试

运行所有测试：
```bash
go test ./...
```

运行特定测试：
```bash
go test ./tests -v
``` 