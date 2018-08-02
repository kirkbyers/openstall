package sockets

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/kirkbyers/openstall-master/db"
)

// Hub maintains the set of active cleints and broadcasts message to the clients
type Hub struct {
	// Registered Clients
	clients map[*Client]bool
	// Inbound msgs from clients
	broadcast chan []byte
	// Register reqs from clients
	register chan *Client
	// Unregister requests from clients
	unregister chan *Client
}

func newHub() *Hub {
	return &Hub{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
		case message := <-h.broadcast:
			// Decode JSON to ensure message is correct structure
			var m db.Monitor
			d := json.NewDecoder(bytes.NewBuffer(message))
			if err := d.Decode(&m); err != nil {
				// if message cannot be decoded as a monitor
				fmt.Println("There was a message not matching type Monitor sent:", err)
			} else {
				// else update the db entry
				if err := db.UpdateExistingMonitor(&m); err != nil {
					fmt.Println("There was an issue updating monitor in DB from message:", err)
				}
			}
		}
	}
}
