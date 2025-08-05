#!/bin/bash

# OKX API 测试脚本
# 测试系统时间接口和持仓查询接口

BASE_URL="http://localhost:8080"

echo "=== OKX API 测试 ==="
echo ""

# 1. 测试系统时间接口
echo "1. 测试系统时间接口..."
echo "GET $BASE_URL/api/v1/okx/system-time"
response=$(curl -s "$BASE_URL/api/v1/okx/system-time")
echo "响应: $response"
echo ""

# 2. 测试获取当前持仓信息
echo "2. 测试获取当前持仓信息..."
echo "GET $BASE_URL/api/v1/account/positions"
response=$(curl -s "$BASE_URL/api/v1/account/positions")
echo "响应: $response"
echo ""

# 3. 测试获取特定产品的持仓信息
echo "3. 测试获取BTC-USDT的持仓信息..."
echo "GET $BASE_URL/api/v1/account/positions?instId=BTC-USDT"
response=$(curl -s "$BASE_URL/api/v1/account/positions?instId=BTC-USDT")
echo "响应: $response"
echo ""

# 4. 测试获取永续合约的持仓信息
echo "4. 测试获取永续合约的持仓信息..."
echo "GET $BASE_URL/api/v1/account/positions?instType=SWAP"
response=$(curl -s "$BASE_URL/api/v1/account/positions?instType=SWAP")
echo "响应: $response"
echo ""

# 5. 测试获取历史持仓信息
echo "5. 测试获取历史持仓信息..."
echo "GET $BASE_URL/api/v1/account/positions-history"
response=$(curl -s "$BASE_URL/api/v1/account/positions-history")
echo "响应: $response"
echo ""

# 6. 测试获取特定产品的历史持仓
echo "6. 测试获取BTC-USDT的历史持仓..."
echo "GET $BASE_URL/api/v1/account/positions-history?instId=BTC-USDT&limit=5"
response=$(curl -s "$BASE_URL/api/v1/account/positions-history?instId=BTC-USDT&limit=5")
echo "响应: $response"
echo ""

echo "=== 测试完成 ===" 