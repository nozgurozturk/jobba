package controller

import (
	"github.com/gofiber/fiber"
	authService "github.com/nozgurozturk/jobba/api/auth/service"
	"github.com/nozgurozturk/jobba/api/jobs/domain/job"
	"github.com/nozgurozturk/jobba/api/jobs/service"
	errors "github.com/nozgurozturk/jobba/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
)

func Create(c *fiber.Ctx) {
	var j job.Model

	if err := c.BodyParser(&j); err != nil {
		jsonErr := errors.BadRequest("invalid job body")
		c.Status(jsonErr.Status).JSON(jsonErr)
		return
	}
	rt := c.Cookies("refreshToken")
	userID, err := authService.Auth.GetUserID(rt)
	if err != nil {
		c.Status(err.Status).JSON(err)
		return
	}

	var idErr error
	j.UserID, idErr = primitive.ObjectIDFromHex(userID)
	if idErr != nil {
		c.Status(http.StatusUnauthorized).JSON(errors.Unauthorized("Unauthorized user"))
		return
	}
	createdJob, err := service.Job.Create(j)

	if err != nil {
		c.Status(err.Status).JSON(err)
		return
	}
	c.Status(http.StatusCreated).JSON(createdJob.JobResponse())
}

func FindOne(c *fiber.Ctx) {
	var job job.Request
	if err := c.BodyParser(&job); err != nil {
		jsonErr := errors.BadRequest("Invalid job body")
		c.Status(jsonErr.Status).JSON(jsonErr)
		return
	}

	rt := c.Cookies("refreshToken")
	userID, err := authService.Auth.GetUserID(rt)
	if err != nil {
		c.Status(err.Status).JSON(err)
		return
	}

	foundJob, err := service.Job.FindOne(job.ID)
	jobUserID := foundJob.UserID.Hex()
	if (jobUserID != userID) {
		c.Status(http.StatusForbidden).JSON(errors.Forbiden("You don't have access"))
		return
	}
	if err != nil {
		c.Status(err.Status).JSON(err)
		return
	}
	c.Status(http.StatusFound).JSON(foundJob.JobResponse())
}

func FindAll(c *fiber.Ctx) {
	rt := c.Cookies("refreshToken")
	if rt == "" {
		c.Status(http.StatusUnauthorized).JSON(errors.Unauthorized("Unauthorized user"))
		return
	}
	userID, err := authService.Auth.GetUserID(rt)
	if err != nil {
		c.Status(err.Status).JSON(err)
		return
	}

	jobs, err := service.Job.FindAll(userID)
	if err != nil {
		c.Status(err.Status).JSON(err)
	}

	c.Status(http.StatusFound).JSON(jobs)
}

func Update(c *fiber.Ctx) {
	var job job.Request
	if err := c.BodyParser(&job); err != nil {
		jsonErr := errors.BadRequest("Invalid job body")
		c.Status(jsonErr.Status).JSON(jsonErr)
		return
	}
	rt := c.Cookies("refreshToken")
	if rt == "" {
		c.Status(http.StatusUnauthorized).JSON(errors.Unauthorized("Unauthorized user"))
		return
	}
	userID, err := authService.Auth.GetUserID(rt)
	if err != nil {
		c.Status(err.Status).JSON(err)
		return
	}
	foundJob, err := service.Job.FindOne(job.ID)
	if err != nil {
		c.Status(err.Status).JSON(err)
		return
	}
	jobUserID := foundJob.UserID.Hex()

	if (jobUserID != userID) {
		c.Status(http.StatusForbidden).JSON(errors.Forbiden("You don't have access"))
		return
	}

	updatedJob, err := service.Job.Update(&job)

	if err != nil {
		c.Status(err.Status).JSON(err)
	}

	c.Status(http.StatusOK).JSON(updatedJob.JobResponse())

}