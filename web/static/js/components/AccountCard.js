// 账户余额卡片组件
class AccountCard {
    constructor() {
        this.currentCurrency = 'USDT';
        this.currentPeriod = '1d';
        this.supportedCurrencies = [];
        this.accountData = null;
        this.elements = this.getElements();
        this.init();
    }

    getElements() {
        return {
            totalAssets: document.getElementById('total-assets'),
            currencySymbol: document.getElementById('currency-symbol'),
            currencySelector: document.getElementById('currency-selector'),
            profitAmount: document.getElementById('profit-amount'),
            profitPercent: document.getElementById('profit-percent'),
            profitPeriod: document.getElementById('profit-period'),
            periodButtons: document.querySelectorAll('.period-btn'),
            loadingIndicator: document.getElementById('account-loading')
        };
    }

    async init() {
        try {
            await this.loadSupportedCurrencies();
            await this.loadDefaultCurrency();
            this.setupEventListeners();
            await this.updateAccountData();
        } catch (error) {
            console.error('账户卡片初始化失败:', error);
            this.showError('初始化失败');
        }
    }

    async loadSupportedCurrencies() {
        try {
            const response = await fetch('/api/v1/account/currencies');
            const result = await response.json();
            
            if (result.success) {
                this.supportedCurrencies = result.data;
                this.updateCurrencySelector();
            }
        } catch (error) {
            console.error('加载支持币种失败:', error);
        }
    }

    async loadDefaultCurrency() {
        try {
            const response = await fetch('/api/v1/account/currency');
            const result = await response.json();
            
            if (result.success) {
                this.currentCurrency = result.data.currency;
                this.updateCurrencyDisplay();
            }
        } catch (error) {
            console.error('加载默认币种失败:', error);
        }
    }

    updateCurrencySelector() {
        if (!this.elements.currencySelector) return;

        this.elements.currencySelector.innerHTML = '';
        
        this.supportedCurrencies.forEach(currency => {
            const option = document.createElement('option');
            option.value = currency.currency;
            option.textContent = `${currency.symbol} ${currency.currency}`;
            option.selected = currency.currency === this.currentCurrency;
            this.elements.currencySelector.appendChild(option);
        });
    }

    updateCurrencyDisplay() {
        const currency = this.supportedCurrencies.find(c => c.currency === this.currentCurrency);
        if (currency && this.elements.currencySymbol) {
            this.elements.currencySymbol.textContent = currency.symbol;
        }
    }

    setupEventListeners() {
        // 币种选择器事件
        if (this.elements.currencySelector) {
            this.elements.currencySelector.addEventListener('change', async (e) => {
                await this.changeCurrency(e.target.value);
            });
        }

        // 时间周期按钮事件
        this.elements.periodButtons.forEach(btn => {
            btn.addEventListener('click', (e) => {
                this.changePeriod(e.target.dataset.period);
            });
        });
    }

    async changeCurrency(newCurrency) {
        if (newCurrency === this.currentCurrency) return;

        try {
            this.showLoading();
            
            // 设置新的默认币种
            const response = await fetch('/api/v1/account/currency', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({ currency: newCurrency })
            });

            const result = await response.json();
            if (result.success) {
                this.currentCurrency = newCurrency;
                this.updateCurrencyDisplay();
                await this.updateAccountData();
                
                // 触发全局币种变化事件
                this.dispatchCurrencyChangeEvent(newCurrency);
            } else {
                throw new Error(result.message || '设置币种失败');
            }
        } catch (error) {
            console.error('切换币种失败:', error);
            this.showError('切换币种失败');
        }
    }

    changePeriod(newPeriod) {
        if (newPeriod === this.currentPeriod) return;

        this.currentPeriod = newPeriod;
        
        // 更新按钮状态
        this.elements.periodButtons.forEach(btn => {
            btn.classList.toggle('active', btn.dataset.period === newPeriod);
        });

        // 更新盈亏显示
        this.updateProfitDisplay();
    }

    async updateAccountData() {
        try {
            this.showLoading();
            
            const response = await fetch(`/api/v1/account/summary/${this.currentCurrency}`);
            const result = await response.json();
            
            if (result.success) {
                this.accountData = result.data;
                this.updateDisplay();
            } else {
                throw new Error(result.message || '获取账户数据失败');
            }
        } catch (error) {
            console.error('更新账户数据失败:', error);
            this.showError('获取数据失败');
        } finally {
            this.hideLoading();
        }
    }

    updateDisplay() {
        if (!this.accountData) return;

        // 更新总资产
        if (this.elements.totalAssets) {
            const totalAssets = parseFloat(this.accountData.balance.totalEquity);
            this.elements.totalAssets.textContent = this.formatNumber(totalAssets);
        }

        // 更新盈亏显示
        this.updateProfitDisplay();
    }

    updateProfitDisplay() {
        if (!this.accountData || !this.accountData.profitLoss) return;

        const profitData = this.accountData.profitLoss.find(p => p.period === this.currentPeriod);
        if (!profitData) return;

        // 更新盈亏金额
        if (this.elements.profitAmount) {
            const amount = parseFloat(profitData.profitAmount);
            const formattedAmount = this.formatNumber(Math.abs(amount));
            const sign = amount >= 0 ? '+' : '-';
            
            this.elements.profitAmount.textContent = `${sign}${formattedAmount}`;
            this.elements.profitAmount.className = `profit-amount ${amount >= 0 ? 'positive' : 'negative'}`;
        }

        // 更新盈亏百分比
        if (this.elements.profitPercent) {
            const percent = parseFloat(profitData.profitPercent);
            const sign = percent >= 0 ? '+' : '';
            
            this.elements.profitPercent.textContent = `(${sign}${percent.toFixed(2)}%)`;
            this.elements.profitPercent.className = `profit-percent ${percent >= 0 ? 'positive' : 'negative'}`;
        }

        // 更新时间周期显示
        if (this.elements.profitPeriod) {
            const periodNames = {
                '1d': '1日收益',
                '1w': '1周收益',
                '1m': '1月收益',
                '6m': '半年收益'
            };
            this.elements.profitPeriod.textContent = periodNames[this.currentPeriod] || '收益';
        }
    }

    formatNumber(number) {
        if (isNaN(number)) return '0.00';
        
        // 根据币种格式化数字
        switch (this.currentCurrency) {
            case 'BTC':
                return number.toFixed(8);
            case 'CNY':
            case 'USD':
            case 'USDT':
                return number.toLocaleString('en-US', {
                    minimumFractionDigits: 2,
                    maximumFractionDigits: 2
                });
            default:
                return number.toFixed(2);
        }
    }

    showLoading() {
        if (this.elements.loadingIndicator) {
            this.elements.loadingIndicator.style.display = 'block';
        }
        
        if (this.elements.totalAssets) {
            this.elements.totalAssets.innerHTML = '<span class="loading"></span>';
        }
    }

    hideLoading() {
        if (this.elements.loadingIndicator) {
            this.elements.loadingIndicator.style.display = 'none';
        }
    }

    showError(message) {
        if (this.elements.totalAssets) {
            this.elements.totalAssets.textContent = '--';
        }
        
        if (this.elements.profitAmount) {
            this.elements.profitAmount.textContent = '--';
        }
        
        if (this.elements.profitPercent) {
            this.elements.profitPercent.textContent = '--';
        }
        
        console.error('AccountCard Error:', message);
    }

    dispatchCurrencyChangeEvent(currency) {
        const event = new CustomEvent('currencyChanged', {
            detail: { currency, symbol: this.getCurrencySymbol(currency) }
        });
        window.dispatchEvent(event);
    }

    getCurrencySymbol(currency) {
        const currencyData = this.supportedCurrencies.find(c => c.currency === currency);
        return currencyData ? currencyData.symbol : '$';
    }

    // 公共方法
    getCurrentCurrency() {
        return this.currentCurrency;
    }

    async refresh() {
        await this.updateAccountData();
    }
}

// 导出组件
if (typeof module !== 'undefined' && module.exports) {
    module.exports = AccountCard;
} else {
    window.AccountCard = AccountCard;
}