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

// WebsocketHTTPHandler is the handler for starting up ws client connections
func WebsocketHTTPHandler(w http.ResponseWriter, r *http.Request) {
	serveWs(hub, w, r)
}
