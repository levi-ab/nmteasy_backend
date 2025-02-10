package sockets

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"nmteasy_backend/models"
	"nmteasy_backend/models/migrated_models"
	"nmteasy_backend/utils"
	"strings"
	"time"
)

const OPPONENT_DISCONECTED string = "opponent_disconnected"
const ERROR string = "error"
const ANSWER string = "answer"
const SKIP_QUESTION string = "skip_question"
const QUESTION string = "question"
const RESULT string = "result"
const FINISHED string = "finished"
const MATCH_FOUND string = "match_found"

type Room struct {
	Clients   map[*Client]bool `json:"-"`
	GameState GameState        `json:"gameState"`
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

	//Match Making queue
	questionTypeQueues map[string][]*Client

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
		broadcast:          make(chan *Message),
		register:           make(chan *Client),
		unregister:         make(chan *Client),
		clients:            make(map[*Client]bool),
		questionTypeQueues: make(map[string][]*Client, 0),
		rooms:              make(map[string]Room),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				h.removeFromRooms(client)
				h.removeFromMatchmakingQueue(client)
				delete(h.clients, client)
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
						h.finishTheGame(room, sender, anotherClient, false)
						continue
					}

					room.GameState.CurrentIndex = room.GameState.CurrentIndex + 1
					room.GameState.ClientRightCounts[sender.clientID] = room.GameState.ClientRightCounts[sender.clientID] + 2

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

				if room.GameState.CurrentIndex+1 == len(room.GameState.Questions) {
					//finishing the game
					var firstClient *Client
					var secondClient *Client

					for client := range room.Clients {
						if firstClient == nil {
							firstClient = client
						} else {
							secondClient = client
						}
					}
					h.finishTheGame(room, firstClient, secondClient, true)
					continue
				}

				room.GameState.CurrentIndex = room.GameState.CurrentIndex + 1
				h.rooms[message.RoomID] = room

				sendQuestion(room.GameState.Questions[room.GameState.CurrentIndex], room.Clients)
			}

		case <-time.After(time.Second):
			for questionType, queue := range h.questionTypeQueues { // Adjust the interval as needed
				if len(queue) >= 2 {
					// Take the first two clients from the queue
					client1 := queue[0]
					client2 := queue[1]

					// Remove them from the queue
					h.questionTypeQueues[questionType] = queue[2:]

					room := client1.clientID.String() + client2.clientID.String() + questionType

					// Notify clients that they have found a match
					msg := Message{
						RoomID:      room,
						Message:     client2.clientName,
						MessageType: MATCH_FOUND,
					}

					messageToSend, err := json.Marshal(msg)
					if err != nil {
						fmt.Println("Error marshaling message:", err)
						continue
					}
					client1.send <- messageToSend

					msg.Message = client1.clientName

					messageToSend, err = json.Marshal(msg)
					if err != nil {
						fmt.Println("Error marshaling message:", err)
						continue
					}

					client2.send <- messageToSend

					client2.queueName = ""
					client1.queueName = ""

					connections := make(map[*Client]bool)
					connections[client1] = true
					connections[client2] = true

					//then here we query the questions and send the first question
					questions, err := utils.GetRandomQuestionsByType(questionType, 10)

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
}

func sendQuestion(question models.Question, clients map[*Client]bool) {
	stringQuestion, _ := json.Marshal(question)

	for client := range clients {
		messageToSend, _ := json.Marshal(Message{
			Message:     string(stringQuestion),
			MessageType: QUESTION,
		})
		client.send <- messageToSend
	}
}

func (h *Hub) finishTheGame(room Room, sender *Client, anotherClient *Client, skippedLastQuestion bool) {
	var userResult int
	var UserAnsweredLast bool

	if skippedLastQuestion {
		userResult = room.GameState.ClientRightCounts[sender.clientID]
	} else {
		userResult = room.GameState.ClientRightCounts[sender.clientID] + 2 //+1 cause he answered this question right as well
		UserAnsweredLast = true
	}

	result := struct {
		UserResult     int
		OpponentResult int
		AnsweredLast   bool
	}{
		UserResult:     userResult,
		OpponentResult: room.GameState.ClientRightCounts[anotherClient.clientID],
		AnsweredLast:   UserAnsweredLast,
	}

	resultMessage, err := json.Marshal(result)

	msg := Message{
		Message:     string(resultMessage),
		MessageType: FINISHED,
	}

	messageToSend, err := json.Marshal(msg)
	if err != nil {
		fmt.Println("Error marshaling message:", err)
		return
	}

	sender.send <- messageToSend

	result.UserResult, result.OpponentResult = result.OpponentResult, result.UserResult
	result.AnsweredLast = false

	resultMessage, err = json.Marshal(result)

	msg = Message{
		Message:     string(resultMessage),
		MessageType: FINISHED,
	}

	messageToSend, err = json.Marshal(msg)
	if err != nil {
		fmt.Println("Error marshaling message:", err)
		return
	}

	anotherClient.send <- messageToSend

	var firstUser migrated_models.User
	models.DB.Where("id = ?", sender.clientID).First(&firstUser)
	firstUser.Points += result.OpponentResult
	models.DB.Save(&firstUser)

	var secondUser migrated_models.User
	models.DB.Where("id = ?", anotherClient.clientID).First(&secondUser)
	secondUser.Points += result.UserResult
	models.DB.Save(&secondUser)

	h.unregister <- sender
	h.unregister <- anotherClient
}

func (h *Hub) removeFromRooms(clientToRemove *Client) {
	// Create a new map for updated rooms
	updatedRooms := make(map[string]Room)

	// Iterate over each room
	for roomID, room := range h.rooms {
		// Check if the clientToRemove is in the room
		if !strings.Contains(roomID, clientToRemove.clientID.String()) {
			// If the room ID doesn't contain the clientToRemove's clientID, include it in the updated map
			updatedRooms[roomID] = room
		} else {
			for anotherClientConn := range room.Clients {
				if anotherClientConn.conn != clientToRemove.conn {

					result := struct {
						UserResult     int
						OpponentResult int
					}{
						UserResult:     room.GameState.ClientRightCounts[anotherClientConn.clientID],
						OpponentResult: room.GameState.ClientRightCounts[clientToRemove.clientID],
					}

					resultMessage, _ := json.Marshal(result)

					messageToSend := Message{
						Message:     string(resultMessage),
						MessageType: OPPONENT_DISCONECTED,
						RoomID:      roomID,
					}

					jsonMessage, _ := json.Marshal(messageToSend)

					anotherClientConn.send <- jsonMessage

					var firstUser migrated_models.User
					models.DB.Where("id = ?", clientToRemove.clientID).First(&firstUser)
					firstUser.Points += result.OpponentResult
					models.DB.Save(&firstUser)

					var secondUser migrated_models.User
					models.DB.Where("id = ?", anotherClientConn.clientID).First(&secondUser)
					secondUser.Points += result.UserResult
					models.DB.Save(&secondUser)

					delete(h.clients, anotherClientConn)
					break // Assuming there's only one other clientToRemove in the room
				}
			}
		}
	}

	// Replace the original rooms with the updated map
	h.rooms = updatedRooms
	clientToRemove.conn.Close()
}

func (h *Hub) removeFromMatchmakingQueue(client *Client) {
	// Create a new queue without the client
	var updatedQueue []*Client
	for _, c := range h.questionTypeQueues[client.queueName] {
		if c.clientID != client.clientID {
			updatedQueue = append(updatedQueue, c)
		}
	}
	h.questionTypeQueues[client.queueName] = updatedQueue
}
