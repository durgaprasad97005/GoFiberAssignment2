package routes

import (
	"github.com/durgaprasad97005/GoFiberAssignment2/src/controllers"
	"github.com/gofiber/fiber/v3"
)

func StudentsRoutes(app *fiber.App) {
	// Grouping the route for students
	students := app.Group("/students")

	// Creating routes
	students.Get("/", controllers.GetStudents)
	students.Get("/:id", controllers.GetStudentById)
	students.Post("/", controllers.CreateStudent)
	students.Put("/:id", controllers.UpdateStudentById)
	students.Delete("/:id", controllers.DeleteStudentById)
}