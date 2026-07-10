package routes

import (
	"github.com/durgaprasad97005/GoFiberAssignment2/src/controllers"
	"github.com/durgaprasad97005/GoFiberAssignment2/src/middleware"
	"github.com/gofiber/fiber/v3"
)

// Student routes
func StudentRoutes(app *fiber.App) {
	students := app.Group("/students", middleware.Auth, middleware.Admin)

	students.Get("/", controllers.Find)
	students.Get("/:id", controllers.Get)
	students.Post("/", controllers.Create)
	students.Put("/:id", controllers.Update)
	students.Delete("/:id", controllers.Delete)
}