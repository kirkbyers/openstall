package door

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

// Test registers the test as door monitor and
// writes the current time every second to the ws conenction
func Test() {
	// Register as monitor
	regURL := url.URL{Scheme: "http", Host: "localhost:4567", Path: "/register"}
	m := &db.Monitor{
		ID:     "test-01",
		Name:   "test-monit-01",
		Type:   "door",
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

	// Start ticker for every other second
	ticker := time.NewTicker(2 * time.Second)
	tracker := 0
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			var err error
			if tracker == 0 {
				m.Status = "open"
			} else {
				m.Status = "close"
			}
			mJSON, err = json.Marshal(m)
			if err != nil {
				fmt.Println("There was a problem converting monitor to JSON:", err)
			}
			// Write the ticker time every message
			err = c.WriteMessage(websocket.TextMessage, mJSON)
			if err != nil {
				fmt.Println("There was a problem writting message:", err)
				return
			}
			tracker = (tracker + 1) % 2
		}
	}
}
