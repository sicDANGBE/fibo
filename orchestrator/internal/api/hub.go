package api

import (
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

type Hub struct {
	Clients   map[*websocket.Conn]bool
	Broadcast chan interface{}
	Mu        sync.Mutex
}

func NewHub() *Hub {
	return &Hub{
		Clients:   make(map[*websocket.Conn]bool),
		Broadcast: make(chan interface{}),
	}
}

func (h *Hub) Run() {
	for msg := range h.Broadcast {
		h.Mu.Lock()
		for client := range h.Clients {
			if err := client.WriteJSON(msg); err != nil {
				client.Close()
				delete(h.Clients, client)
			}
		}
		h.Mu.Unlock()
	}
}

func ServeWs(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	hub.Mu.Lock()
	hub.Clients[conn] = true
	hub.Mu.Unlock()
}

func (h *Hub) BroadcastMessage(msg interface{}) {
	h.Broadcast <- msg
}
