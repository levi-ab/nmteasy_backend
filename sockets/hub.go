package sockets

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"nmteasy_backend/models"
	"strings"
	"time"
)

const INFO string = "info"
const ANSWER string = "answer"
const QUESTION string = "question"

type Room struct {
	Connections map[*websocket.Conn]bool `json:"-"`
	GameState   GameState                `json:"gameState"`
}

type GameState struct {
	Questions         []models.Question `json:"questions"`
	CurrentIndex      int               `json:"currentIndex"`
	ClientRightCounts map[uuid.UUID]int `json:"clientRightCounts"`
}

type Hub struct {
	// Registered clients.
	clients map[*Client]bool

	// Mapping of user IDs to their respective clients.
	userClients map[uuid.UUID]*Client

	//Match Making queue
	matchmakingQueue []*Client

	rooms map[string]Room

	// Inbound messages from the clients.
	broadcast chan *Message // Use a Message type to include sender and target information.

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client
}

type Message struct {
	Message     string
	MessageType string
}

func NewHub() *Hub {
	return &Hub{
		broadcast:        make(chan *Message),
		register:         make(chan *Client),
		unregister:       make(chan *Client),
		clients:          make(map[*Client]bool),
		userClients:      make(map[uuid.UUID]*Client),
		matchmakingQueue: make([]*Client, 0),
		rooms:            make(map[string]Room),
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
				h.removeFromRooms(client)
				h.removeFromMatchmakingQueue(client)
				delete(h.clients, client)
				delete(h.userClients, client.clientID)
				close(client.send)
			}
		case message := <-h.broadcast:
			//here we accept answers and send the next questions???
			println(message)

		case <-time.After(time.Second * 5): // Adjust the interval as needed
			if len(h.matchmakingQueue) >= 2 {
				// Take the first two clients from the queue
				client1 := h.matchmakingQueue[0]
				client2 := h.matchmakingQueue[1]

				// Remove them from the queue
				h.matchmakingQueue = h.matchmakingQueue[2:]

				// Notify clients that they have found a match
				client1.send <- []byte("Match found! i will give u an id of user to communicate with")
				client2.send <- []byte("Match found! i will give u an id of user to communicate with")

				client2.IsInQueue = false
				client1.IsInQueue = false
				client1.targetUserID = client2.clientID
				client2.targetUserID = client1.clientID

				room := client1.clientID.String() + client2.clientID.String()

				connections := make(map[*websocket.Conn]bool)
				connections[client1.conn] = true
				connections[client2.conn] = true

				h.rooms[room] = Room{
					Connections: connections,
					GameState:   GameState{},
				}

				//then here we query the questions and send the first question
			}
		}
	}
}

func (h *Hub) removeFromRooms(client *Client) {
	// Create a new map for updated rooms
	updatedRooms := make(map[string]Room)

	// Iterate over each room
	for roomID, clients := range h.rooms {
		// Check if the client is in the room
		if !strings.Contains(roomID, client.clientID.String()) {
			// If the room ID doesn't contain the client's clientID, include it in the updated map
			updatedRooms[roomID] = clients
		} else {
			for otherClientConn := range clients.Connections {
				if otherClientConn != client.conn {
					message := "The other user has disconnected."

					messageToSend := Message{
						Message:     message,
						MessageType: INFO,
					}

					jsonMessage, _ := json.Marshal(messageToSend)

					otherClientConn.WriteMessage(websocket.TextMessage, jsonMessage)
					otherClientConn.Close()
					break // Assuming there's only one other client in the room
				}
			}
		}
	}

	// Replace the original rooms with the updated map
	h.rooms = updatedRooms
}

func (h *Hub) removeFromMatchmakingQueue(client *Client) {
	// Create a new queue without the client
	var updatedQueue []*Client
	for _, c := range h.matchmakingQueue {
		if c.clientID != client.clientID {
			updatedQueue = append(updatedQueue, c)
		}
	}
	h.matchmakingQueue = updatedQueue
}
