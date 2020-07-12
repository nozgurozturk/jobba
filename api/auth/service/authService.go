package service

import (
	"github.com/nozgurozturk/jobba/api/auth/domain/token"
	"github.com/nozgurozturk/jobba/api/user/domain/user"
	"github.com/nozgurozturk/jobba/api/user/service"
	errors "github.com/nozgurozturk/jobba/utils"
)

type authService struct{}

type authServiceInterface interface {
	Login(user.Request) (*token.Response, *errors.ApplicationError)
	SignUp(user.Model) (*token.Response, *errors.ApplicationError)
	GetUserID(string) (string, *errors.ApplicationError)
	ExtractUserID(string) (string, *errors.ApplicationError)
}

var (
	Auth authServiceInterface = &authService{}
)

func (s *authService) Login(u user.Request) (*token.Response, *errors.ApplicationError) {

	if u.Password == "" {
		return nil, errors.BadRequest("password is required")
	}

	if u.UserName == "" {
		return nil, errors.BadRequest("username is required")
	}
	currentUser, err := service.User.FindByName(u.UserName)

	if err != nil {
		return nil, err
	}

	if passErr := user.VerifyPassword(currentUser.Password, u.Password); passErr != nil {
		return nil, errors.BadRequest("Password is not correct")
	}
	userId := currentUser.ID.Hex()
	accessToken, atError := token.CreateAccessToken(userId)
	if atError != nil {
		return nil, atError
	}

	refreshToken, rtError := token.CreateRefreshToken(userId)
	if rtError != nil {
		return nil, rtError
	}

	if err := token.SetToken(refreshToken); err != nil {
		return nil, err
	}

	return &token.Response{
		AccessToken:  accessToken.Token,
		RefreshToken: refreshToken.Token,
	}, nil
}

func (s *authService) SignUp(u user.Model) (*token.Response, *errors.ApplicationError) {
	existingUser, _ := service.User.FindByEmail(u.Email)

	if existingUser != nil {
		return nil, errors.AlreadyExist("Email is already exist")
	}

	existingUser, _ = service.User.FindByName(u.UserName)

	if existingUser != nil {
		return nil, errors.AlreadyExist("User is already exist")
	}

	createdUser, err := service.User.Create(u)
	if err != nil {
		return nil, err
	}

	userID := createdUser.ID.Hex()
	accessToken, atError := token.CreateAccessToken(userID)
	if atError != nil {
		return nil, atError
	}

	refreshToken, rtError := token.CreateRefreshToken(userID)
	if rtError != nil {
		return nil, rtError
	}
	if err := token.SetToken(refreshToken); err != nil {
		return nil, err
	}

	return &token.Response{
		AccessToken:  accessToken.Token,
		RefreshToken: refreshToken.Token,
	}, nil
}

func (s *authService) ExtractUserID(tokenString string)(string, *errors.ApplicationError) {
	if tokenString == "" {
		return "", errors.Unauthorized("token can not be empty")
	}
	t, err := token.ExtractTokenClaims(tokenString, "accessToken")
	if err != nil {
		return "", err
	}
	return t.UserId, nil
}

func (s *authService) GetUserID(tokenString string)(string, *errors.ApplicationError) {
	if tokenString == "" {
		return "", errors.Unauthorized("token is not found")
	}
	userID, err := token.GetUserId(tokenString)
	if err != nil {
		return "", err
	}

	return userID, nil
}
