package service

import (
	"chat/conf"
	"chat/e"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
)

func (manager *ClientManager) Start() {
	for {
		fmt.Println("-----监听管道通信-----")
		select {
		case conn := <-Manager.Register:
			fmt.Printf("有新连接 %s", conn.ID)
			Manager.Clients[conn.ID] = conn // 把连接放到用户管理上
			replyMsg := ReplyMsg{
				Code:    e.WebsocketSuccess,
				Content: "已经连接到服务器了",
			}
			msg, _ := json.Marshal(replyMsg)
			_ = conn.Socket.WriteMessage(websocket.TextMessage, msg)
		case conn := <-Manager.Unregister:
			fmt.Printf("连接失败 %s", conn.ID)
			if _, ok := Manager.Clients[conn.ID]; ok {
				replyMsg := &ReplyMsg{
					Code:    e.WebsocketEnd,
					Content: "连接中断",
				}
				msg, _ := json.Marshal(replyMsg)
				_ = conn.Socket.WriteMessage(websocket.TextMessage, msg)
				close(conn.Send)
				delete(Manager.Clients, conn.ID)
			}
		case broadcast := <-Manager.Broadcast: // 1->2 1发送给2
			message := broadcast.Message
			sendId := broadcast.Client.SendID // 2->1 2接收1
			flag := false                     // 默认对方是不在线的
			for id, conn := range Manager.Clients {
				if id != sendId {
					continue
				}
				select {
				case conn.Send <- message:
					flag = true
				default:
					close(conn.Send)
					delete(manager.Clients, conn.ID)
				}
			}
			id := broadcast.Client.ID // 1->2
			if flag {
				replMsg := &ReplyMsg{
					Code:    e.WebsocketOnlineReply,
					Content: "对方在线应答",
				}
				msg, _ := json.Marshal(replMsg)
				_ = broadcast.Client.Socket.WriteMessage(websocket.TextMessage, msg)
				err := InsertMsg(conf.MongoDBName, id, string(message), 1, int64(3*month)) // 1 已经读了
				if err != nil {
					fmt.Println("InsertOne Err", err)
				}
			} else {
				replMsg := &ReplyMsg{
					Code:    e.WebsocketOnlineReply,
					Content: "对方未在线应答",
				}
				msg, _ := json.Marshal(replMsg)
				_ = broadcast.Client.Socket.WriteMessage(websocket.TextMessage, msg)
				err := InsertMsg(conf.MongoDBName, id, string(message), 0, int64(3*month))
				if err != nil {
					fmt.Println("InsertOne Err", err)
				}
			}
		}
	}
}
