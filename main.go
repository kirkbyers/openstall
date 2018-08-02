package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/kirkbyers/openstall-master/db"
	door "github.com/kirkbyers/openstall-master/door-monitor-example"
	"github.com/kirkbyers/openstall-master/routes"
	"github.com/kirkbyers/openstall-master/sockets"
	homedir "github.com/mitchellh/go-homedir"
)

func main() {
	var err error
	// flags
	port := flag.Int("port", 4567, "port to start server on")
	flag.Parse()

	// Init sockets before http.ListenAndServe
	sockets.Init()

	// Init DB
	var home string
	home, err = homedir.Dir()
	must(err)
	dbPath := filepath.Join(home, "openstall.db")
	must(db.Init(dbPath))

	handler := http.NewServeMux()
	handler.Handle("/ws", routes.WebsocketHandler{})
	handler.Handle("/register", routes.RegisterHandler{})

	fmt.Printf("Server Listening on port %v\n", *port)
	go http.ListenAndServe(fmt.Sprintf(":%v", *port), handler)
	door.Test()
	// if err != nil {
	// 	fmt.Println("There was a problem listing and serving", err)
	// }
	for {
	}
}

func must(err error) {
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
