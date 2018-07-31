package main

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/kirkbyers/openstall-master/sockets"
)

func main() {
	port := flag.Int("port", 4567, "port to start server on")
	flag.Parse()

	// Init sockets before http.ListenAndServe
	sockets.Init()

	http.HandleFunc("/ws", sockets.WebsocketHTTPHandler)

	fmt.Printf("Server Listening on port %v\n", *port)
	err := http.ListenAndServe(fmt.Sprintf(":%v", *port), nil)
	if err != nil {
		fmt.Println("There was a problem listing and serving", err)
	}
}
