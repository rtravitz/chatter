package user

import (
	"time"
)

type User struct {
	ID           string
	Name         string
	Email        string
	PasswordHash []byte
	DateCreated  time.Time
	DateUpdated  time.Time
}

type NewUser struct {
	Name            string
	Email           string
	Password        string
	PasswordConfirm string
}
