package main

import (
	"time"

	"github.com/gorilla/websocket"
)

// client is one chatting user
type client struct {
	// socket is websocket for this client
	socket *websocket.Conn
	// send is channel that message is sent
	send chan *message
	// room is ChatRoom that this clinet is joinning
	room *room
	// userData takes information about User.
	userData map[string]interface{}
}

func (c *client) read() {
	for {
		var msg *message
		if err := c.socket.ReadJSON(&msg); err == nil {
			msg.When = time.Now()
			msg.Name = c.userData["name"].(string)
			c.room.forward <- msg
		} else {
			break
		}
	}
	c.socket.Close()
}

func (c *client) write() {
	for msg := range c.send {
		if err := c.socket.WriteJSON(msg); err != nil {
			break
		}
	}
	c.socket.Close()
}
