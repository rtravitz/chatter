package handlers

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/rtravitz/chatter/room"
	"log"
	"net/http"
)

const (
	MESSAGE = "MESSAGE"
	USERID  = "USERID"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// Room holds extra resources that its routes need to access
type Room struct {
}

// OutputMessage ...
type OutputMessage struct {
	Data   string `json:"data"`
	Sender string `json:"sender"`
	Type   string `json:"type"`
}

// Listen listens and sends messages over websockets
func (r *Room) Listen(sub, unsub chan room.User, pub chan room.Message) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println(err)
			return
		}

		user := room.NewUser()
		output := OutputMessage{
			Data:   string(user.ID.String()),
			Sender: user.ID.String(),
			Type:   USERID,
		}
		val, err := json.Marshal(output)
		if err != nil {
			log.Printf("Failed creating message for user %v\n", user.ID)
			return
		}
		//send user back their id

		conn.WriteMessage(websocket.TextMessage, val)

		log.Printf("Made new user: %+v\n", user)
		sub <- user

		// write messages back to the client in a separate goroutine that
		// closes when the user's messages channel is closed
		go func() {
			for mess := range user.Messages {
				log.Printf("Waiting for messages for user %v\n", user.ID)
				output := OutputMessage{
					Data:   string(mess.Data),
					Sender: mess.Sender.ID.String(),
					Type:   MESSAGE,
				}
				val, err := json.Marshal(output)
				if err != nil {
					log.Printf("Failed creating message for user %v\n", user.ID)
					continue
				}
				conn.WriteMessage(mess.MessageType, val)
			}
		}()

		// read messages from the websocket until the socket is
		// closed by the client
		for {
			messageType, p, err := conn.ReadMessage()
			if err != nil {
				log.Println(err)
				unsub <- user
				close(user.Messages)
				return
			}

			log.Printf("Sending message for user %v\n", user.ID)
			// broadcast messages to all other users
			pub <- room.Message{MessageType: messageType, Data: p, Sender: user}
		}
	}
}
