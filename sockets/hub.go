package sockets

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"nmteasy_backend/models"
	"nmteasy_backend/models/migrated_models"
	"strings"
	"time"
)

const INFO string = "info"
const ERROR string = "error"
const ANSWER string = "answer"
const SKIP_QUESTION string = "skip_question"
const QUESTION string = "question"
const GET_NEXT_QUESTION string = "get_next_question"
const RESULT string = "result"
const FINISHED string = "finished"
const MATCH_FOUND string = "match_found"

type Room struct {
	Clients   map[*Client]bool `json:"-"`
	GameState GameState        `json:"gameState"`
}

type GameState struct {
	Questions         []migrated_models.HistoryQuestion `json:"questions"`
	CurrentIndex      int                               `json:"currentIndex"`
	ClientRightCounts map[uuid.UUID]int                 `json:"clientRightCounts"`
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
	RoomID      string
	MessageType string
}

type AnswerMessage struct {
	Answer string
	UserID uuid.UUID
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
			if message.MessageType == ANSWER {
				room, ok := h.rooms[message.RoomID]
				if !ok {
					// Room not found
					fmt.Println("room not found")
					continue
				}

				var parsedAnswerMessage AnswerMessage

				if err := json.Unmarshal([]byte(message.Message), &parsedAnswerMessage); err != nil {
					fmt.Println("failed to parse the answer")
					continue
				}

				var sender *Client
				var anotherClient *Client
				//i dont like how i do this, it is very confusing, spend like 5 min to understand wtf is going on todo
				for client := range room.Clients {
					if client.clientID == parsedAnswerMessage.UserID {
						sender = client
					} else {
						anotherClient = client
					}
				}

				if sender == nil {
					fmt.Println("failed to determine the sender")
					continue
				}

				if room.GameState.Questions[room.GameState.CurrentIndex].RightAnswer == parsedAnswerMessage.Answer {
					if room.GameState.CurrentIndex+1 == len(room.GameState.Questions) {
						//finishing the game
						msg := Message{
							Message:     "Congratulation the game is finished",
							MessageType: FINISHED,
						}

						messageToSend, err := json.Marshal(msg)
						if err != nil {
							fmt.Println("Error marshaling message:", err)
							continue
						}

						sender.send <- messageToSend
						anotherClient.send <- messageToSend

						h.unregister <- sender
						h.unregister <- anotherClient
						continue
					}

					room.GameState.CurrentIndex = room.GameState.CurrentIndex + 1
					h.rooms[message.RoomID] = room

					msg := Message{
						Message:     "correct",
						MessageType: RESULT,
					}

					messageToSend, err := json.Marshal(msg)
					if err != nil {
						fmt.Println("Error marshaling message:", err)
						continue
					}

					sender.send <- messageToSend
					msg.Message = "other_answered"

					messageToSend, err = json.Marshal(msg)
					if err != nil {
						fmt.Println("Error marshaling message:", err)
						continue
					}

					anotherClient.send <- messageToSend //sending message when another user answered right

					sendQuestion(room.GameState.Questions[room.GameState.CurrentIndex], room.Clients)

				} else {
					msg := Message{
						Message:     "wrong",
						MessageType: RESULT,
					}

					messageToSend, err := json.Marshal(msg)
					if err != nil {
						fmt.Println("Error marshaling message:", err)
						continue
					}
					sender.send <- messageToSend
				}
			}

			if message.MessageType == SKIP_QUESTION {
				room, ok := h.rooms[message.RoomID]
				if !ok {
					// Room not found
					fmt.Println("room not found")
					continue
				}
				room.GameState.CurrentIndex = room.GameState.CurrentIndex + 1
				h.rooms[message.RoomID] = room

				sendQuestion(room.GameState.Questions[room.GameState.CurrentIndex], room.Clients)
			}

		case <-time.After(time.Second * 5): // Adjust the interval as needed
			if len(h.matchmakingQueue) >= 2 {
				// Take the first two clients from the queue
				client1 := h.matchmakingQueue[0]
				client2 := h.matchmakingQueue[1]

				// Remove them from the queue
				h.matchmakingQueue = h.matchmakingQueue[2:]

				room := client1.clientID.String() + client2.clientID.String()

				// Notify clients that they have found a match
				msg := Message{
					Message:     room,
					MessageType: MATCH_FOUND,
				}

				messageToSend, err := json.Marshal(msg)
				if err != nil {
					fmt.Println("Error marshaling message:", err)
					continue
				}
				client1.send <- messageToSend
				client2.send <- messageToSend

				client2.IsInQueue = false
				client1.IsInQueue = false

				connections := make(map[*Client]bool)
				connections[client1] = true
				connections[client2] = true

				//then here we query the questions and send the first question
				var questions []migrated_models.HistoryQuestion
				models.DB.Limit(20).Find(&questions)

				h.rooms[room] = Room{
					Clients: connections,
					GameState: GameState{
						Questions:         questions,
						ClientRightCounts: make(map[uuid.UUID]int),
					},
				}

				sendQuestion(questions[0], connections)
			}
		}
	}
}

func sendQuestion(question migrated_models.HistoryQuestion, clients map[*Client]bool) {
	stringQuestion, _ := json.Marshal(question)

	for client := range clients {
		messageToSend, _ := json.Marshal(Message{
			Message:     string(stringQuestion),
			MessageType: QUESTION,
		})
		client.send <- messageToSend
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
			for otherClientConn := range clients.Clients {
				if otherClientConn.conn != client.conn {
					message := "The other user has disconnected."

					messageToSend := Message{
						Message:     message,
						MessageType: INFO,
						RoomID:      roomID,
					}

					jsonMessage, _ := json.Marshal(messageToSend)

					otherClientConn.send <- jsonMessage
					otherClientConn.conn.Close()
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
