package routes

import (
	"fmt"
	"net/http"

	"github.com/kirkbyers/openstall-master/sockets"
)

// Init wires up all http routes
func Init() {

}

// ChatExampleHandler handles http req for serving the example chat app
type ChatExampleHandler struct{}

func (chat ChatExampleHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.URL)
	if r.URL.Path != "/" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(w, r, "chat-example.html")
}

// WebsocketHandler is the handler for starting up ws client connections
type WebsocketHandler struct{}

func (w WebsocketHandler) ServeHTTP(wr http.ResponseWriter, r *http.Request) {
	// if r.Method != "POST" {
	// 	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	// 	return
	// }
	// var m db.Monitor
	// d := json.NewDecoder(r.Body)
	// if err := d.Decode(&m); err != nil {
	// 	// w.Write()
	// 	return
	// }
	sockets.ServeWs(sockets.PublicHub, wr, r)
}
