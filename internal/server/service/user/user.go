package user

import (
	"errors"
	"regexp"
	"time"
)

var (
	ErrNotFound = errors.New("user not found")
	ErrBadEmail = errors.New("bad email")
)

type User struct {
	ID          string    `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	DisplayName string    `json:"display_name"`
	Email       string    `json:"email"`
}

func NewUser(name, email string) (User, error) {
	if !isValidEmail(email) {
		return User{}, ErrBadEmail
	}

	user := User{
		DisplayName: name,
		Email:       email,
	}

	return user, nil
}

func isValidEmail(email string) bool {
	regex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	r := regexp.MustCompile(regex)
	return r.MatchString(email)
}
