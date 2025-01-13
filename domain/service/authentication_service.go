package service

import (
	"fmt"

	"github.com/To-ge/gr_backend_go/domain/repository"
)

type IAuthenticationService interface {
	SignIn(username string, password string) (string, error)
	SignOut(sessionKey string) error
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

func (as *authenticationService) SignIn(username string, password string) (string, error) {
	user, err := as.userRepo.FindOne(username, password)
	if err != nil {
		return "", err
	}
	if user == nil {
		return "", fmt.Errorf("user is not found")
	}
	key := "session:" + string(user.ID)
	if err := as.authRepo.SignIn(key); err != nil {
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
