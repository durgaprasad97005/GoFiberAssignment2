package routes

import (
	"github.com/durgaprasad97005/GoFiberAssignment2/src/controllers"
	"github.com/gofiber/fiber/v3"
)

// Authorization routes
func AuthRoutes(app *fiber.App) {
	app.Post("/signup", controllers.SignUp)
	app.Post("/login", controllers.Login)
}