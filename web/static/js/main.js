// 应用程序主类
class AlphaArkApp {
    constructor() {
        this.priceCard = null;
        this.wsService = null;
        this.apiService = null;
        this.isInitialized = false;
    }

    async init() {
        if (this.isInitialized) return;

        console.log('AlphaArk 量化交易系统初始化...');

        try {
            // 初始化服务
            this.apiService = new PriceApiService();
            this.priceCard = new PriceCard('BTC-USDT');
            this.accountCard = new AccountCard();
            
            // 初始化WebSocket
            const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
            const wsUrl = `${protocol}//${window.location.host}/ws/price`;
            this.wsService = new WebSocketService(wsUrl);

            // 设置事件监听器
            this.setupEventListeners();

            // 显示加载状态
            this.priceCard.showLoading();

            // 获取初始价格数据
            await this.fetchInitialPrice();

            // 连接WebSocket
            this.wsService.connect();

            this.isInitialized = true;
            console.log('系统初始化完成');

        } catch (error) {
            console.error('系统初始化失败:', error);
        }
    }

    setupEventListeners() {
        // WebSocket事件监听
        this.wsService.on('open', () => {
            this.priceCard.setConnectionStatus(true);
        });

        this.wsService.on('close', () => {
            this.priceCard.setConnectionStatus(false);
        });

        this.wsService.on('message', (data) => {
            this.priceCard.updatePrice(data);
        });

        this.wsService.on('error', (error) => {
            console.error('WebSocket错误:', error);
            this.priceCard.setConnectionStatus(false);
        });

        // 页面卸载时清理资源
        window.addEventListener('beforeunload', () => {
            this.cleanup();
        });

        // 页面可见性变化处理
        document.addEventListener('visibilitychange', () => {
            if (document.hidden) {
                console.log('页面隐藏，暂停WebSocket连接');
            } else {
                console.log('页面显示，恢复WebSocket连接');
                if (!this.wsService.isConnected()) {
                    this.wsService.connect();
                }
            }
        });

        // 全局币种变化事件监听
        window.addEventListener('currencyChanged', (event) => {
            console.log('全局币种已变更:', event.detail);
            // 这里可以通知其他组件更新显示单位
        });
    }

    async fetchInitialPrice() {
        try {
            const response = await this.apiService.getPrice('BTC-USDT');
            if (response.success && response.data) {
                this.priceCard.updatePrice(response.data);
            }
        } catch (error) {
            console.error('获取初始价格失败:', error);
        }
    }

    cleanup() {
        if (this.wsService) {
            this.wsService.close();
        }
    }

    // 公共API方法
    async getSystemStatus() {
        try {
            const config = await this.apiService.getOKXConfig();
            return {
                wsConnected: this.wsService ? this.wsService.isConnected() : false,
                apiConfigured: config.success && config.data.hasApiKey,
                initialized: this.isInitialized
            };
        } catch (error) {
            console.error('获取系统状态失败:', error);
            return {
                wsConnected: false,
                apiConfigured: false,
                initialized: this.isInitialized
            };
        }
    }
}

// 全局应用实例
const app = new AlphaArkApp();

// 页面初始化
document.addEventListener('DOMContentLoaded', function() {
    // 延迟初始化，确保所有资源加载完成
    setTimeout(() => {
        app.init();
    }, 1000);
});

// 导出全局对象供调试使用
window.AlphaArk = {
    app,
    version: '1.0.0',
    
    // 调试方法
    async getStatus() {
        return await app.getSystemStatus();
    },
    
    reconnect() {
        if (app.wsService) {
            app.wsService.connect();
        }
    },
    
    disconnect() {
        if (app.wsService) {
            app.wsService.close();
        }
    }
};