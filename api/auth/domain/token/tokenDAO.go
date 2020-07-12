package token

import (
	"context"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	tokenDB "github.com/nozgurozturk/jobba/api/auth/database/redis"
	errors "github.com/nozgurozturk/jobba/utils"
	"os"
	"time"
)

func CreateAccessToken(userId string) (*Model, *errors.ApplicationError) {
	accessToken := &Model{}
	e := godotenv.Load()
	if e != nil {
		fmt.Print(e)
	}

	AccessSecret := os.Getenv("ACCESS_SECRET")

	accessToken.Expires = time.Now().Add(time.Minute * 15).Unix()
	accessToken.Uuid = uuid.New().String()
	accessToken.UserId = userId

	claims := jwt.MapClaims{
		"exp":    accessToken.Expires,
		"uuid":   accessToken.Uuid,
		"userId": accessToken.UserId,
	}
	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	var err error
	accessToken.Token, err = rt.SignedString([]byte(AccessSecret))
	if err != nil {
		return nil, errors.InternalServer("Access Token can't created")
	}
	return accessToken, nil
}

func CreateRefreshToken(userId string) (*Model, *errors.ApplicationError) {
	refreshToken := &Model{}
	e := godotenv.Load()
	if e != nil {
		fmt.Print(e)
	}

	RefreshSecret := os.Getenv("REFRESH_SECRET")

	refreshToken.Expires = time.Now().Add(time.Hour * 24 * 7).Unix()
	refreshToken.Uuid = uuid.New().String()
	refreshToken.UserId = userId

	claims := jwt.MapClaims{
		"exp":    refreshToken.Expires,
		"uuid":   refreshToken.Uuid,
		"userId": refreshToken.UserId,
	}
	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	var err error
	refreshToken.Token, err = rt.SignedString([]byte(RefreshSecret))
	if err != nil {
		return nil, errors.InternalServer("Refresh Token can't created")
	}
	return refreshToken, nil
}

func SetToken(token *Model) *errors.ApplicationError {
	expires := time.Unix(token.Expires, 0)
	now := time.Now()
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err := tokenDB.Client.Set(ctx, token.Uuid, token.UserId, expires.Sub(now)).Err()
	if err != nil {
		return errors.InternalServer("Error when setting access token")
	}
	return nil
}

func GetUserId(token string) (string, *errors.ApplicationError) {
	claims, claimErr := ExtractTokenClaims(token, "refreshToken")
	if claimErr != nil {
		return "", claimErr
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	userid, err := tokenDB.Client.Get(ctx, claims.Uuid).Result()
	if err != nil {
		return "", errors.Unauthorized(err.Error())
	}
	userID := userid
	return userID, nil
}
