package service

import (
	"fmt"

	"github.com/To-ge/gr_backend_go/domain/repository"
	"github.com/google/uuid"
)

type IAuthenticationService interface {
	SignIn(email string, password string) (string, error)
	SignOut(sessionKey string) error
	RefreshSessionExpiration(sessionKey string) error
}

type authenticationService struct {
	userRepo repository.IUserRepository
	authRepo repository.IAuthenticationRepository
}

func NewAuthenticationService(userRepo repository.IUserRepository, authRepo repository.IAuthenticationRepository) IAuthenticationService {
	return &authenticationService{
		userRepo: userRepo,
		authRepo: authRepo,
	}
}

func (as *authenticationService) SignIn(email string, password string) (string, error) {
	user, err := as.userRepo.FindOne(email, password)
	if err != nil {
		return "", err
	}
	if user == nil {
		return "", fmt.Errorf("user is not found")
	}
	key := "session:" + uuid.NewString()
	if err := as.authRepo.SignIn(key, user.ID.String()); err != nil {
		return "", err
	}
	return key, nil
}

func (as *authenticationService) SignOut(sessionKey string) error {
	if err := as.authRepo.SignOut(sessionKey); err != nil {
		return err
	}
	return nil
}

func (as *authenticationService) RefreshSessionExpiration(sessionKey string) error {
	if err := as.authRepo.RefreshSessionExpiration(sessionKey); err != nil {
		return err
	}
	return nil
}
