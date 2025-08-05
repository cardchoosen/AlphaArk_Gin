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
OKX_API_KEY=6e5bb45e-bffe-42d7-932f-9c828dc3a533
OKX_SECRET_KEY=CB3D7E7D80F9FD3E3CB838017AC0CA1F
OKX_PASSPHRASE=DWQIUD.e39081
OKX_IP=
OKX_REMARK=Gin项目
OKX_PERMISSIONS=读取/提现/交易
OKX_BASE_URL=https://www.okx.com
OKX_IS_TEST=false
```

## 开发指南

- 遵循Go官方代码规范
- 使用依赖注入管理依赖关系
- 编写单元测试和集成测试
- 使用日志记录关键操作 