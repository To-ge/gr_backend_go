package entity

import (
	"strings"

	"github.com/google/uuid"
)

type User struct {
	ID       uuid.UUID
	Name     string
	Password string
	Email    string
	IsAdmin  bool
}

func NewUser(name string, email string) *User {
	return &User{
		Name:  name,
		Email: email,
	}
}

func (u *User) IsValidEmail() bool {
	return strings.Contains(u.Email, "@")
}
