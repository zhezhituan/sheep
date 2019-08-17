package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	model "sheep/server/model"
	"sheep/server/util"
	wb "sheep/server/websocket"

	"github.com/gorilla/websocket"
)

type ResMsg struct {
	Code int         `json:"Code"`
	Data interface{} `json:"Data"`
}

type loginReq struct {
	Name string
	Pw   string
}

func Login(w http.ResponseWriter, r *http.Request) {
	lr := loginReq{}
	json.NewDecoder(r.Body).Decode(&lr)
	user, err := model.CurrentUserDao.Login(lr.Name, lr.Pw)
	fmt.Println(user, err)
	if err == nil {
		res := ResMsg{Code: 0, Data: user}
		str, _ := json.Marshal(res)
		sess, ok := util.GlobalSessions.SessionStart(w, r)
		if ok != nil {
			fmt.Println("Cookie Error")
		}
		sess.Set("user", user)
		w.Write([]byte(str))

	} else {
		res := ResMsg{Code: 1, Data: user}
		str, _ := json.Marshal(res)
		w.Write([]byte(str))
	}
}

type resgisterReq struct {
	Name    string
	Pw      string
	Pw_sure string
}

func Register(w http.ResponseWriter, r *http.Request) {
	rr := resgisterReq{}
	json.NewDecoder(r.Body).Decode(&rr)
	user, err := model.CurrentUserDao.Register(rr.Name, rr.Pw, rr.Pw_sure)
	fmt.Print(user, err)
	if err == nil {
		res := ResMsg{Code: 0, Data: user}
		str, _ := json.Marshal(res)
		w.Write([]byte(str))
	} else {
		res := ResMsg{Code: 1, Data: user}
		str, _ := json.Marshal(res)
		w.Write([]byte(str))
	}
}

func Chatroom(w http.ResponseWriter, r *http.Request) {
	conn, error := (&websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}).Upgrade(w, r, nil)
	if error != nil {
		http.NotFound(w, r)
		return
	}
	c, err := r.Cookie("gosessionid")
	if err != nil {
		fmt.Print("cookie error")
		return
	}
	client := &(wb.Client{Id: c.Value, Socket: conn, Send: make(chan []byte)})
	wb.Manager.Register <- client
	go client.Read()
	go client.Write()
}
