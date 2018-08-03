package sockets

import (
	"bytes"
	"fmt"
	"time"

	"github.com/gorilla/websocket"
)

// SubClient is the middleman type connecting to Hub
// Only has access to readPump
type SubClient struct {
	hub *Hub
	// WS connection
	conn *websocket.Conn
	// Buffer channel of outbound messages
	send chan []byte
}

// readPump pumps message from the websocker connection to the hub
func (c *SubClient) readPump() {
	defer func() {
		c.hub.unregisterSub <- c
		c.conn.Close()
	}()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				fmt.Println("Error:", err)
			}
			break
		}
		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
		c.hub.broadcast <- message
	}
}
