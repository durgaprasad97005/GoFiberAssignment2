package src

import (
	"log"

	"github.com/durgaprasad97005/GoFiberAssignment2/src/db"
	"github.com/durgaprasad97005/GoFiberAssignment2/src/routes"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/joho/godotenv"
)

func SetupApp() *fiber.App {
	// Creating new fiber app
	app := fiber.New()

	// Adding global middleware
	app.Use(logger.New())
	app.Use(cors.New())

	// Loading environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("Error loading environment variables from .env file. Using system environment variables.")
	}

	// Connecting to database
	db.ConnectDB()

	// Declaring routes
	routes.StudentsRoutes(app)
	
	return app
}