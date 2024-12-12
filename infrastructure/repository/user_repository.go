package repository

import (
	"fmt"

	"github.com/To-ge/gr_backend_go/domain/entity"
	domainRepo "github.com/To-ge/gr_backend_go/domain/repository"
	"github.com/To-ge/gr_backend_go/infrastructure/database"
	"github.com/To-ge/gr_backend_go/infrastructure/database/model"
)

type userRepository struct {
	dbc *database.DBConnector
}

func NewUserRepository(dbc *database.DBConnector) domainRepo.IUserRepository {
	return userRepository{
		dbc: dbc,
	}
}

func (ur userRepository) CreateUser(user entity.User) error {
	record := &model.User{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}
	if err := ur.dbc.Conn.Create(record).Error; err != nil {
		return fmt.Errorf("new user can't create, %s", err.Error())
	}
	return nil
}
