package model

import "net/http"

// SignIn
type SignInInput struct {
	ResponseWriter http.ResponseWriter
	Request        *http.Request
	Email          string `json:"email"`
	Password       string `json:"password"`
}
type SignInOutput struct{}

// SignOut
type SignOutInput struct {
	ResponseWriter http.ResponseWriter
	Request        *http.Request
}
type SignOutOutput struct{}

type SessionCheckInput struct{}
type SessionCheckOutput struct{}
