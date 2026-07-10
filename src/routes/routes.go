package routes

import "github.com/gofiber/fiber/v3"

// Setting up routes
func SetupRoutes(app *fiber.App) {
	AuthRoutes(app)
	StudentRoutes(app)
}