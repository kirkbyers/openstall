package sockets

import (
	"bytes"
	"fmt"
	"time"

	"github.com/gorilla/websocket"
)

// PubClient is the middleman type connecting to Hub
// Only has access to readPump
type PubClient struct {
	hub *Hub
	// WS connection
	conn *websocket.Conn
	// Buffer channel of outbound messages
	send chan []byte
}

// readPump pumps message from the websocker connection to the hub
func (c *PubClient) readPump() {
	defer func() {
		c.hub.unregisterPub <- c
		c.conn.Close()
	}()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	c.conn.SetWriteDeadline(time.Now().Add(writeWait))
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
