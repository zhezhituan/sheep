package main

import (
	"fmt"
	"net/http"
	"sheep/config"
	handler "sheep/server/handler"
	model "sheep/server/model"
	wb "sheep/server/websocket"
	"time"
)

func init() {
	// 初始化 redis 连接池，全局唯一
	redisInfo := config.Configuration.RedisInfo
	fmt.Println("redisInfo", redisInfo)
	initRedisPool(redisInfo.MaxIdle, redisInfo.MaxActive, time.Second*(redisInfo.IdleTimeout), redisInfo.Host)

	// 创建 userDao 用于操作用户信息
	// 全局唯一 UserDao 实例：model.CurrentUserDao
	model.CurrentUserDao = model.InitUserDao(pool)
	//全局唯一web管理器
	go wb.Manager.Start()
	//全局唯一session管理

}

func main() {

	http.Handle("/", http.FileServer(http.Dir("C:/goproject/src/sheep/server/web/")))
	//http.Handle("/chatroom.html", http.FileServer(http.Dir("C:/goproject/src/sheep/server/web/chatroom.html")))
	http.HandleFunc("/login", handler.Login)
	http.HandleFunc("/register", handler.Register)
	http.HandleFunc("/chatroom", handler.Chatroom)
	http.ListenAndServe(":9090", nil) //监听8080端口
}
