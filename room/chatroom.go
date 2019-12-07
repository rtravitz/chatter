package room

import (
	"fmt"
)

// Message represents a chatroom message
type Message struct {
	MessageType int
	Data        []byte
	Sender      User
}

// New creates a new chatroom with subscribe, unsubscribe, and publish hooks
func New(subscribe chan User, unsubscribe chan User, publish chan Message) {
	var subs []User

	for {
		select {
		case sub := <-subscribe:
			fmt.Printf("Adding subr: %+v\n", sub)
			subs = append(subs, sub)
		case mes := <-publish:
			fmt.Printf("Sending message to %d subs: %+v\n", len(subs), mes)
			for _, sub := range subs {
				sub.Messages <- mes
			}
		case unsub := <-unsubscribe:
			fmt.Printf("Unsubscribing user: %v\nCurrent sub list: %v\n", unsub.ID, subs)
			if match, ok := findUser(unsub.ID, subs); ok {
				subs[match] = subs[len(subs)-1]
				subs = subs[:len(subs)-1]
			}
			fmt.Printf("After sub list: %v\n", subs)
		}
	}
}
