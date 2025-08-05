// WebSocket服务类
class WebSocketService {
    constructor(url, options = {}) {
        this.url = url;
        this.options = {
            maxReconnectAttempts: 5,
            reconnectInterval: 3000,
            ...options
        };
        
        this.ws = null;
        this.reconnectAttempts = 0;
        this.listeners = {
            open: [],
            message: [],
            close: [],
            error: []
        };
    }

    connect() {
        try {
            this.ws = new WebSocket(this.url);
            this.setupEventListeners();
        } catch (error) {
            console.error('WebSocket连接失败:', error);
            this.handleError(error);
        }
    }

    setupEventListeners() {
        this.ws.onopen = (event) => {
            console.log('WebSocket连接已建立');
            this.reconnectAttempts = 0;
            this.emit('open', event);
        };

        this.ws.onmessage = (event) => {
            try {
                const data = JSON.parse(event.data);
                this.emit('message', data);
            } catch (error) {
                console.error('解析WebSocket消息失败:', error);
                this.emit('error', error);
            }
        };

        this.ws.onclose = (event) => {
            console.log('WebSocket连接已关闭');
            this.emit('close', event);
            this.handleReconnect();
        };

        this.ws.onerror = (error) => {
            console.error('WebSocket错误:', error);
            this.emit('error', error);
        };
    }

    handleReconnect() {
        if (this.reconnectAttempts < this.options.maxReconnectAttempts) {
            this.reconnectAttempts++;
            const delay = this.options.reconnectInterval * this.reconnectAttempts;
            
            console.log(`尝试重连... (${this.reconnectAttempts}/${this.options.maxReconnectAttempts})`);
            
            setTimeout(() => {
                this.connect();
            }, delay);
        } else {
            console.error('达到最大重连次数，停止重连');
        }
    }

    send(data) {
        if (this.ws && this.ws.readyState === WebSocket.OPEN) {
            this.ws.send(JSON.stringify(data));
        } else {
            console.warn('WebSocket未连接，无法发送消息');
        }
    }

    close() {
        if (this.ws) {
            this.ws.close();
            this.ws = null;
        }
    }

    // 事件监听器管理
    on(event, callback) {
        if (this.listeners[event]) {
            this.listeners[event].push(callback);
        }
    }

    off(event, callback) {
        if (this.listeners[event]) {
            const index = this.listeners[event].indexOf(callback);
            if (index > -1) {
                this.listeners[event].splice(index, 1);
            }
        }
    }

    emit(event, data) {
        if (this.listeners[event]) {
            this.listeners[event].forEach(callback => {
                try {
                    callback(data);
                } catch (error) {
                    console.error(`事件监听器执行失败 (${event}):`, error);
                }
            });
        }
    }

    // 获取连接状态
    getReadyState() {
        return this.ws ? this.ws.readyState : WebSocket.CLOSED;
    }

    isConnected() {
        return this.ws && this.ws.readyState === WebSocket.OPEN;
    }
}

// 导出服务
if (typeof module !== 'undefined' && module.exports) {
    module.exports = WebSocketService;
} else {
    window.WebSocketService = WebSocketService;
}