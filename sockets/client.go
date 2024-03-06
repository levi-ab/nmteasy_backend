package sockets

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"nmteasy_backend/common"
	"nmteasy_backend/models/migrated_models"
	"nmteasy_backend/utils"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 60 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 520 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	hub            *Hub
	clientID       uuid.UUID //client id
	clientName     string
	queueName      string          //whether the client is in queue
	conn           *websocket.Conn // the websocket connection.
	CorrectAnswers int
	send           chan []byte // channel of outbound messages.
	writeMutex     sync.Mutex
}

// readPump pumps messages from the websocket connection to the hub.
//
// The application runs readPump in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}

		messageData := bytes.TrimSpace(bytes.Replace(message, newline, space, -1))

		var targetMessage Message

		if err = json.Unmarshal(messageData, &targetMessage); err != nil {
			errorMessage := Message{
				Message:     "failed to parse the message",
				MessageType: ERROR,
			}

			messageToSend, err := json.Marshal(errorMessage)
			if err != nil {
				fmt.Println("Error marshaling message:", err)
				continue
			}

			c.send <- messageToSend
			continue
		}

		// Send the message to the hub for processing
		c.hub.broadcast <- &targetMessage

		if targetMessage.MessageType == "join_matchmaking" {
			// Set the client state to indicate they are in the matchmaking queue
			c.queueName = targetMessage.Message

			// Add the client to the matchmaking queue
			c.hub.questionTypeQueues[targetMessage.Message] = append(c.hub.questionTypeQueues[targetMessage.Message], c)
			continue
		}
	}
}

// writePump pumps messages from the hub to the websocket connection.
//
// A goroutine running writePump is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			c.writeMutex.Lock()
			if !ok {
				// The hub closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				c.writeMutex.Unlock()
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				c.writeMutex.Unlock()
				return
			}
			w.Write(message)

			if err := w.Close(); err != nil {
				c.writeMutex.Unlock()
				return
			}
			c.writeMutex.Unlock()

		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			c.writeMutex.Lock()
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				c.writeMutex.Unlock()
				return
			}
			c.writeMutex.Unlock()
		}
	}
}

// ServeWs handles websocket requests from the peer.
func ServeWs(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	var user *migrated_models.User

	token := mux.Vars(r)["token"] //todo huge security risk but i dont give a damn cuz stupid js wont change one thing in 14 years and neither will i
	mySigningKey := []byte(common.SECRET_KEY)
	parsedToken, err := utils.ParseToken(token, mySigningKey)
	if err != nil || parsedToken == nil {
		log.Println(err)
		return
	}

	if claims, ok := parsedToken.Claims.(jwt.MapClaims); ok {
		if email, ok := claims["email"].(string); ok && email != "" {
			_user, err := utils.GetUserByEmail(email)
			if err == nil {
				user = _user
			}
		}
	}

	client := &Client{hub: hub, conn: conn, send: make(chan []byte, 256), clientID: user.ID, clientName: user.FirstName + " " + user.LastName}
	client.hub.register <- client

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.writePump()
	go client.readPump()
}
