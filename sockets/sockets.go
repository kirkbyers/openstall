package sockets

import (
	"net/http"
)

var hub *Hub

// Init starts the server WS hub
func Init() {
	hub = newHub()
	go hub.run()
}

// WebsocketHandler is the handler for starting up ws client connections
type WebsocketHandler struct{}

func (w WebsocketHandler) ServeHTTP(wr http.ResponseWriter, r *http.Request) {
	serveWs(hub, wr, r)
}
