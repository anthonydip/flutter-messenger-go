package ws

import (
	"strings"
)

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// Registered clients.
	clients map[*Client]bool

	// Inbound messages from the clients.
	broadcast chan []byte

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client

	// Mapping of user id to clients
	userIds map[string]*Client
}

func NewHub() *Hub {
	return &Hub{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
		userIds:    make(map[string]*Client),
	}
}

func (h *Hub) Run() {
	for {
		select {
		// When a new client registers with the hub
		case client := <-h.register:
			h.clients[client] = true
			h.userIds[client.userId] = client
		// When a client requests to unregister (disconnects)
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
		// When a message is broadcasted to all connected clients
		case message := <-h.broadcast:
			// Check if the message is a private message
			// Private message is of the form "/msg <sender id> <recipient id> <message>"
			if strings.HasPrefix(string(message), "/msg ") {
				// Extract the user id from the message
				parts := strings.SplitN(string(message[5:]), " ", 3)
				targetId := parts[1]

				// Find the target client and send the private messsage
				if targetId, ok := h.userIds[targetId]; ok {
					select {
					case targetId.send <- []byte(message[5:]):
					default:
						close(targetId.send)
						delete(h.clients, targetId)
						delete(h.userIds, targetId.userId)
					}
				}
			}
		}
	}
}
