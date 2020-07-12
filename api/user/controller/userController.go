package controller

import (
	"github.com/gofiber/fiber"
	"github.com/nozgurozturk/jobba/api/user/domain/user"
	"github.com/nozgurozturk/jobba/api/user/service"
	authService "github.com/nozgurozturk/jobba/api/auth/service"
	errors "github.com/nozgurozturk/jobba/utils"
	"net/http"
)

func Create(c *fiber.Ctx) {
	var u user.Model

	if err := c.BodyParser(&u); err != nil {
		jsonError := errors.BadRequest("invalid json body")
		c.Status(jsonError.Status).JSON(jsonError)
		return
	}

	createdUser, serviceErr := service.User.Create(u)
	if serviceErr != nil {
		c.Status(serviceErr.Status).JSON(serviceErr)
		return
	}

	c.Status(http.StatusCreated).JSON(createdUser.UserResponse())
}

func Info(c *fiber.Ctx) {
	at := c.Cookies("accessToken")
	if at == ""{
		c.Status(http.StatusBadRequest).JSON(errors.BadRequest("Unauthorized User"))
		return
	}

	userID, err := authService.Auth.GetUserID(at)
	if err != nil {
		c.Status(http.StatusBadRequest).JSON(errors.BadRequest("Unauthorized User"))
	}

	userInfo, infoErr := service.User.FindByUserID(userID)

	if infoErr != nil {
		c.Status(infoErr.Status).JSON(infoErr)
		return
	}

	c.Status(http.StatusOK).JSON(userInfo.UserResponse())
}