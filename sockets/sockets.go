package sockets

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the perr
	writeWait = 10 * time.Second
	// Time allowed to read the next pong message from the peer
	pongWait = 60 * time.Second
	// Send pings to the peer with this period. Must be < pongWait
	pingPeriod = (pongWait * 9) / 10
	// Max message size allowed from peer
	maxMessageSize = 512
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// PublicHub is the exported instance
var PublicHub *Hub

// Init starts the server WS hub
func Init() {
	PublicHub = newHub()
	go PublicHub.run()
}

// ServeWsPub handles websock requests from the peer
func ServeWsPub(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	client := &PubClient{hub: hub, conn: conn, send: make(chan []byte, 256)}
	client.hub.registerPub <- client

	// Allow collection of memory referenced by the caller by doing all work in new goroutines
	go client.readPump()
}

// ServeWsSub handles websock requests from the peer
func ServeWsSub(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	client := &SubClient{hub: hub, conn: conn, send: make(chan []byte, 256)}
	client.hub.registerSub <- client

	// Allow collection of memory referenced by the caller by doing all work in new goroutines
	go client.writePump()
}
