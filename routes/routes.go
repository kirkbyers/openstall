package routes

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/kirkbyers/openstall-master/db"
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

// RegisterHandler handles registering Monitors
type RegisterHandler struct{}

func (RegisterHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var m db.Monitor
	d := json.NewDecoder(r.Body)
	if err := d.Decode(&m); err != nil {
		http.Error(w, "Malformed body", http.StatusBadRequest)
		return
	}
	if _, err := db.UpdateMonitor(&m); err != nil {
		fmt.Println("There was an error saving monitor to db:", err)
		return
	}
	return
}

// WSPubHandler is the handler for starting up ws client connections
type WSPubHandler struct{}

func (WSPubHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	sockets.ServeWsPub(sockets.PublicHub, w, r)
}

// WSSubHandler is the handler for starting up ws client connections
type WSSubHandler struct{}

func (WSSubHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	sockets.ServeWsSub(sockets.PublicHub, w, r)
}
