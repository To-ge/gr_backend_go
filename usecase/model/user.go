package model

import "github.com/google/uuid"

type CreateUserInput struct {
	Name     string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
	IsAdmin  bool   `json:"is_admin"`
}

type CreateUserOutput struct {
	Name  string `json:"username"`
	Email string `json:"email"`
}

type User struct {
	ID      uuid.UUID
	Name    string
	Email   string
	IsAdmin bool
}
