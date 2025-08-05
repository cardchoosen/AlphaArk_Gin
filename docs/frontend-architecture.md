# 前端架构文档

## 📁 项目结构

```
web/
├── static/
│   ├── css/
│   │   └── main.css                 # 主样式文件
│   └── js/
│       ├── services/               # 服务层
│       │   ├── ApiService.js       # API服务封装
│       │   └── WebSocketService.js # WebSocket服务封装
│       ├── components/             # 组件层
│       │   └── PriceCard.js        # 价格卡片组件
│       └── main.js                 # 应用主文件
└── templates/
    └── index.html                  # HTML模板
```

## 🏗️ 架构设计

### 分层架构

1. **表现层 (Presentation Layer)**
   - `index.html` - HTML结构
   - `main.css` - 样式定义

2. **组件层 (Component Layer)**
   - `PriceCard.js` - 价格显示组件
   - 可扩展其他UI组件

3. **服务层 (Service Layer)**
   - `ApiService.js` - HTTP API调用封装
   - `WebSocketService.js` - WebSocket连接管理

4. **应用层 (Application Layer)**
   - `main.js` - 应用程序主逻辑和协调

## 🔧 核心组件

### AlphaArkApp (主应用类)

**职责:**
- 应用程序生命周期管理
- 服务协调和初始化
- 全局状态管理

**主要方法:**
- `init()` - 初始化应用程序
- `setupEventListeners()` - 设置事件监听器
- `fetchInitialPrice()` - 获取初始价格数据
- `getSystemStatus()` - 获取系统状态

### PriceCard (价格卡片组件)

**职责:**
- 价格数据显示
- 价格变化动画
- 连接状态指示

**主要方法:**
- `updatePrice(data)` - 更新价格显示
- `setConnectionStatus(connected)` - 设置连接状态
- `showLoading()` - 显示加载状态

### WebSocketService (WebSocket服务)

**职责:**
- WebSocket连接管理
- 自动重连机制
- 事件监听器管理

**主要方法:**
- `connect()` - 建立连接
- `on(event, callback)` - 添加事件监听器
- `emit(event, data)` - 触发事件

### PriceApiService (价格API服务)

**职责:**
- HTTP API调用封装
- 价格数据获取
- 配置信息获取

**主要方法:**
- `getPrice(symbol)` - 获取价格数据
- `getOKXConfig()` - 获取OKX配置
- `getInstruments(instType)` - 获取交易对信息

## 🎨 样式架构

### CSS组织结构

1. **基础样式** - 重置和全局样式
2. **布局样式** - 头部、主容器、网格布局
3. **组件样式** - 卡片、按钮等UI组件
4. **动画样式** - 过渡效果和关键帧动画
5. **响应式样式** - 媒体查询和适配

### 设计系统

**颜色主题:**
- 主背景: `#0a0a0a` (深黑)
- 卡片背景: `#1a1a1a` -> `#0f0f0f` (渐变)
- 主色调: `#00d4ff` (青蓝)
- 成功色: `#00ff88` (绿)
- 错误色: `#ff4757` (红)

**字体系统:**
- 主字体: SF Pro Display, 系统字体栈
- 等宽字体: SF Mono, Monaco (价格显示)

## 🔄 数据流

### 实时价格更新流程

1. **初始化阶段**
   ```
   页面加载 -> 应用初始化 -> API获取初始价格 -> WebSocket连接
   ```

2. **实时更新阶段**
   ```
   WebSocket接收数据 -> 数据解析 -> PriceCard更新 -> UI重新渲染
   ```

3. **错误处理**
   ```
   连接断开 -> 自动重连 -> 状态指示器更新 -> 用户通知
   ```

## 🚀 扩展指南

### 添加新组件

1. 在 `web/static/js/components/` 创建组件文件
2. 实现组件类，包含必要的生命周期方法
3. 在 `index.html` 中引入组件脚本
4. 在 `main.js` 中初始化和使用组件

### 添加新服务

1. 在 `web/static/js/services/` 创建服务文件
2. 继承基础服务类或实现标准接口
3. 在应用初始化时注册服务
4. 通过依赖注入使用服务

### 样式扩展

1. 在 `main.css` 中添加新的样式规则
2. 遵循现有的命名约定和结构
3. 使用CSS自定义属性进行主题化
4. 确保响应式设计兼容性

## 🐛 调试工具

### 浏览器控制台

```javascript
// 获取系统状态
await AlphaArk.getStatus()

// 手动重连
AlphaArk.reconnect()

// 断开连接
AlphaArk.disconnect()

// 访问应用实例
AlphaArk.app
```

### 开发者工具

- Network面板：监控API调用和WebSocket连接
- Console面板：查看日志和错误信息
- Elements面板：检查DOM结构和样式
- Application面板：查看存储和缓存

## 📈 性能优化

### 已实现的优化

1. **资源分离** - CSS/JS文件独立，便于缓存
2. **模块化加载** - 按需加载组件和服务
3. **事件优化** - 防抖和节流处理
4. **连接管理** - 智能重连和页面可见性检测

### 未来优化方向

1. **代码分割** - 动态导入大型组件
2. **服务工作者** - 离线支持和缓存策略
3. **虚拟滚动** - 大数据列表优化
4. **预加载** - 关键资源预加载

## 🔒 安全考虑

1. **XSS防护** - 输入验证和输出编码
2. **CSRF防护** - 请求令牌验证
3. **内容安全策略** - CSP头部配置
4. **HTTPS强制** - 生产环境安全传输

## 📚 技术栈

- **HTML5** - 语义化标记
- **CSS3** - 现代样式特性
- **ES6+** - 现代JavaScript特性
- **WebSocket** - 实时通信
- **Fetch API** - HTTP请求
- **JSON** - 数据交换格式