package main

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/kirkbyers/openstall-master/sockets"
)

func serveChatExample(w http.ResponseWriter, r *http.Request) {
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

func main() {
	port := flag.Int("port", 4567, "port to start server on")
	flag.Parse()

	// Init sockets before http.ListenAndServe
	sockets.Init()

	http.HandleFunc("/", serveChatExample)
	http.HandleFunc("/ws", sockets.WebsocketHTTPHandler)

	fmt.Printf("Server Listening on port %v\n", *port)
	err := http.ListenAndServe(fmt.Sprintf(":%v", *port), nil)
	if err != nil {
		fmt.Println("There was a problem listing and serving", err)
	}
}
