package ws

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

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

type Response struct {
	Status        string `json:"status"`
	StatusCode    int    `json:"statusCode"`
	StatusMessage string `json:"statusMessage,omitempty"`
}

// Client is a middleman between the websocket connection and the hub
type Client struct {
	hub *Hub

	// Websocket connection
	conn *websocket.Conn

	// Buffered channel of outbound messages
	send chan []byte

	// Hold user ID for the client
	userId string
}

// readPump pumps messages from the websocket connection to the hub.
//
// The application runs readPump in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
func (c *Client) ReadPump() {
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
				log.Info().Msgf("[/ws] WebSocket connection closed for %s", c.userId)
			}
			break
		}
		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
		c.hub.broadcast <- message
	}
}

// writePump pumps messages from the hub to the websocket connection.
//
// A goroutine running writePump is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func (c *Client) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
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
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func ServeWs(hub *Hub, w http.ResponseWriter, r *http.Request, id string) {
	if id == "" {
		log.Error().Msgf("[GET /ws] Invalid user id %s", id)
		res := Response{
			Status:        "INTERNAL SERVER ERROR",
			StatusCode:    500,
			StatusMessage: "Unable to parse user from token",
		}

		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(&res)
		return
	}

	// TO-DO: Maybe remove later?
	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Error().Msgf("[GET /ws] Error upgrading to WebSocket connection for %s", id)
		res := Response{
			Status:        "INTERNAL SERVER ERROR",
			StatusCode:    500,
			StatusMessage: "Unable to parse user from token",
		}

		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(&res)
		return
	}

	client := &Client{hub: hub, conn: conn, send: make(chan []byte, 256), userId: id}
	client.hub.register <- client

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines
	go client.WritePump()
	go client.ReadPump()

	log.Info().Msgf("[GET /ws] Established WebSocket connection for %s", id)
}
