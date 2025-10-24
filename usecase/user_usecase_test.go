package usecase

import (
	"reflect"
	"testing"

	"github.com/To-ge/gr_backend_go/domain/entity"
	"github.com/To-ge/gr_backend_go/domain/mock_repository"
	"github.com/To-ge/gr_backend_go/usecase/model"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
)

var (
	id       = uuid.New()
	name     = "testname"
	password = "testpass"
	email    = "test@mail"
	isAdmin  = true
)

func TestCreateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	arg1 := entity.User{
		Name:     name,
		Password: password,
		Email:    email,
		IsAdmin:  isAdmin,
	}
	var err error

	mockSample := mock_repository.NewMockIUserRepository(ctrl)
	mockSample.EXPECT().CreateUser(arg1).Return(err)

	arg2 := &model.CreateUserInput{
		Name:     name,
		Password: password,
		Email:    email,
		IsAdmin:  isAdmin,
	}
	expected := &model.CreateUserOutput{
		Name:  name,
		Email: email,
	}
	userUsecase := NewUserUsecase(mockSample)
	got, err := userUsecase.CreateUser(arg2)

	if err != nil {
		t.Errorf("CreateUser() err = %v, want nil", err)
	}
	if !reflect.DeepEqual(*got, *expected) {
		t.Errorf("CreateUser() got = %v, want %v", got, expected)
	}
}

func TestFindOneById(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testUser := &entity.User{
		ID:       id,
		Name:     name,
		Password: password,
		Email:    email,
		IsAdmin:  isAdmin,
	}
	var err error

	mockSample := mock_repository.NewMockIUserRepository(ctrl)
	mockSample.EXPECT().FindOneById(id.String()).Return(testUser, err)

	arg := id.String()
	expected := &model.User{
		ID:      id,
		Name:    name,
		Email:   email,
		IsAdmin: isAdmin,
	}

	userUsecase := NewUserUsecase(mockSample)
	got, err := userUsecase.FindOneById(arg)

	if err != nil {
		t.Errorf("FindOneById() err = %v, want nil", err)
	}
	if !reflect.DeepEqual(*got, *expected) {
		t.Errorf("FindOneById() got = %v, want %v", got, expected)
	}
}
