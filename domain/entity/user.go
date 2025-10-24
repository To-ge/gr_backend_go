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

func NewUser(name string, email string, password string, isAdmin bool) *User {
	return &User{
		Name:     name,
		Email:    email,
		Password: password,
		IsAdmin:  isAdmin,
	}
}

func (u *User) IsValidEmail() bool {
	return strings.Contains(u.Email, "@")
}
