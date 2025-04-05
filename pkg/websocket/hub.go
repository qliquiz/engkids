package websocket

import (
	"github.com/gofiber/websocket"
)

type Hub struct {
	Clients   map[*websocket.Conn]bool
	Broadcast chan []byte
}

func NewHub() *Hub {
	return &Hub{
		Clients:   make(map[*websocket.Conn]bool),
		Broadcast: make(chan []byte),
	}
}

func (h *Hub) Run() {
	for {
		msg := <-h.Broadcast
		for client := range h.Clients {
			err := client.WriteMessage(websocket.TextMessage, msg)
			if err != nil {
				return
			}
		}
	}
}
