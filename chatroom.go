package main

import "fmt"

type message struct {
	messageType int
	data        []byte
	sender      user
}

func chatroom(subscribe chan user, unsubscribe chan user, publish chan message) {
	var subs []user

	for {
		select {
		case sub := <-subscribe:
			fmt.Printf("Adding subr: %+v\n", sub)
			subs = append(subs, sub)
		case mes := <-publish:
			fmt.Printf("Sending message to %d subs: %+v\n", len(subs), mes)
			for _, sub := range subs {
				sub.messages <- mes
			}
		case unsub := <-unsubscribe:
			fmt.Printf("Unsubscribing user: %v\nCurrent sub list: %v\n", unsub.id, subs)
			if match, ok := findUser(unsub.id, subs); ok {
				subs[match] = subs[len(subs)-1]
				subs = subs[:len(subs)-1]
			}
			fmt.Printf("After sub list: %v\n", subs)
		}
	}
}
