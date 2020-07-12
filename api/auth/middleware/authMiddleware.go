package middleware

import (
	"github.com/gofiber/fiber"
	"github.com/nozgurozturk/jobba/api/auth/domain/token"
	errors "github.com/nozgurozturk/jobba/utils"
	"net/http"
)

func AuthMiddleware() func(*fiber.Ctx) {
	return func(ctx *fiber.Ctx) {
		// Except Login and Signup
		if ctx.Path() == "/api/login" || ctx.Path() == "/api/signup" {
			ctx.Next()
			return
		}
		at := ctx.Cookies("accessToken")
		if at == "" {
			ctx.Status(http.StatusUnauthorized).JSON(errors.Unauthorized("Access token is not found"))
			return
		}

		if errAt := token.ValidateToken(at, "accessToken"); errAt != nil {
			ctx.Status(errAt.Status).JSON(errAt)
			return
		}

		rt := ctx.Cookies("refreshToken")
		if rt == "" {
			ctx.Status(http.StatusUnauthorized).JSON(errors.Unauthorized("Refresh token is not found"))
			return
		}

		if errRt := token.ValidateToken(rt, "refreshToken"); errRt != nil {
			ctx.Status(errRt.Status).JSON(errRt)
			return
		}

		ctx.Next()
	}
}
