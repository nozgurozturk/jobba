package controller

import (
	"github.com/gofiber/fiber"
	"github.com/nozgurozturk/jobba/api/auth/service"
	"github.com/nozgurozturk/jobba/api/user/domain/user"
	errors "github.com/nozgurozturk/jobba/utils"
	"net/http"
	"time"
)

func SignUp(c *fiber.Ctx) {

	var u user.Model
	if err := c.BodyParser(&u); err != nil {
		jsonError := errors.BadRequest("invalid json body")
		c.Status(jsonError.Status).JSON(jsonError)
		return
	}
	tokens, err := service.Auth.SignUp(u)
	if err != nil {
		c.Status(err.Status).JSON(err)
		return
	}

	accessToken := new(fiber.Cookie)
	accessToken.Name = "accessToken"
	accessToken.Value = tokens.AccessToken
	accessToken.Expires = time.Now().Add(time.Minute * 15)
	accessToken.HTTPOnly = true
	c.Cookie(accessToken)

	refreshToken := new(fiber.Cookie)
	refreshToken.Name = "refreshToken"
	refreshToken.Value = tokens.RefreshToken
	refreshToken.Expires = time.Now().Add(time.Hour * 24 * 7)
	refreshToken.HTTPOnly = true
	c.Cookie(refreshToken)

	c.Status(http.StatusOK).JSON(tokens)
}

func Login(c *fiber.Ctx) {
	var u user.Request
	if err := c.BodyParser(&u); err != nil {
		jsonError := errors.BadRequest("invalid json body")
		c.Status(jsonError.Status).JSON(jsonError)
		return
	}
	tokens, err := service.Auth.Login(u)
	if err != nil {
		c.Status(err.Status).JSON(err)
		return
	}

	accessToken := new(fiber.Cookie)
	accessToken.Name = "accessToken"
	accessToken.Value = tokens.AccessToken
	accessToken.Expires = time.Now().Add(time.Minute * 15)
	accessToken.HTTPOnly = true
	c.Cookie(accessToken)

	refreshToken := new(fiber.Cookie)
	refreshToken.Name = "refreshToken"
	refreshToken.Value = tokens.RefreshToken
	refreshToken.Expires = time.Now().Add(time.Hour * 24 * 7)
	refreshToken.HTTPOnly = true
	c.Cookie(refreshToken)

	c.Status(http.StatusOK).JSON(tokens)
}
