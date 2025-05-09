package usecase

import (
	"github.com/To-ge/gr_backend_go/domain/entity"
	"github.com/To-ge/gr_backend_go/domain/repository"
	"github.com/To-ge/gr_backend_go/usecase/model"
)

type IUserUsecase interface {
	CreateUser(input *model.CreateUserInput) (*model.CreateUserOutput, error)
	FindOneById(id string) (*model.User, error)
}

type userUsecase struct {
	repo repository.IUserRepository
}

func NewUserUsecase(ur repository.IUserRepository) IUserUsecase {
	return &userUsecase{
		repo: ur,
	}
}

func (uu *userUsecase) CreateUser(input *model.CreateUserInput) (*model.CreateUserOutput, error) {
	user := entity.NewUser(input.Name, input.Email)
	if err := uu.repo.CreateUser(*user); err != nil {
		return nil, err
	}

	return &model.CreateUserOutput{
		Name:  user.Name,
		Email: user.Email,
	}, nil
}

func (uu *userUsecase) FindOneById(id string) (*model.User, error) {
	user, err := uu.repo.FindOneById(id)
	if err != nil {
		return nil, err
	}

	return &model.User{
		ID:      user.ID,
		Name:    user.Name,
		Email:   user.Email,
		IsAdmin: user.IsAdmin,
	}, nil
}
