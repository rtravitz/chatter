package room

import "github.com/google/uuid"

// User represents a chatroom user
type User struct {
	ID       uuid.UUID
	Messages chan Message
}

// NewUser creates a new user
func NewUser() User {
	id := uuid.New()
	mc := make(chan Message)
	return User{ID: id, Messages: mc}
}

// Find gets a user out of a list based on the provided uuid
func findUser(id uuid.UUID, users []User) (int, bool) {
	for i, u := range users {
		if u.ID == id {
			return i, true
		}
	}

	return 0, false
}
