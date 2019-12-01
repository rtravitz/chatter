package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
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

// OutputMessage ...
type OutputMessage struct {
	Data   string `json:"data"`
	Sender string `json:"sender"`
	Type   string `json:"type"`
}

func roomHandler(sub, unsub chan user, pub chan message) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println(err)
			return
		}

		user := newUser()
		output := OutputMessage{
			Data:   string(user.id.String()),
			Sender: user.id.String(),
			Type:   USERID,
		}
		val, err := json.Marshal(output)
		if err != nil {
			log.Printf("Failed creating message for user %v\n", user.id)
			return
		}
		//send user back their id

		conn.WriteMessage(websocket.TextMessage, val)

		log.Printf("Made new user: %+v\n", user)
		sub <- user

		// write messages back to the client in a separate goroutine that
		// closes when the user's messages channel is closed
		go func() {
			for mess := range user.messages {
				log.Printf("Waiting for messages for user %v\n", user.id)
				output := OutputMessage{
					Data:   string(mess.data),
					Sender: mess.sender.id.String(),
					Type:   MESSAGE,
				}
				val, err := json.Marshal(output)
				if err != nil {
					log.Printf("Failed creating message for user %v\n", user.id)
					continue
				}
				conn.WriteMessage(mess.messageType, val)
			}
		}()

		// read messages from the websocket until the socket is
		// closed by the client
		for {
			messageType, p, err := conn.ReadMessage()
			if err != nil {
				log.Println(err)
				unsub <- user
				close(user.messages)
				return
			}

			log.Printf("Sending message for user %v\n", user.id)
			// broadcast messages to all other users
			pub <- message{messageType: messageType, data: p, sender: user}
		}
	}
}

func getPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		return "localhost:5000"
	}

	return ":" + port
}

func main() {
	sub := make(chan user)
	pub := make(chan message)
	unsub := make(chan user)
	// initialize the chatroom to receive and redistribute messages
	go chatroom(sub, unsub, pub)

	r := mux.NewRouter()
	r.HandleFunc("/room", roomHandler(sub, unsub, pub))

	log.Println("Starting on port 5000")
	log.Fatal(http.ListenAndServe(getPort(), r))
}
