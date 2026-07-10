package middleware

import "github.com/gofiber/fiber/v3"

// Admin middleware
func Admin(c fiber.Ctx) error {
	// check for the role
	if c.Locals("role").(string) != "admin" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false, 
			"message": "Unauthorized", 
			"error": "Found invalid role. This route requires admin role", 
		})
	}

	return c.Next()
}