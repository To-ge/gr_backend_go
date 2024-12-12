package model

type CreateUserInput struct {
	Name  string `json:"username"`
	Email string `json:"email"`
}

type CreateUserOutput struct {
	ID    string `json:"id"`
	Name  string `json:"username"`
	Email string `json:"email"`
}
