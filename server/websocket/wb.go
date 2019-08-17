package websocket

import (
	"encoding/json"
	"sheep/server/model"
	"sheep/server/util"

	"github.com/gorilla/websocket"
)

type ClientManager struct {
	Clients    map[string]*Client
	User       map[string]string
	Broadcast  chan []byte
	Register   chan *Client
	Unregister chan *Client
}

type Client struct {
	Id     string
	Socket *websocket.Conn
	Send   chan []byte
}

type Message struct {
	Sender    string `json:"sender"`
	Recipient string `json:"recipient"`
	Content   string `json:"content"`
}

//全局唯一
var Manager = ClientManager{
	Broadcast:  make(chan []byte),
	Register:   make(chan *Client),
	Unregister: make(chan *Client),
	Clients:    make(map[string]*Client),
	User:       make(map[string]string),
}

func (Manager *ClientManager) Start() {

	for {
		select {
		case conn := <-Manager.Register:
			Manager.Clients[conn.Id] = conn
			sess, _ := util.GlobalSessions.GetSessionStore(conn.Id)
			user := sess.Get("user")
			tempu := user.(model.User)
			Manager.User[tempu.Name] = conn.Id
			println(conn.Id, tempu.Name)
			//Manager.Clients[conn] = true
			jsonMessage, _ := json.Marshal(&Message{Content: "/A new socket has connected."})
			Manager.Send(jsonMessage, conn)
		case conn := <-Manager.Unregister:
			if _, ok := Manager.Clients[conn.Id]; ok {
				close(conn.Send)
				delete(Manager.Clients, conn.Id)
				sess, _ := util.GlobalSessions.GetSessionStore(conn.Id)
				user := sess.Get("user")
				tempu := user.(model.User)
				delete(Manager.User, tempu.Name)
				jsonMessage, _ := json.Marshal(&Message{Content: "/A socket has disconnected."})
				Manager.Send(jsonMessage, conn)
			}
		case message := <-Manager.Broadcast:
			for _, conn := range Manager.Clients {
				select {
				case conn.Send <- message:
				default:
					close(conn.Send)
					delete(Manager.Clients, conn.Id)
				}
			}
		}
	}
}

func (Manager *ClientManager) Send(message []byte, ignore *Client) {
	for _, conn := range Manager.Clients {
		if conn != ignore {
			conn.Send <- message
		}
	}
}

func (c *Client) Read() {
	defer func() {
		Manager.Unregister <- c
		c.Socket.Close()
	}()

	for {
		_, mes, err := c.Socket.ReadMessage()
		println(string(mes))
		if err != nil {
			println("Read Message Error")
			Manager.Unregister <- c
			c.Socket.Close()
			break
		}
		temp_mes := Message{}
		err = json.Unmarshal(mes, &temp_mes)
		// sess, _ := util.GlobalSessions.GetSessionStore(c.Id)
		// user := sess.Get("user")
		// tempu := user.(model.User)
		// temp_mes :=user
		if temp_mes.Recipient == "all" {
			Manager.Broadcast <- mes
		} else {
			recName := temp_mes.Recipient
			println(recName)
			recId := Manager.User[recName]
			println(recName, recId)
			conn := Manager.Clients[recId]
			if conn == nil {
				c.Send <- []byte("无此人")
			} else {
				conn.Send <- mes
			}
		}
	}
}

func (c *Client) Write() {
	defer func() {
		c.Socket.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			if !ok {
				c.Socket.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			c.Socket.WriteMessage(websocket.TextMessage, message)
		}
	}
}
