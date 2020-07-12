package app

import (
	authController "github.com/nozgurozturk/jobba/api/auth/controller"
	jobController "github.com/nozgurozturk/jobba/api/jobs/controller"
	userController "github.com/nozgurozturk/jobba/api/user/controller"
)

func (s *Server) InitRoutes() {
	prefix := s.Router.Group("/api")
	// Auth
	prefix.Post("/login", authController.Login)
	prefix.Post("/signup", authController.SignUp)
	// User
	prefix.Get("/user", userController.Info)
	prefix.Post("/user", userController.Create)
	// Job
	prefix.Get("/jobs", jobController.FindAll)
	prefix.Post("/job", jobController.FindOne)
	prefix.Post("/job/create", jobController.Create)
	prefix.Put("/job", jobController.Update)
}
