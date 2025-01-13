package model

import "net/http"

// SignIn
type SignInInput struct {
	ResponseWriter http.ResponseWriter
	Request        *http.Request
	Username       string `json:"username"`
	Password       string `json:"password"`
}
type SignInOutput struct{}

// SignOut
type SignOutInput struct {
	ResponseWriter http.ResponseWriter
	Request        *http.Request
}
type SignOutOutput struct{}
