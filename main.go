package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// WebSocket upgrader , upgrade html request to websocket
var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1 << 10,
		WriteBufferSize: 1 << 10,
	}
)

func homepage(c *gin.Context) {
	c.HTML(http.StatusOK, "index.tmpl", gin.H{
		"title": "Home Page",
	})
}

func reader(conn *websocket.Conn) {
	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}

		for {
			log.Println(string(p))
			if err := conn.WriteMessage(messageType, p); err != nil {
				log.Println(err)
			}
			time.Sleep(3 * time.Second)
		}
	}

}

func wsEndpoint(c *gin.Context) {
	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}

	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println(err)
	}
	log.Println("Client Successfully connected")

	go reader(ws)
	//c.String(http.StatusOK, "web socket")
}

func main() {
	fmt.Println("gin-websocket")

	route := gin.Default()
	route.LoadHTMLGlob("templates/*")
	route.GET("/", homepage)
	route.GET("/ws", wsEndpoint)

	route.Run(":8080")
}
