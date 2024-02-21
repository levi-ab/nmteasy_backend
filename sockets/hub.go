package sockets

import "github.com/google/uuid"

type Hub struct {
	// Registered clients.
	clients map[*Client]bool

	// Mapping of user IDs to their respective clients.
	userClients map[uuid.UUID]*Client

	// Inbound messages from the clients.
	broadcast chan *Message // Use a Message type to include sender and target information.

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client
}

type Message struct {
	SenderID uuid.UUID
	TargetID uuid.UUID
	Message  []byte
}

func NewHub() *Hub {
	return &Hub{
		broadcast:   make(chan *Message),
		register:    make(chan *Client),
		unregister:  make(chan *Client),
		clients:     make(map[*Client]bool),
		userClients: make(map[uuid.UUID]*Client),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
			h.userClients[client.clientID] = client
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				delete(h.userClients, client.clientID)
				close(client.send)
			}
		case message := <-h.broadcast:
			if targetClient, ok := h.userClients[message.TargetID]; ok {
				select {
				case targetClient.send <- message.Message:
				default:
					close(targetClient.send)
					delete(h.clients, targetClient)
					delete(h.userClients, targetClient.clientID)
				}
			}
		}
	}
}
