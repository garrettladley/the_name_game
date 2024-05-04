package handlers

import (
	"log"
	"time"

	"github.com/gofiber/contrib/websocket"
)

func Join(c *websocket.Conn) {
	defer c.Close()

	log.Println(c.Params("id"))

	timer := time.NewTimer(5 * time.Second)

	<-timer.C
	c.WriteMessage(websocket.TextMessage, []byte("Goodbye!"))
}
