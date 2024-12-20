package main

import "fmt"

// Hub maintains the set of active clients and broadcasts messages to the clients.
// Represents a single chatroom or lobby.
type Hub struct {
	// Registered clients.
	clients map[*Client]struct{}

	// Inbound messages from the clients.
	broadcast chan []byte

	// Registration requests from the clients.
	register chan *Client

	// Unregister requests from the clients.
	unregister chan *Client
}

func newHub() *Hub {
	return &Hub{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]struct{}),
	}
}

func (h *Hub) run() {
	for {
		select {
		// A client to register
		case client := <-h.register:
			h.clients[client] = struct{}{}
			// A client to unregister (peer has disconnected)
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client) // remove client from the map
				close(client.send)        // close the associated channel)
			}
		// a message to broadcast to all peers/clients
		case message := <-h.broadcast:
			fmt.Println("broadcasting message", message)
			for client := range h.clients {
				select {
				// send this message if the channel (fifo blocking queue) is not full
				// if its full we will block
				case client.send <- message: // send this message to the client's message buffer
				default:
					// if the send channel is full of messages then close the channel,
					// we don't want it to block us forever from sending messages to the rest of the peers.
					// its most likely an attack or spam
					close(client.send)        // close their channel
					delete(h.clients, client) // remove from map
				}

			}
		}
	}
}
