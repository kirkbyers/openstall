package sockets

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/kirkbyers/openstall-master/db"
)

// Hub maintains the set of active cleints and broadcasts message to the clients
type Hub struct {
	// Inbound msgs from clients
	broadcast chan []byte
	// Registered PubClients
	pubClients map[*PubClient]bool
	// registerPub reqs from pub clients
	registerPub chan *PubClient
	// Unregister requests from pub clients
	unregisterPub chan *PubClient
	// Registered SubClients
	subClients map[*SubClient]bool
	// registerPub reqs from pub clients
	registerSub chan *SubClient
	// Unregister requests from pub clients
	unregisterSub chan *SubClient
}

func newHub() *Hub {
	return &Hub{
		broadcast:     make(chan []byte),
		registerPub:   make(chan *PubClient),
		unregisterPub: make(chan *PubClient),
		pubClients:    make(map[*PubClient]bool),
		registerSub:   make(chan *SubClient),
		unregisterSub: make(chan *SubClient),
		subClients:    make(map[*SubClient]bool),
	}
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.registerPub:
			h.pubClients[client] = true
		case client := <-h.registerSub:
			h.subClients[client] = true
		case client := <-h.unregisterPub:
			if _, ok := h.pubClients[client]; ok {
				delete(h.pubClients, client)
				close(client.send)
			}
		case client := <-h.unregisterSub:
			if _, ok := h.subClients[client]; ok {
				delete(h.subClients, client)
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
				for client := range h.subClients {
					select {
					case client.send <- message:
					default:
						close(client.send)
						delete(h.subClients, client)
					}
				}
			}
		}
	}
}
