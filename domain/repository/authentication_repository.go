package repository

type IAuthenticationRepository interface {
	SignIn(key string, value string) error
	SignOut(key string) error
	RefreshSessionExpiration(key string) error
}
