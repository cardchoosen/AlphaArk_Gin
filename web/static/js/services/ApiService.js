// API服务类
class ApiService {
    constructor(baseURL = '') {
        this.baseURL = baseURL;
        this.defaultHeaders = {
            'Content-Type': 'application/json',
            'Accept': 'application/json'
        };
    }

    async request(url, options = {}) {
        const config = {
            headers: { ...this.defaultHeaders, ...options.headers },
            ...options
        };

        try {
            const response = await fetch(this.baseURL + url, config);
            
            if (!response.ok) {
                throw new Error(`HTTP Error: ${response.status} ${response.statusText}`);
            }

            const data = await response.json();
            return data;
        } catch (error) {
            console.error('API请求失败:', error);
            throw error;
        }
    }

    // GET请求
    async get(url, params = {}) {
        const queryString = new URLSearchParams(params).toString();
        const fullUrl = queryString ? `${url}?${queryString}` : url;
        
        return this.request(fullUrl, {
            method: 'GET'
        });
    }

    // POST请求
    async post(url, data = {}) {
        return this.request(url, {
            method: 'POST',
            body: JSON.stringify(data)
        });
    }

    // PUT请求
    async put(url, data = {}) {
        return this.request(url, {
            method: 'PUT',
            body: JSON.stringify(data)
        });
    }

    // DELETE请求
    async delete(url) {
        return this.request(url, {
            method: 'DELETE'
        });
    }
}

// 价格API服务
class PriceApiService extends ApiService {
    constructor() {
        super('/api/v1');
    }

    // 获取指定交易对价格
    async getPrice(symbol) {
        return this.get(`/price/${symbol}`);
    }

    // 获取OKX配置信息
    async getOKXConfig() {
        return this.get('/okx/config');
    }

    // 获取交易对信息
    async getInstruments(instType = 'SPOT') {
        return this.get('/okx/instruments', { instType });
    }
}

// 导出服务
if (typeof module !== 'undefined' && module.exports) {
    module.exports = { ApiService, PriceApiService };
} else {
    window.ApiService = ApiService;
    window.PriceApiService = PriceApiService;
}