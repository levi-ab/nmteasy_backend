package sockets

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"log"
	"net/http"
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
}

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	hub *Hub

	// The websocket connection.
	clientID       uuid.UUID //client id
	IsInQueue      bool      //weather the client is in queue
	conn           *websocket.Conn
	CorrectAnswers int
	// Buffered channel of outbound messages.
	send       chan []byte
	writeMutex sync.Mutex
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
			c.IsInQueue = true

			// Add the client to the matchmaking queue
			c.hub.matchmakingQueue = append(c.hub.matchmakingQueue, c)
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

			// Add queued chat messages to the current websocket message.
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-c.send)
			}

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

	user := utils.GetCurrentUser(r)

	client := &Client{hub: hub, conn: conn, send: make(chan []byte, 256), clientID: user.ID}
	client.hub.register <- client

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.writePump()
	go client.readPump()
}
