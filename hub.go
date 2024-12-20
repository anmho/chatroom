package main

// Hub maintains the set of active clients and broadcasts messages to the clients.
// Represents a single chatroom or lobby.
type Hub struct {
	// Registered clients.
	clients map[*Client]bool

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
		clients:    make(map[*Client]bool),
	}
}
