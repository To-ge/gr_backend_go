package entity

import (
	"strings"

	"github.com/google/uuid"
)

type User struct {
	ID       string
	Name     string
	Password string
	Email    string
}

func NewUser(name string, email string) *User {
	uuid := uuid.NewString()
	return &User{
		ID:    uuid,
		Name:  name,
		Email: email,
	}
}

func (u *User) IsValidEmail() bool {
	return strings.Contains(u.Email, "@")
}
