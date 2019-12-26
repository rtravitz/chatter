package user

import (
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
	"time"
)

// User is someone with access to the system
type User struct {
	ID           string    `db:"user_id" json:"id"`
	Name         string    `db:"name" json:"name"`
	Email        string    `db:"email" json:"email"`
	PasswordHash []byte    `db:"password_hash" json:"-"`
	DateCreated  time.Time `db:"date_created" json:"date_created"`
	DateUpdated  time.Time `db:"date_updated" json:"date_updated"`
}

// NewUser contains information needed to create a new user
type NewUser struct {
	Name            string
	Email           string
	Password        string
	PasswordConfirm string
}

// Create inserts a new user in the database
func Create(db *sqlx.DB, n NewUser, now time.Time) (*User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(n.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	u := User{
		ID:           uuid.New().String(),
		Name:         n.Name,
		Email:        n.Email,
		PasswordHash: hash,
		DateCreated:  now.UTC(),
		DateUpdated:  now.UTC(),
	}

	const q = `INSERT INTO users
		(user_id, name, email, password_hash, date_created, date_updated)
		VALUES ($1, $2, $3, $4, $5, $6)`
	_, err = db.Exec(
		q,
		u.ID, u.Name, u.Email,
		u.PasswordHash, u.DateCreated, u.DateUpdated,
	)

	if err != nil {
		return nil, err
	}

	return &u, nil
}
