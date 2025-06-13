package sockets

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"math/rand"
	"nmteasy_backend/models"
	"nmteasy_backend/utils"
	"strings"
	"sync"
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

const DEFAULT_BOT_TIME_WAIT = 20 * time.Second

type Room struct {
	Clients   map[*Client]bool `json:"-"`
	GameState GameState        `json:"gameState"`
	mutex     sync.RWMutex     `json:"-"`
}

type GameState struct {
	Questions         []models.Question `json:"questions"`
	CurrentIndex      int               `json:"currentIndex"`
	ClientRightCounts map[uuid.UUID]int `json:"clientRightCounts"`
}

type ConcurrentHub struct {
	// Thread-safe client storage
	clients sync.Map // map[*Client]bool

	// Thread-safe room storage
	rooms sync.Map // map[string]*Room

	// Thread-safe queue storage
	questionTypeQueues sync.Map // map[string][]*Client
	queueMutex         sync.RWMutex

	// Channels for different operations
	broadcast  chan *Message
	register   chan *Client
	unregister chan *Client

	// Worker configuration
	messageWorkers    int
	matchmakingTicker *time.Ticker

	// Shutdown coordination
	shutdown chan struct{}
	wg       sync.WaitGroup
}

type Message struct {
	Message     string
	RoomID      string
	MessageType string
	ClientID    uuid.UUID
}

type AnswerMessage struct {
	Answer string
	UserID uuid.UUID
}

func NewHub() *ConcurrentHub {
	return &ConcurrentHub{
		broadcast:         make(chan *Message, 1000),
		register:          make(chan *Client, 100),
		unregister:        make(chan *Client, 100),
		messageWorkers:    10,
		matchmakingTicker: time.NewTicker(time.Second),
		shutdown:          make(chan struct{}),
	}
}

func (h *ConcurrentHub) Run() {
	// Start registration handle
	go h.handleRegistrations()

	// Start message processing workers
	for i := 0; i < h.messageWorkers; i++ {
		h.wg.Add(1)
		go h.messageWorker(i)
	}

	go h.handleMatchmaking()

	// Wait for shutdown
	<-h.shutdown

	// Close channels and wait for goroutines to finish
	close(h.broadcast)
	close(h.register)
	close(h.unregister)
	h.matchmakingTicker.Stop()

	h.wg.Wait()
}

func (h *ConcurrentHub) Shutdown() {
	close(h.shutdown)
}

func (h *ConcurrentHub) handleRegistrations() {
	defer h.wg.Done()

	for {
		select {
		case client := <-h.register:
			if client != nil {
				client.joinedAt = time.Now()
				h.clients.Store(client, true)
			}

		case client := <-h.unregister:
			if client != nil {
				if _, ok := h.clients.Load(client); ok {
					h.removeFromRooms(client)
					h.removeFromMatchmakingQueue(client)
					h.clients.Delete(client)
					//close(client.send)
				}
			}

		case <-h.shutdown:
			return
		}
	}
}

func (h *ConcurrentHub) messageWorker(workerID int) {
	defer h.wg.Done()

	for {
		select {
		case message := <-h.broadcast:
			if message != nil {
				h.processMessage(message)
			}
		case <-h.shutdown:
			return
		}
	}
}

func (h *ConcurrentHub) processMessage(message *Message) {
	switch message.MessageType {
	case ANSWER:
		h.handleAnswer(message)
	case SKIP_QUESTION:
		h.handleSkipQuestion(message)
	}
}

func (h *ConcurrentHub) handleAnswer(message *Message) {
	roomInterface, ok := h.rooms.Load(message.RoomID)
	if !ok {
		fmt.Println("room not found")
		return
	}

	room := roomInterface.(*Room)
	room.mutex.Lock()
	defer room.mutex.Unlock()

	var parsedAnswerMessage AnswerMessage
	if err := json.Unmarshal([]byte(message.Message), &parsedAnswerMessage); err != nil {
		fmt.Println("failed to parse the answer")
		return
	}

	var sender *Client
	var anotherClient *Client

	for client := range room.Clients {
		if client.clientID == parsedAnswerMessage.UserID {
			sender = client
		} else {
			anotherClient = client
		}
	}

	if sender == nil {
		fmt.Println("failed to determine the sender")
		return
	}

	if room.GameState.Questions[room.GameState.CurrentIndex].RightAnswer == parsedAnswerMessage.Answer {
		if room.GameState.CurrentIndex+1 == len(room.GameState.Questions) {
			// Finishing the game
			h.finishTheGame(room, sender, anotherClient, false)
			return
		}

		room.GameState.CurrentIndex = room.GameState.CurrentIndex + 1
		room.GameState.ClientRightCounts[sender.clientID] = room.GameState.ClientRightCounts[sender.clientID] + 2

		msg := Message{
			Message:     "correct",
			MessageType: RESULT,
		}

		messageToSend, err := json.Marshal(msg)
		if err != nil {
			fmt.Println("Error marshaling message:", err)
			return
		}

		sender.send <- messageToSend
		msg.Message = "other_answered"

		messageToSend, err = json.Marshal(msg)
		if err != nil {
			fmt.Println("Error marshaling message:", err)
			return
		}

		anotherClient.send <- messageToSend

		h.sendQuestion(room.GameState.Questions[room.GameState.CurrentIndex], room.Clients)

	} else {
		msg := Message{
			Message:     "wrong",
			MessageType: RESULT,
		}

		messageToSend, err := json.Marshal(msg)
		if err != nil {
			fmt.Println("Error marshaling message:", err)
			return
		}
		sender.send <- messageToSend

		if sender.clientName != DEFAULT_BOT_NANE {
			hasBotClient := false
			for client := range room.Clients {
				if client.clientName == DEFAULT_BOT_NANE {
					hasBotClient = true
					break
				}
			}

			// If there's a bot, simulate its answer with a delay
			if hasBotClient {
				go func(r *Room, roomID string) {
					// Random delay before bot answers (between 1-4 seconds)
					time.Sleep(time.Duration(rand.Intn(3000)+1000) * time.Millisecond)
					h.simulateBotAnswer(r, roomID)
				}(room, message.RoomID)
			}
		}
	}
}

func (h *ConcurrentHub) handleSkipQuestion(message *Message) {
	roomInterface, ok := h.rooms.Load(message.RoomID)
	if !ok {
		fmt.Println("room not found")
		return
	}

	room := roomInterface.(*Room)
	room.mutex.Lock()
	defer room.mutex.Unlock()

	if room.GameState.CurrentIndex+1 == len(room.GameState.Questions) {
		// Finishing the game
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
		return
	}

	room.GameState.CurrentIndex = room.GameState.CurrentIndex + 1

	h.sendQuestion(room.GameState.Questions[room.GameState.CurrentIndex], room.Clients)
}

func (h *ConcurrentHub) handleMatchmaking() {
	defer h.wg.Done()

	for {
		select {
		case <-h.matchmakingTicker.C:
			h.processMatchmaking()
		case <-h.shutdown:
			return
		}
	}
}

func (h *ConcurrentHub) processMatchmaking() {
	h.queueMutex.Lock()
	defer h.queueMutex.Unlock()

	h.questionTypeQueues.Range(func(key, value interface{}) bool {
		questionType := key.(string)
		queue := value.([]*Client)
		println(len(queue))
		if len(queue) == 1 {
			println(queue[0].clientName)
			println(queue[0].clientName)
		}

		if len(queue) == 1 {
			client := queue[0]
			if time.Since(client.joinedAt) >= DEFAULT_BOT_TIME_WAIT {
				// Create a bot client
				botClient := h.createBotClient(questionType)

				// Add bot to the queue
				queue = append(queue, botClient)
				h.questionTypeQueues.Store(questionType, queue)
				println("client with name " + client.clientName + " is playing bot")
			}
		} else if len(queue) >= 2 {
			// Take the first two clients from the queue
			client1 := queue[0]
			client2 := queue[1]

			// Update queue
			h.questionTypeQueues.Store(questionType, queue[2:])

			// Create match
			h.createMatch(client1, client2, questionType)
			println("client with name " + client1.clientName + " is playing real player player" + client2.clientName)
		}
		return true
	})
}

func (h *ConcurrentHub) createMatch(client1, client2 *Client, questionType string) {
	room := client1.clientID.String() + client2.clientID.String() + questionType

	// Notify clients that they have found a match
	msg := Message{
		RoomID:      room,
		Message:     client2.clientName,
		ClientID:    client1.clientID,
		MessageType: MATCH_FOUND,
	}

	// Send match found messages
	if client1.conn != nil {
		messageToSend, _ := json.Marshal(msg)
		client1.send <- messageToSend
	}

	msg.Message = client1.clientName
	msg.ClientID = client2.clientID

	if client2.conn != nil {
		messageToSend, _ := json.Marshal(msg)
		client2.send <- messageToSend
	}

	client2.queueName = ""
	client1.queueName = ""

	connections := make(map[*Client]bool)
	connections[client1] = true
	connections[client2] = true

	println("created client1 and client2" + client1.clientName + client2.clientName)

	// Query questions
	questions, err := utils.GetRandomQuestionsByType(questionType, 10)
	if err != nil {
		fmt.Println("Error getting questions:", err)
		return
	}

	newRoom := &Room{
		Clients: connections,
		GameState: GameState{
			Questions:         questions,
			ClientRightCounts: make(map[uuid.UUID]int),
		},
	}

	h.rooms.Store(room, newRoom)

	h.sendQuestion(questions[0], connections)
}

func (h *ConcurrentHub) sendQuestion(question models.Question, clients map[*Client]bool) {
	stringQuestion, _ := json.Marshal(question)

	for client := range clients {
		messageToSend, _ := json.Marshal(Message{
			Message:     string(stringQuestion),
			MessageType: QUESTION,
		})
		client.send <- messageToSend
	}
}

func (h *ConcurrentHub) finishTheGame(room *Room, sender *Client, anotherClient *Client, skippedLastQuestion bool) {
	var userResult int
	var UserAnsweredLast bool

	if skippedLastQuestion {
		userResult = room.GameState.ClientRightCounts[sender.clientID]
	} else {
		userResult = room.GameState.ClientRightCounts[sender.clientID] + 2
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

	UpdateUserPoints(sender.clientID, result.OpponentResult)
	UpdateUserPoints(anotherClient.clientID, result.UserResult)

	h.unregister <- sender
	h.unregister <- anotherClient
}

func (h *ConcurrentHub) removeFromRooms(clientToRemove *Client) {
	var roomsToDelete []string

	h.rooms.Range(func(key, value interface{}) bool {
		roomID := key.(string)
		room := value.(*Room)

		// Check if the clientToRemove is in the room
		if strings.Contains(roomID, clientToRemove.clientID.String()) {
			room.mutex.Lock()

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

					UpdateUserPoints(clientToRemove.clientID, result.OpponentResult)
					UpdateUserPoints(anotherClientConn.clientID, result.UserResult)

					h.clients.Delete(anotherClientConn)
					break
				}
			}

			room.mutex.Unlock()
			roomsToDelete = append(roomsToDelete, roomID)
		}
		return true
	})

	// Delete rooms outside of the Range loop
	for _, roomID := range roomsToDelete {
		h.rooms.Delete(roomID)
	}

	if clientToRemove.conn != nil {
		clientToRemove.conn.Close()
	}
}

func (h *ConcurrentHub) removeFromMatchmakingQueue(client *Client) {
	if client.queueName == "" {
		return
	}

	h.queueMutex.Lock()
	defer h.queueMutex.Unlock()

	queueInterface, ok := h.questionTypeQueues.Load(client.queueName)
	if !ok {
		return
	}

	queue := queueInterface.([]*Client)
	var updatedQueue []*Client

	for _, c := range queue {
		if c.clientID != client.clientID {
			updatedQueue = append(updatedQueue, c)
		}
	}

	h.questionTypeQueues.Store(client.queueName, updatedQueue)
}

func (h *ConcurrentHub) addToQueue(questionType string, client *Client) {
	h.queueMutex.Lock()
	defer h.queueMutex.Unlock()

	queueInterface, _ := h.questionTypeQueues.LoadOrStore(questionType, []*Client{})
	queue := queueInterface.([]*Client)
	queue = append(queue, client)
	h.questionTypeQueues.Store(questionType, queue)
}
