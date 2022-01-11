package service

import (
	"chat/cache"
	"chat/e"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	"net/http"
	"time"
)

const month = 60 * 60 * 24 * 30 // 30天

type SendMsg struct {
	Type    int    `json:"type"`
	Content string `json:"content"`
}

type ReplyMsg struct {
	From    string `json:"from"`
	Code    int    `json:"code"`
	Content string `json:"content"`
}

type Client struct {
	ID     string
	SendID string
	Socket *websocket.Conn
	Send   chan []byte
}

type Broadcast struct {
	Client  *Client
	Message []byte
	Type    int
}

type ClientManager struct {
	Clients    map[string]*Client
	Broadcast  chan *Broadcast
	Reply      chan *Client
	Register   chan *Client
	Unregister chan *Client
}

type Message struct {
	Sender    string `json:"sender,omitempty"`
	Recipient string `json:"recipient,omitempty"`
	Content   string `json:"content,omitempty"`
}

var Manager = ClientManager{
	Clients:    make(map[string]*Client), // 参与连接的用户，出于性能的考虑，需要设置最大连接数
	Broadcast:  make(chan *Broadcast),
	Register:   make(chan *Client),
	Reply:      make(chan *Client),
	Unregister: make(chan *Client),
}

func CreateID(uid, toUid string) string {
	return uid + "->" + toUid // 1->2
}

func Handler(c *gin.Context) {
	uid := c.Query("uid")
	toUid := c.Query("toUid")
	// 升级成ws协议
	conn, err := (&websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}).
		Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		http.NotFound(c.Writer, c.Request)
		return
	}
	// 创建用户实例
	client := &Client{
		ID:     CreateID(uid, toUid), // 1->2
		SendID: CreateID(toUid, uid), // 2->1
		Socket: conn,
		Send:   make(chan []byte),
	}
	// 用户注册到用户管理上
	Manager.Register <- client

	go client.Read()
	go client.Write()
}

func (c *Client) Read() {
	// 避免资源浪费
	defer func() {
		// 结束的时候把信息传进去
		Manager.Unregister <- c // 用户进行离线操作
		_ = c.Socket.Close()
	}()
	for {
		c.Socket.PongHandler()
		sendMsg := new(SendMsg)
		err := c.Socket.ReadJSON(&sendMsg)
		if err != nil {
			log.Println("数据格式不正确", err)
			Manager.Unregister <- c
			_ = c.Socket.Close()
			break
		}
		// 发送消息 1->2
		if sendMsg.Type == 1 {
			r1, _ := cache.RedisClient.Get("c.ID").Result()     // 1->2
			r2, _ := cache.RedisClient.Get("c.SendID").Result() // 2->1
			// r1给r2发消息超过3条为回应，则不允许1再次发送
			if r1 > "3" && r2 == "" {
				replyMsg := ReplyMsg{
					Code:    e.WebsocketLimit,
					Content: "达到限制",
				}
				msg, _ := json.Marshal(replyMsg) // 序列化
				_ = c.Socket.WriteMessage(websocket.TextMessage, msg)
				continue
			} else {
				cache.RedisClient.Incr(c.ID)
				// 防止过快”分手“，建立连接三个月过期
				_, _ = cache.RedisClient.Expire(c.ID, time.Hour*24*30*3).Result()
			}
			Manager.Broadcast <- &Broadcast{
				Client:  c,
				Message: []byte(sendMsg.Content), // 发送过来的消息
			}
		}
	}
}

func (c *Client) Write() {
	defer func() {
		_ = c.Socket.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send: // 把发送的消息读出来
			if !ok {
				// 不成功，关闭消息
				_ = c.Socket.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			// 构建消息回复
			replyMsg := ReplyMsg{
				Code:    e.WebsocketSuccessMessage,
				Content: fmt.Sprintf("%s", string(message)),
			}
			msg, _ := json.Marshal(replyMsg)
			_ = c.Socket.WriteMessage(websocket.TextMessage, msg)
		}
	}

}
