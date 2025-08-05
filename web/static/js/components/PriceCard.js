// 价格卡片组件
class PriceCard {
    constructor(symbol = 'BTC-USDT') {
        this.symbol = symbol;
        this.lastPrice = null;
        this.elements = this.getElements();
    }

    getElements() {
        return {
            price: document.getElementById('btc-price'),
            change: document.getElementById('price-change'),
            time: document.getElementById('price-time'),
            status: document.getElementById('connection-status')
        };
    }

    updatePrice(data) {
        if (data.symbol !== this.symbol || !data.price) return;

        const currentPrice = parseFloat(data.price);
        const formattedPrice = this.formatPrice(currentPrice);

        this.elements.price.textContent = `$${formattedPrice}`;
        this.updateChange(data);
        this.updateTime();
        this.animatePrice(currentPrice);

        this.lastPrice = currentPrice;
    }

    formatPrice(price) {
        return price.toLocaleString('en-US', {
            minimumFractionDigits: 2,
            maximumFractionDigits: 2
        });
    }

    updateChange(data) {
        if (!data.change24h) return;

        const change = parseFloat(data.change24h);
        const changePercent = data.changePercent24h ? parseFloat(data.changePercent24h) : 0;
        
        this.elements.change.textContent = `${change >= 0 ? '+' : ''}${change.toFixed(2)} (${changePercent.toFixed(2)}%)`;
        this.elements.change.className = `price-change ${change >= 0 ? 'positive' : 'negative'}`;
    }

    updateTime() {
        this.elements.time.textContent = `更新时间: ${new Date().toLocaleTimeString()}`;
    }

    animatePrice(currentPrice) {
        if (this.lastPrice === null) return;

        let color = '#00d4ff';
        if (currentPrice > this.lastPrice) {
            color = '#00ff88';
        } else if (currentPrice < this.lastPrice) {
            color = '#ff4757';
        }

        this.elements.price.style.color = color;
        
        setTimeout(() => {
            this.elements.price.style.color = '#00d4ff';
        }, 1000);
    }

    setConnectionStatus(connected) {
        if (connected) {
            this.elements.status.classList.remove('disconnected');
            this.elements.time.textContent = '实时连接';
        } else {
            this.elements.status.classList.add('disconnected');
            this.elements.time.textContent = '连接断开';
        }
    }

    showLoading() {
        this.elements.price.innerHTML = '<span class="loading"></span>';
        this.elements.change.textContent = '--';
        this.elements.time.textContent = '连接中...';
    }
}

// 导出组件
if (typeof module !== 'undefined' && module.exports) {
    module.exports = PriceCard;
} else {
    window.PriceCard = PriceCard;
}