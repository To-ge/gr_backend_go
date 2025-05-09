package repository

import "github.com/To-ge/gr_backend_go/domain/entity"

type IUserRepository interface {
	CreateUser(entity.User) error
	FindOne(email string, password string) (*entity.User, error)
	FindOneById(id string) (*entity.User, error)
}
