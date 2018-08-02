package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/gorilla/websocket"
	"github.com/kirkbyers/openstall-master/db"
)

func main() {
	// Register as monitor
	regURL := url.URL{Scheme: "http", Host: "localhost:4567", Path: "/register"}
	m := &db.Monitor{
		ID:     "test-00",
		Name:   "test-monit-00",
		Type:   "test",
		Status: "open",
	}
	mJSON, err := json.Marshal(m)
	if err != nil {
		fmt.Println("There was an err marshling JSON", err)
		os.Exit(1)
	}
	http.Post(regURL.String(), "json", bytes.NewBuffer(mJSON))

	// Create dial-up URL at server address
	u := url.URL{Scheme: "ws", Host: "localhost:4567", Path: "/ws"}
	// Create ws *Conn
	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		fmt.Println("There was a problem dialing WS server:", err)
	}
	defer c.Close()

	// Start ticker for every second
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case t := <-ticker.C:
			// Write the ticker time every message
			err := c.WriteMessage(websocket.TextMessage, []byte(t.String()))
			if err != nil {
				fmt.Println("There was a problem writting message:", err)
				return
			}
		}
	}
}
