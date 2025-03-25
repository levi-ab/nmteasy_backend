package sockets

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"math/rand"
)

const DEFAULT_BOT_NANE = "AI Bot"

func (h *Hub) createBotClient(questionType string) *Client {
	botID := uuid.New()
	botClient := &Client{
		hub:      h,
		clientID: botID,
		//clientName:     "AI " + questionType + " Bot ",
		clientName:     DEFAULT_BOT_NANE,
		queueName:      questionType,
		conn:           nil,
		send:           make(chan []byte, 256),
		CorrectAnswers: 0,
	}

	return botClient
}

func (h *Hub) simulateBotAnswer(room Room, roomID string) {
	// Find the bot client in the room
	var botClient *Client
	var humanClient *Client
	for client := range room.Clients {
		if client.clientName == DEFAULT_BOT_NANE {
			botClient = client
		} else {
			humanClient = client
		}
	}

	if botClient == nil || humanClient == nil {
		return
	}

	// Get the current question
	currentQuestion := room.GameState.Questions[room.GameState.CurrentIndex]

	// Simulate bot answer with some randomness
	var botAnswer string

	if rand.Float32() < 0.1 {
		botAnswer = currentQuestion.RightAnswer
	} else {
		botAnswer = "incorrectAnswers[rand.Intn(len(incorrectAnswers))]"
	}

	// Prepare the bot answer message
	botAnswerMessage := AnswerMessage{
		Answer: botAnswer,
		UserID: botClient.clientID,
	}

	// Marshal the answer message
	botAnswerBytes, err := json.Marshal(botAnswerMessage)
	if err != nil {
		fmt.Println("Error marshaling bot answer:", err)
		return
	}

	// Broadcast the bot's answer
	h.broadcast <- &Message{
		Message:     string(botAnswerBytes),
		MessageType: ANSWER,
		RoomID:      roomID,
	}
}
