package service

import (
	"github.com/nozgurozturk/jobba/api/user/domain/user"
	errors "github.com/nozgurozturk/jobba/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type userService struct{}

type userServiceInterface interface {
	Create(user.Model) (*user.Model, *errors.ApplicationError)
	FindByName(string) (*user.Model, *errors.ApplicationError)
	FindByUserID(string) (*user.Model, *errors.ApplicationError)
	FindByEmail(string) (*user.Model, *errors.ApplicationError)
}

var (
	User userServiceInterface = &userService{}
)

func (s *userService) Create(u user.Model) (*user.Model, *errors.ApplicationError) {

	if err := u.Validate(); err != nil {
		return nil, err
	}

	hashedPassword, err := user.HashPassword(u.Password)
	if err != nil {
		return nil, errors.InternalServer(err.Error())
	}

	u.Password = string(hashedPassword)
	u.CreatedAt = time.Now().Unix()
	if err := u.Create(); err != nil {
		return nil, err
	}

	return &u, nil
}
func (s *userService) FindByName(userName string) (*user.Model, *errors.ApplicationError) {

	if userName == "" {
		return nil, errors.BadRequest("user name can not be empty")
	}
	u := &user.Model{UserName: userName}
	currentUser, errResult := u.FindByName();
	if errResult != nil {
		return nil, errResult
	}

	return currentUser, nil
}

func (s *userService) FindByUserID(userId string) (*user.Model, *errors.ApplicationError) {

	_id, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return nil, errors.BadRequest("user id can not be empty")
	}
	u := &user.Model{ID: _id}
	currentUser, errResult := u.FindByUserID();
	if errResult != nil {
		return nil, errResult
	}

	return currentUser, nil
}

func (s *userService) FindByEmail(email string) (*user.Model, *errors.ApplicationError) {

	if email == "" {
		return nil, errors.BadRequest("Email can not be empty")
	}
	u := &user.Model{Email: email}
	currentUser, errResult := u.FindByEmail();
	if errResult != nil {
		return nil, errResult
	}

	return currentUser, nil
}




