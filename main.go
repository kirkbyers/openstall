package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/kirkbyers/openstall-master/db"
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
	handler.Handle("/pub", routes.WSPubHandler{})
	handler.Handle("/sub", routes.WSSubHandler{})
	handler.Handle("/register", routes.RegisterHandler{})
	handler.Handle("/", routes.StaticHandler{Path: "/", HTMLDoc: "doors.html"})
	handler.Handle("/status", routes.MonitorStatusHandler{})

	fmt.Printf("Server Listening on port %v\n", *port)
	go http.ListenAndServe(fmt.Sprintf(":%v", *port), handler)
	// door.Test() // Test that toggles the status of a monitor every 2s
	// client.Test() // There was a message not matching type Monitor sent: json: cannot unmarshal number into Go value of type db.Monitor
	for {
	}
}

func must(err error) {
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
