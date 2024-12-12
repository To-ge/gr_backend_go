package repository

import "github.com/To-ge/gr_backend_go/domain/entity"

type IUserRepository interface {
	CreateUser(entity.User) error
}
