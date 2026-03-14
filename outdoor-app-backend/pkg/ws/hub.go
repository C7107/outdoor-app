package ws

import (
	"encoding/json"
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

// ClientManager 管理所有的 WebSocket 连接
type ClientManager struct {
	// 存储: UserID -> WebSocket 连接
	// ⚠️ 注意：实际生产中一个用户可能多端登录，这里为了毕设简化，只存一个连接
	Clients    map[uint]*websocket.Conn
	ClientsMap sync.RWMutex // 读写锁，防止高并发下 Map 崩溃
}

// 全局唯一的 Hub 实例
var Hub = &ClientManager{
	Clients: make(map[uint]*websocket.Conn),
}

// AddClient 用户上线，保存连接
func (manager *ClientManager) AddClient(userID uint, conn *websocket.Conn) {
	manager.ClientsMap.Lock()
	defer manager.ClientsMap.Unlock()
	manager.Clients[userID] = conn
	log.Printf("用户 %d 已连接 WebSocket，当前在线人数: %d", userID, len(manager.Clients))
}

// RemoveClient 用户下线，删除连接
func (manager *ClientManager) RemoveClient(userID uint) {
	manager.ClientsMap.Lock()
	defer manager.ClientsMap.Unlock()
	if conn, ok := manager.Clients[userID]; ok {
		conn.Close()
		delete(manager.Clients, userID)
		log.Printf("用户 %d WebSocket 已断开", userID)
	}
}

// SendMessageToUser 🌟 核心：给指定用户发送实时消息
func SendMessageToUser(userID uint, message interface{}) {
	Hub.ClientsMap.RLock()
	conn, ok := Hub.Clients[userID]
	Hub.ClientsMap.RUnlock()

	// 如果该用户当前在线（连着 WebSocket）
	if ok {
		// 将结构体转为 JSON 字符串
		msgBytes, _ := json.Marshal(message)
		// 发送给前端
		err := conn.WriteMessage(websocket.TextMessage, msgBytes)
		if err != nil {
			log.Printf("给用户 %d 推送消息失败: %v", userID, err)
			Hub.RemoveClient(userID) // 发送失败通常是断网了，清理掉
		}
	}
}
