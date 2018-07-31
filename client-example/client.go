package main

import (
	"fmt"
	"net/url"
	"time"

	"github.com/gorilla/websocket"
)

func main() {
	// Create dial-up URL at server address
	u := url.URL{Scheme: "ws", Host: "localhost:4567", Path: "/ws"}

	// Create c *Conn
	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		fmt.Println("There was a problem dialing WS server:", err)
	}
	defer c.Close()

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case t := <-ticker.C:
			err := c.WriteMessage(websocket.TextMessage, []byte(t.String()))
			if err != nil {
				fmt.Println("There was a problem writting message:", err)
				return
			}
		}
	}
}
