package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"
	"github.com/kirkbyers/openstall-master/db"
	"gobot.io/x/gobot/drivers/gpio"
	"gobot.io/x/gobot/platforms/raspi"
)

func main() {
	// Flags for config
	mID := flag.String("id", "", "String for id of monitor")
	mName := flag.String("name", "", "String for display-able name")
	masterHost := flag.String("master-host", "", "String for host name of master address")
	masterPort := flag.Int("master-port", 0, "Int for port of master address")
	flag.Parse()
	if *mID == "" || *mName == "" || *masterHost == "" || *masterPort == 0 {
		fmt.Println("--id, --name, --master-host, --master-port flags required")
		os.Exit(1)
	}

	// Attach to rPi GPIO pin readout
	r := raspi.NewAdaptor()
	pin := gpio.NewDirectPinDriver(r, "7")

	// Get initial reading on pin
	pinVal, err := pin.DigitalRead()
	must(err)
	// Build the current status of the monitor
	mStatus := doorMonitorStatus(pinVal)

	// Build the monitor
	m := &db.Monitor{
		ID:     *mID,
		Name:   *mName,
		Type:   "door",
		Status: mStatus,
	}
	// Register as monitor
	regURL := url.URL{Scheme: "http", Host: fmt.Sprintf("%s:%d", *masterHost, *masterPort), Path: "/register"}
	mJSON, err := json.Marshal(m)
	if err != nil {
		fmt.Println("There was an err marshling JSON", err)
		os.Exit(1)
	}
	http.Post(regURL.String(), "json", bytes.NewBuffer(mJSON))

	// Create dial-up URL at server address
	u := url.URL{Scheme: "ws", Host: "localhost:4567", Path: "/pub"}
	// Create ws *Conn
	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		fmt.Println("There was a problem dialing WS server:", err)
	}
	defer c.Close()

	// Set up ticker to check pin on an interval
	ticker := time.NewTicker(250 * time.Millisecond)
	defer ticker.Stop()

	// Handle interuptions
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	for {
		// Never ending loop to listen to ticker to check pin out
		select {
		case <-ticker.C:
			// Get pin readout
			pinVal, err = pin.DigitalRead()
			must(err)
			// Build the current status of the monitor
			mStatus := doorMonitorStatus(pinVal)
			if mStatus != m.Status {
				m.Status = mStatus
				err = sendMonitorUpdate(c, m)
				if err != nil {
					c.Close()
					c, _, err = websocket.DefaultDialer.Dial(u.String(), nil)
					if err != nil {
						fmt.Println("There was a problem dialing WS server:", err)
					}
				}
			}
		case <-interrupt:
			log.Println("interrupt")

			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
				return
			}
		}
	}
}

func sendMonitorUpdate(c *websocket.Conn, m *db.Monitor) error {
	// Send monitor as Text Message to WS conn
	var err error
	// Convert the monitor to JSON
	// Write the monitor Update
	err = c.WriteJSON(&m)
	if err != nil {
		fmt.Println("There was an issue writting JSON")
		return err
	}
	return nil
}

func doorMonitorStatus(s int) string {
	// Convert pin readout to string
	var mStatus string
	if s == 1 {
		mStatus = "open"
	} else {
		mStatus = "close"
	}
	return mStatus
}

func must(err error) {
	// Generic error handling
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
