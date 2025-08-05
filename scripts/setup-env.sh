#!/bin/bash

# OKX API 环境配置脚本
# 使用此脚本可以快速创建 .env 配置文件

echo "==================================="
echo "OKX API 环境配置向导"
echo "==================================="
echo

# 检查是否已存在 .env 文件
if [ -f ".env" ]; then
    echo "⚠️  检测到已存在 .env 文件"
    read -p "是否要覆盖现有配置？(y/n): " confirm
    if [ "$confirm" != "y" ] && [ "$confirm" != "Y" ]; then
        echo "❌ 配置已取消"
        exit 0
    fi
fi

echo "请输入您的OKX API信息："
echo

# 获取API密钥信息
read -p "🔑 OKX API Key: " api_key
if [ -z "$api_key" ]; then
    echo "❌ API Key不能为空"
    exit 1
fi

read -s -p "🔐 OKX Secret Key: " secret_key
echo
if [ -z "$secret_key" ]; then
    echo "❌ Secret Key不能为空"
    exit 1
fi

read -s -p "🔒 OKX Passphrase: " passphrase
echo
if [ -z "$passphrase" ]; then
    echo "❌ Passphrase不能为空"
    exit 1
fi

# 环境选择
echo
echo "请选择环境："
echo "1) 正式环境 (推荐)"
echo "2) 测试环境"
read -p "选择 (1-2, 默认1): " env_choice

base_url="https://www.okx.com"
is_test="false"
if [ "$env_choice" = "2" ]; then
    base_url="https://www.okx.com"
    is_test="true"
    echo "✅ 已选择测试环境"
else
    echo "✅ 已选择正式环境"
fi

# 创建 .env 文件
cat > .env << EOF
# 应用配置
ENVIRONMENT=development
PORT=8080

# 数据库配置
DATABASE_URL=postgres://username:password@localhost:5432/dbname

# JWT配置
JWT_SECRET=$(openssl rand -base64 32)

# 日志配置
LOG_LEVEL=debug

# 跨域配置
CORS_ALLOW_ORIGIN=*

# OKX API配置
OKX_API_KEY=$api_key
OKX_SECRET_KEY=$secret_key
OKX_PASSPHRASE=$passphrase
OKX_IP=
OKX_REMARK=AlphaArk_Gin项目
OKX_PERMISSIONS=读取/提现/交易
OKX_BASE_URL=$base_url
OKX_IS_TEST=$is_test
EOF

echo
echo "✅ 配置文件已创建！"
echo
echo "📋 配置摘要："
echo "   - API Key: ${api_key:0:8}..."
echo "   - 环境: $([ "$is_test" = "true" ] && echo "测试环境" || echo "正式环境")"
echo "   - 配置文件: .env"
echo
echo "🚀 现在可以启动应用："
echo "   go run cmd/server/main.go"
echo
echo "📖 更多信息请查看: docs/okx-api-setup.md"