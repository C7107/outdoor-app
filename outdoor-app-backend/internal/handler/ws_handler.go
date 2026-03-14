package handler

import (
	"log"
	"net/http"
	"outdoor-app-backend/internal/service"
	"outdoor-app-backend/pkg/jwt"
	"outdoor-app-backend/pkg/ws"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// 创建一个结构体，只不过这个结构体内容是一个函数，Upgrade() 方法属于 Upgrader结构体，然后下面就调用
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

// WebSocket：在连接建立时手动检查 token好，不用中间件
func ConnectWebSocket(c *gin.Context) {
	tokenString := c.Query("token")
	if tokenString == "" {
		return
	}

	claims, err := jwt.ParseToken(tokenString)
	if err != nil {
		return
	}
	userID := claims.UserID

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil) //这就是上面说的
	if err != nil {
		return
	}

	// 1. 将连接交给 Hub 管理 (上线)
	ws.Hub.AddClient(userID, conn)

	// 2. 🌟 核心逻辑：建立连接时，自动获取该用户所有的【未读】消息并推送
	go func(uid uint) {
		// 查询数据库中该用户未读的消息 (这里假设用之前写的 GetUserMessages 或专门写个未读列表)
		// 我们可以复用一个查询未读的方法
		unreadMsgs, err := service.GetUnreadMessages(uid)
		if err == nil && len(unreadMsgs) > 0 {
			// 推送给该用户
			ws.SendMessageToUser(uid, map[string]interface{}{
				"type": "unread_list", // 给前端一个类型标识，告诉它这是历史未读消息
				"data": unreadMsgs,
			})
			log.Printf("用户 %d 上线，已推送 %d 条历史未读消息", uid, len(unreadMsgs))
		}
	}(userID)

	defer ws.Hub.RemoveClient(userID)
	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			break
		}
	}
}
