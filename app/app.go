package app

import (
	"fmt"
	"github.com/gofiber/fiber"
	"github.com/nozgurozturk/jobba/api/auth/middleware"
)

type Server struct {
	Router *fiber.App
}

func (s *Server) Run() {
	s.Router = fiber.New()
	s.Router.Use(middleware.AuthMiddleware())

	s.InitRoutes()
	s.Router.Settings = &fiber.Settings{DisableStartupMessage: true}
	fmt.Println("Started listening on 0.0.0.0:8081")
	s.Router.Listen(8081)
}
