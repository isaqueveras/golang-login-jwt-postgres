package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/isaqueveras/golang-login-jwt-postgres/controllers"
)

func Setup(app *fiber.App) {
	app.Post("/api/register", controllers.Register)
	app.Post("/api/login", controllers.Login)
	app.Post("/api/logout", controllers.Logout)
	app.Get("/api/user", controllers.User)
}
