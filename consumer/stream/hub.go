package stream

import "consumer/models"

type Hub struct {
	clients    map[chan models.Post]struct{}
	register   chan chan models.Post
	unregister chan chan models.Post
	broadcast  chan models.Post
}

func NewHub(buffer int) *Hub {
	return &Hub{
		clients:    make(map[chan models.Post]struct{}),
		register:   make(chan chan models.Post),
		unregister: make(chan chan models.Post),
		broadcast:  make(chan models.Post, buffer),
	}
}

func (hub *Hub) Run() {
	for {
		select {
		case ch := <-hub.register:
			hub.clients[ch] = struct{}{}
		case ch := <-hub.unregister:
			if _, ok := hub.clients[ch]; ok {
				delete(hub.clients, ch)
				close(ch)
			}
		case post := <-hub.broadcast:
			for ch := range hub.clients {
				select {
				case ch <- post:
				default:
				}
			}
		}
	}
}

func (hub *Hub) Register(buffer int) chan models.Post {
	ch := make(chan models.Post, buffer)
	hub.register <- ch
	return ch
}

func (hub *Hub) Unregister(ch chan models.Post) {
	hub.unregister <- ch
}

func (hub *Hub) Broadcast(post models.Post) {
	hub.broadcast <- post
}

func (h *Hub) Close() {
	for ch := range h.clients {
		close(ch)
	}
	h.clients = make(map[chan models.Post]struct{})
}
