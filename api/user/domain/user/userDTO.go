package user

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
	errors "github.com/nozgurozturk/jobba/utils"
	"regexp"
	"strings"
)

type Model struct {
	ID        primitive.ObjectID  `bson:"_id,omitempty" json:"id"`
	UserName  string `bson:"userName,omitempty" json:"userName"`
	Email     string `bson:"email,omitempty" json:"email"`
	Password  string `bson:"password,omitempty" json:"password"`
	CreatedAt int64 `bson:"createdAt" json:"createdAt"`
}

type Request struct {
	UserName string `json:"userName"`
	Password string `json:"password"`
}

type Response struct {
	UserName string `json:"userName"`
	Email    string `json:"email"`
}

func (user *Model) Validate() *errors.ApplicationError {

	user.UserName = strings.TrimSpace(user.UserName)
	if user.UserName == "" {

		return errors.BadRequest("User name is required")
	}

	if len(user.UserName) < 6 {
		return errors.BadRequest("user name must be minimum six character")
	}

	user.Password = strings.TrimSpace(user.Password)
	if user.Password == "" {
		return errors.BadRequest("password is required")
	}

	if len(user.Password) < 8 {
		return errors.BadRequest("password must be minimum eight character")
	}

	user.Email = strings.TrimSpace(user.Email)
	emailRegex := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

	if user.Email == "" {
		return errors.BadRequest("email is required")
	}

	if !emailRegex.MatchString(user.Email) {
		return errors.BadRequest("email address is not valid")
	}

	return nil
}

func HashPassword(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

func VerifyPassword(hashedPassword string, userPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(userPassword))
}

func (user *Model) UserResponse() interface{} {
	return Response{
		UserName: user.UserName,
		Email:    user.Email,
	}
}
