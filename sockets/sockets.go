package sockets

// PublicHub is the exported instance
var PublicHub *Hub

// Init starts the server WS hub
func Init() {
	PublicHub = newHub()
	go PublicHub.run()
}
