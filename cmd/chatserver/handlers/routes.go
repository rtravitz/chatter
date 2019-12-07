package handlers

import (
	"github.com/gorilla/mux"
	"github.com/rtravitz/chatter/room"
	"net/http"
)

// API establishes the routes for the chatserver
func API() http.Handler {
	sub := make(chan room.User)
	pub := make(chan room.Message)
	unsub := make(chan room.User)
	// initialize the chatroom to receive and redistribute messages
	go room.New(sub, unsub, pub)

	ro := Room{}

	r := mux.NewRouter()
	r.HandleFunc("/room", ro.Listen(sub, unsub, pub))

	return r
}
