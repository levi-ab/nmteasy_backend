package sockets

import (
	"encoding/json"
	"github.com/google/uuid"
	"math/rand"
	"time"
)

const DEFAULT_BOT_NANE = "AI Bot"

// createBotClient creates a bot client for matchmaking
func (h *ConcurrentHub) createBotClient(questionType string) *Client {
	botClient := &Client{
		hub:        h,
		clientID:   uuid.New(),
		clientName: DEFAULT_BOT_NANE,
		queueName:  questionType,
		conn:       nil, // Bot doesn't have a real connection
		send:       make(chan []byte, 256),
		joinedAt:   time.Now(),
	}

	// Register the bot client
	h.clients.Store(botClient, true)

	return botClient
}

// simulateBotAnswer simulates a bot answering a question
func (h *ConcurrentHub) simulateBotAnswer(room *Room, roomID string) {
	room.mutex.RLock()
	defer room.mutex.RUnlock()

	// Find the bot client
	var botClient *Client
	for client := range room.Clients {
		if client.clientName == DEFAULT_BOT_NANE {
			botClient = client
			break
		}
	}

	if botClient == nil {
		return
	}

	// Get current question
	if room.GameState.CurrentIndex >= len(room.GameState.Questions) {
		return
	}

	currentQuestion := room.GameState.Questions[room.GameState.CurrentIndex]

	// Bot answers correctly 30% of the time
	var answer string
	if rand.Float32() < 0.3 {
		answer = currentQuestion.RightAnswer
	} else {
		// Pick a random wrong answer
		answer = "wrong"
	}

	// Create answer message
	answerMsg := AnswerMessage{
		Answer: answer,
		UserID: botClient.clientID,
	}

	answerJSON, _ := json.Marshal(answerMsg)

	message := &Message{
		Message:     string(answerJSON),
		RoomID:      roomID,
		MessageType: ANSWER,
		ClientID:    botClient.clientID,
	}

	// Send the bot's answer to the broadcast channel
	select {
	case h.broadcast <- message:
	default:
		// Channel is full, skip this bot answer
	}
}
