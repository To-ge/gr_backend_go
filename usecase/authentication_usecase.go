package usecase

import (
	"fmt"

	"github.com/To-ge/gr_backend_go/config"
	"github.com/To-ge/gr_backend_go/domain/service"
	"github.com/To-ge/gr_backend_go/usecase/model"
	"github.com/gorilla/sessions"
)

type IAuthenticationUsecase interface {
	SignIn(*model.SignInInput) error
	SignOut(*model.SignOutInput) error
	RefreshSessionExpiration(key string) error
}

type authenticationUsecase struct {
	srv   service.IAuthenticationService
	store *sessions.CookieStore
}

func NewAuthenticationUsecase(as service.IAuthenticationService) IAuthenticationUsecase {
	return &authenticationUsecase{
		srv:   as,
		store: config.SessionStore,
	}
}

func (au *authenticationUsecase) SignIn(input *model.SignInInput) error {
	key, err := au.srv.SignIn(input.Email, input.Password)
	if err != nil {
		return err
	}
	sess, err := au.store.Get(input.Request, "session")
	if err != nil {
		return err
	}
	sess.Values[config.SessionKey] = key
	return sess.Save(input.Request, input.ResponseWriter)
}

func (au *authenticationUsecase) SignOut(input *model.SignOutInput) error {
	sess, err := au.store.Get(input.Request, "session")
	if err != nil {
		return err
	}
	value := sess.Values[config.SessionKey]
	if value != nil {
		// key := strings.TrimPrefix(value.(string), "session:")
		key := value.(string)
		au.srv.SignOut(key)
		sess.Options.MaxAge = -1
	} else {
		return fmt.Errorf("invalid cookie argument")
	}

	return sess.Save(input.Request, input.ResponseWriter)
}

func (au *authenticationUsecase) RefreshSessionExpiration(key string) error {
	if err := au.srv.RefreshSessionExpiration(key); err != nil {
		return err
	}
	return nil
}
