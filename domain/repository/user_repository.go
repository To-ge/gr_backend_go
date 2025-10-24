//go:generate mockgen -source=user_repository.go -destination=../mock_repository/mock_user_repository.go -package=mock_repository

package repository

import "github.com/To-ge/gr_backend_go/domain/entity"

type IUserRepository interface {
	CreateUser(entity.User) error
	FindOne(email string, password string) (*entity.User, error)
	FindOneById(id string) (*entity.User, error)
}
