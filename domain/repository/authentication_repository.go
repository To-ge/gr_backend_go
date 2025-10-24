//go:generate mockgen -source=authentication_repository.go -destination=../mock_repository/mock_authentication_repository.go -package=mock_repository

package repository

type IAuthenticationRepository interface {
	SignIn(key string, value string) error
	SignOut(key string) error
	RefreshSessionExpiration(key string) error
}
