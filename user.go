package main

import "github.com/google/uuid"

type user struct {
	id       uuid.UUID
	messages chan message
}

func newUser() user {
	id := uuid.New()
	mc := make(chan message)
	return user{id: id, messages: mc}
}

func findUser(id uuid.UUID, users []user) (int, bool) {
	for i, u := range users {
		if u.id == id {
			return i, true
		}
	}

	return 0, false
}
