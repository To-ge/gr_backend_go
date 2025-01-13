package repository

type IAuthenticationRepository interface {
	SignIn(key string) error
	SignOut(key string) error
}
