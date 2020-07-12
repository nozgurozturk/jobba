package token

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	errors "github.com/nozgurozturk/jobba/utils"
	"os"
	"strconv"
)

type Model struct {
	Token   string
	Uuid    string
	Expires int64
	UserId  string
}

type Response struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

func VerifyToken(tokenString string, tokenType string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {

		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		if tokenType == "refreshToken" {
			return []byte(os.Getenv("REFRESH_SECRET")), nil
		}
		if tokenType == "accessToken" {
			return []byte(os.Getenv("ACCESS_SECRET")), nil
		}
		return nil, fmt.Errorf("unexpected token: %v", tokenType)
	})
	if err != nil {
		return nil, err
	}
	return token, nil
}

func ValidateToken(tokenString string, tokenType string) *errors.ApplicationError {
	token, err := VerifyToken(tokenString, tokenType)
	if err != nil {
		return errors.Unauthorized(err.Error())
	}
	if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
		return errors.Unauthorized(err.Error())
	}
	return nil
}

func ExtractTokenClaims(tokenString string, tokenType string) (*Model, *errors.ApplicationError) {
	token, err := VerifyToken(tokenString, tokenType)
	if err != nil {
		return nil, errors.InternalServer(err.Error())
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		uuid, ok := claims["uuid"].(string)
		if !ok {
			return nil, errors.InternalServer("can't access uuid in token claims")
		}
		expires, parseErr := strconv.ParseInt(fmt.Sprintf("%.f", claims["exp"]), 10, 64)
		if parseErr != nil {
			return nil, errors.InternalServer("can't access uuid in token claims")
		}
		userId, ok := claims["userId"].(string)
		if !ok {
			return nil, errors.InternalServer("can't user id in token claims")
		}

		return &Model{
			Uuid:    uuid,
			Expires: expires,
			UserId:  userId,
		}, nil
	}
	return nil, errors.Unauthorized(err.Error())
}
