# Web前端架构文档

## 概述

AlphaArk_Gin 的Web前端采用原生JavaScript开发，无框架依赖，提供实时价格监控、账户管理等功能。前端采用组件化架构，支持实时数据更新和响应式设计。

## 目录结构

```
web/
├── templates/
│   └── index.html              # 主页面模板
└── static/
    ├── css/
    │   └── main.css            # 主样式文件
    └── js/
        ├── main.js             # 应用程序主文件
        ├── components/         # 前端组件
        │   ├── AccountCard.js  # 账户卡片组件
        │   └── PriceCard.js    # 价格卡片组件
        └── services/           # 服务层
            ├── ApiService.js   # API服务
            └── WebSocketService.js # WebSocket服务
```

## 核心组件

### 1. 主应用程序 (`main.js`)

**文件**: `web/static/js/main.js`

**功能**:
- 应用程序初始化和生命周期管理
- 组件协调和事件管理
- WebSocket连接管理
- 全局状态管理

**主要类**: `AlphaArkApp`

**核心方法**:
- `init()`: 应用程序初始化
- `setupEventListeners()`: 设置事件监听器
- `fetchInitialPrice()`: 获取初始价格数据
- `cleanup()`: 资源清理
- `getSystemStatus()`: 获取系统状态

**全局对象**: `window.AlphaArk` - 提供调试和开发工具

### 2. 账户卡片组件 (`AccountCard.js`)

**文件**: `web/static/js/components/AccountCard.js`

**功能**:
- 总资产显示和格式化
- 币种切换（USDT、CNY、USD、BTC）
- 收益统计和时间周期选择
- 实时数据更新

**主要特性**:
- 支持多币种显示（精确度自动调整）
- 时间周期切换（1日、1周、1月、半年）
- 自动加载支持币种
- 错误处理和加载状态

**核心方法**:
- `init()`: 组件初始化
- `changeCurrency()`: 切换币种
- `changePeriod()`: 切换时间周期
- `updateAccountData()`: 更新账户数据
- `formatNumber()`: 数字格式化

**支持的币种**:
- **USDT**: 2位小数
- **CNY**: 2位小数
- **USD**: 2位小数
- **BTC**: 5位小数

### 3. 价格卡片组件 (`PriceCard.js`)

**文件**: `web/static/js/components/PriceCard.js`

**功能**:
- 实时价格显示
- WebSocket连接状态管理
- 价格变化指示器
- 时间戳显示

**主要特性**:
- 实时价格更新
- 连接状态可视化
- 价格变化颜色指示
- 自动重连机制

**核心方法**:
- `updatePrice()`: 更新价格显示
- `setConnectionStatus()`: 设置连接状态
- `showLoading()`: 显示加载状态
- `formatPrice()`: 价格格式化

### 4. API服务 (`ApiService.js`)

**文件**: `web/static/js/services/ApiService.js`

**功能**:
- HTTP请求封装
- 错误处理
- 响应数据格式化

**主要类**:
- `ApiService`: 基础HTTP服务
- `PriceApiService`: 价格相关API

**核心方法**:
- `request()`: 通用请求方法
- `get()`: GET请求
- `post()`: POST请求
- `getPrice()`: 获取价格
- `getOKXConfig()`: 获取OKX配置

### 5. WebSocket服务 (`WebSocketService.js`)

**文件**: `web/static/js/services/WebSocketService.js`

**功能**:
- WebSocket连接管理
- 自动重连机制
- 事件分发
- 连接状态监控

**主要特性**:
- 自动重连（指数退避）
- 心跳检测
- 事件监听器管理
- 连接状态跟踪

**核心方法**:
- `connect()`: 建立连接
- `close()`: 关闭连接
- `send()`: 发送消息
- `on()`: 添加事件监听器
- `reconnect()`: 重新连接

## 页面结构

### 主页面 (`index.html`)

**布局结构**:
```
┌─────────────────────────────────────┐
│             头部导航                  │
├─────────────────────────────────────┤
│  ┌─────────────┐ ┌─────────────────┐ │
│  │  价格卡片    │ │    账户卡片      │ │
│  │ BTC/USDT    │ │    总资产        │ │
│  └─────────────┘ └─────────────────┘ │
│  ┌─────────────┐ ┌─────────────────┐ │
│  │  当前持仓    │ │   活跃订单      │ │
│  │  (占位)     │ │    (占位)       │ │
│  └─────────────┘ └─────────────────┘ │
│  ┌─────────────┐ ┌─────────────────┐ │
│  │  市场概览    │ │   交易历史      │ │
│  │  (占位)     │ │    (占位)       │ │
│  └─────────────┘ └─────────────────┘ │
└─────────────────────────────────────┘
```

**当前实现的功能**:
- ✅ 实时价格显示
- ✅ 账户总资产显示
- ✅ 币种切换
- ✅ 收益统计
- ⏳ 持仓信息（占位）
- ⏳ 订单管理（占位）
- ⏳ 市场概览（占位）
- ⏳ 交易历史（占位）

## 样式系统

### 主样式文件 (`main.css`)

**设计原则**:
- 响应式设计
- 现代化UI
- 深色主题
- 卡片式布局

**主要样式类**:
- `.header`: 头部导航
- `.card`: 卡片容器
- `.price-card`: 价格卡片
- `.account-card`: 账户卡片
- `.placeholder-card`: 占位卡片

**颜色系统**:
- 主色调: 深蓝色系
- 成功色: 绿色
- 警告色: 橙色
- 错误色: 红色
- 中性色: 灰色系

## 数据流

### 1. 价格数据流
```
OKX API → 后端服务 → WebSocket → PriceCard → UI更新
```

### 2. 账户数据流
```
用户操作 → AccountCard → API请求 → 后端服务 → 响应 → UI更新
```

### 3. 币种切换流
```
用户选择 → AccountCard.changeCurrency() → API请求 → 后端更新 → UI刷新
```

## 事件系统

### 全局事件
- `currencyChanged`: 币种变化事件
- `priceUpdated`: 价格更新事件
- `connectionStatusChanged`: 连接状态变化

### 组件事件
- `periodChanged`: 时间周期变化
- `dataLoaded`: 数据加载完成
- `error`: 错误事件

## 错误处理

### 网络错误
- API请求失败自动重试
- WebSocket断开自动重连
- 友好的错误提示

### 数据错误
- 数据格式验证
- 默认值处理
- 错误状态显示

## 性能优化

### 1. 资源加载
- 脚本按需加载
- CSS压缩
- 静态资源缓存

### 2. 数据更新
- 防抖处理
- 增量更新
- 缓存机制

### 3. 内存管理
- 事件监听器清理
- 定时器管理
- 对象引用清理

## 浏览器兼容性

### 支持的浏览器
- Chrome 80+
- Firefox 75+
- Safari 13+
- Edge 80+

### 使用的现代特性
- ES6+ 语法
- WebSocket API
- Fetch API
- CSS Grid/Flexbox
- Custom Events

## 开发工具

### 调试功能
```javascript
// 全局调试对象
window.AlphaArk = {
    app,           // 主应用实例
    version,       // 版本信息
    getStatus(),   // 获取系统状态
    reconnect(),   // 重新连接
    disconnect()   // 断开连接
}
```

### 开发模式
- 详细的控制台日志
- 错误堆栈跟踪
- 性能监控
- 网络请求调试

## 扩展指南

### 添加新组件
1. 在 `components/` 目录创建新组件文件
2. 在 `index.html` 中引入脚本
3. 在 `main.js` 中初始化组件
4. 添加相应的CSS样式

### 添加新API
1. 在 `services/ApiService.js` 中添加新方法
2. 在组件中调用API
3. 处理响应数据
4. 更新UI显示

### 添加新页面
1. 创建新的HTML模板
2. 添加路由处理
3. 创建对应的JavaScript文件
4. 更新导航菜单

## 维护指南

### 日常维护
- 检查WebSocket连接状态
- 监控API响应时间
- 查看错误日志
- 更新依赖版本

### 性能监控
- 页面加载时间
- API响应时间
- WebSocket连接稳定性
- 内存使用情况

### 安全考虑
- API密钥安全
- XSS防护
- CSRF防护
- 输入验证

## 未来规划

### 短期目标
- [ ] 实现持仓信息显示
- [ ] 添加订单管理功能
- [ ] 完善市场概览
- [ ] 添加交易历史

### 中期目标
- [ ] 移动端适配
- [ ] PWA支持
- [ ] 主题切换
- [ ] 多语言支持

### 长期目标
- [ ] 微前端架构
- [ ] 组件库建设
- [ ] 自动化测试
- [ ] CI/CD集成 