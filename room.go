package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type room struct {
	// forward is channel to send message for other clients
	forward chan []byte
	// join is channel for clinet to try joinning to ChatRoom
	join chan *client
	// leave is channel for client to try leaving ChatRoom
	leave chan *client
	// clients takes all clients joinning to ChatRoom
	clients map[*client]bool
}

// newRoom create ChatRoom to be able to use right now and return this ChatRoom
func newRoom() *room {
	return &room{
		forward: make(chan []byte),
		join:    make(chan *client),
		leave:   make(chan *client),
		clients: make(map[*client]bool),
	}
}

func (r *room) run() {
	for {
		select {
		case client := <-r.join:
			// Join
			r.clients[client] = true
		case client := <-r.leave:
			// Leave
			delete(r.clients, client)
			close(client.send)
		case msg := <-r.forward:
			// send message to all clients
			for client := range r.clients {
				select {
				case client.send <- msg:
					// send message
				default:
					// failed to send message
					delete(r.clients, client)
					close(client.send)
				}
			}
		}
	}
}

const (
	socketBufferSize  = 1024
	messageBufferSize = 256
)

var upgrader = &websocket.Upgrader{ReadBufferSize: socketBufferSize, WriteBufferSize: socketBufferSize}

func (r *room) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	socket, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		log.Fatal("ServeHTTP:", err)
		return
	}
	client := &client{
		socket: socket,
		send:   make(chan []byte, messageBufferSize),
		room:   r,
	}
	r.join <- client
	defer func() { r.leave <- client }()
	go client.write()
	client.read()
}
