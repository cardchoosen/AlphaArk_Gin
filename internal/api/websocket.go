package api

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/cardchoosen/AlphaArk_Gin/internal/config"
	"github.com/cardchoosen/AlphaArk_Gin/internal/service"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // 允许跨域
	},
}

// WebSocketManager WebSocket连接管理器
type WebSocketManager struct {
	clients      map[*websocket.Conn]bool
	broadcast    chan []byte
	register     chan *websocket.Conn
	unregister   chan *websocket.Conn
	mutex        sync.RWMutex
	priceService service.PriceService
}

// NewWebSocketManager 创建WebSocket管理器
func NewWebSocketManager(cfg *config.OKXConfig) *WebSocketManager {
	return &WebSocketManager{
		clients:      make(map[*websocket.Conn]bool),
		broadcast:    make(chan []byte),
		register:     make(chan *websocket.Conn),
		unregister:   make(chan *websocket.Conn),
		priceService: service.NewPriceService(cfg),
	}
}

// Run 启动WebSocket管理器
func (manager *WebSocketManager) Run() {
	// 启动价格数据流
	manager.priceService.StartPriceStream("BTC-USDT", func(priceData *service.PriceData) {
		data, err := json.Marshal(priceData)
		if err != nil {
			log.Printf("序列化价格数据失败: %v", err)
			return
		}
		
		select {
		case manager.broadcast <- data:
		default:
			// 如果广播通道满了，跳过这次广播
		}
	})

	for {
		select {
		case client := <-manager.register:
			manager.mutex.Lock()
			manager.clients[client] = true
			manager.mutex.Unlock()
			log.Printf("WebSocket客户端连接，当前连接数: %d", len(manager.clients))

		case client := <-manager.unregister:
			manager.mutex.Lock()
			if _, ok := manager.clients[client]; ok {
				delete(manager.clients, client)
				client.Close()
			}
			manager.mutex.Unlock()
			log.Printf("WebSocket客户端断开，当前连接数: %d", len(manager.clients))

		case message := <-manager.broadcast:
			manager.mutex.RLock()
			for client := range manager.clients {
				err := client.WriteMessage(websocket.TextMessage, message)
				if err != nil {
					log.Printf("发送消息失败: %v", err)
					delete(manager.clients, client)
					client.Close()
				}
			}
			manager.mutex.RUnlock()
		}
	}
}

// HandleWebSocket 处理WebSocket连接
func (manager *WebSocketManager) HandleWebSocket(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("WebSocket升级失败: %v", err)
		return
	}

	// 注册新客户端
	manager.register <- conn

	// 发送初始价格数据
	go func() {
		priceData, err := manager.priceService.GetPrice("BTC-USDT")
		if err != nil {
			log.Printf("获取初始价格失败: %v", err)
			return
		}

		data, err := json.Marshal(priceData)
		if err != nil {
			log.Printf("序列化初始价格数据失败: %v", err)
			return
		}

		conn.WriteMessage(websocket.TextMessage, data)
	}()

	// 监听客户端断开
	defer func() {
		manager.unregister <- conn
	}()

	// 保持连接活跃
	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket错误: %v", err)
			}
			break
		}
	}
}

// SetupWebSocketRoutes 设置WebSocket路由
func SetupWebSocketRoutes(r *gin.Engine, cfg *config.Config) {
	manager := NewWebSocketManager(&cfg.OKX)
	
	// 启动WebSocket管理器
	go manager.Run()

	// WebSocket路由
	r.GET("/ws/price", manager.HandleWebSocket)
}