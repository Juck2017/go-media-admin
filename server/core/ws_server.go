package core

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type WSServer struct {
	Port int `json:"port"`
	conn *websocket.Conn
}

// 升级websocket连接
var upgrader = websocket.Upgrader{
	ReadBufferSize:   1024,
	WriteBufferSize:  1024,
	HandshakeTimeout: 5 * time.Second,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (wss *WSServer) Start(handler func(*gin.Context)) {
	gin.SetMode(gin.ReleaseMode)
	gin.DefaultWriter = ioutil.Discard
	r := gin.Default()
	r.GET("/", handler)
	err := r.Run(fmt.Sprintf("0.0.0.0:%d", wss.Port))
	if err != nil {
		panic(err)
	}
}
