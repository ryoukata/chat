package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/ryoukata/trace"
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
	// tracer receives Logs operated on ChatRoom.
	tracer trace.Tracer
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
			r.tracer.Trace("New Client joined.")
		case client := <-r.leave:
			// Leave
			delete(r.clients, client)
			close(client.send)
			r.tracer.Trace("Client leaved.")
		case msg := <-r.forward:
			r.tracer.Trace("received message: ", string(msg))
			// send message to all clients
			for client := range r.clients {
				select {
				case client.send <- msg:
					// send message
					r.tracer.Trace(" -- send to Client.")
				default:
					// failed to send message
					delete(r.clients, client)
					close(client.send)
					r.tracer.Trace(" -- failed to send. clean up Client.")
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
