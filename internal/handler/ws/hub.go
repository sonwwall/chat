package ws

import (
	"chat/internal/model"
	"github.com/gorilla/websocket"
	"log"
	"sync"
)

type Client struct {
	Conn   *websocket.Conn // Websocket连接对象
	UserID uint
	RoomID uint
}

type Hub struct {
	Clients    map[*Client]bool   //客户端集合，使用Map存储
	Broadcast  chan model.Message //广播消息通道
	Register   chan *Client       //注册通道
	Unregister chan *Client       //注销通道
	Mutex      sync.Mutex         //并发保护锁
}

var HubInstance = Hub{
	Clients:    make(map[*Client]bool),
	Broadcast:  make(chan model.Message),
	Register:   make(chan *Client),
	Unregister: make(chan *Client),
}

func (h *Hub) Run() {

	for {
		select {
		case client := <-h.Register: //触发条件：当有新客户端注册时
			h.Mutex.Lock()
			h.Clients[client] = true //将新客户端加入map中
			h.Mutex.Unlock()
		case client := <-h.Unregister:
			h.Mutex.Lock()
			if _, ok := h.Clients[client]; ok {
				err := client.Conn.Close() //关闭WebSocket连接
				if err != nil {
					log.Println(err.Error())
				}
				delete(h.Clients, client) //将该client从map中移除
			}
			h.Mutex.Unlock()
		case message := <-h.Broadcast:
			h.Mutex.Lock()
			for client := range h.Clients {
				if client.RoomID == message.RoomID { //只发送给目标房间
					err := client.Conn.WriteJSON(message) //发送消息
					if err != nil {
						client.Conn.Close()
						delete(h.Clients, client)
						log.Println(err.Error())
					}
				}
			}
			h.Mutex.Unlock()

		}
	}
}
