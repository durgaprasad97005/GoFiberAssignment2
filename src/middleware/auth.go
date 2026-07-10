package middleware

import (
	"log"
	"os"
	"strings"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"
)

var jwtSigningKey = []byte(os.Getenv("JWT_SIGNING_KEY"))

// middleware for authentication
func Auth(c fiber.Ctx) error {
	startTime := time.Now()

	// getting Authorization header
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false, 
			"message": "Unauthorized", 
			"error": "Missing Authorization header",
		})
	}

	// Get token from tokenString
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return jwtSigningKey, nil
	})

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false, 
			"message": "Error parsing token", 
			"error": err.Error(), 
		})
	}

	// Check token expiration 
	claims := token.Claims.(jwt.MapClaims)

	if claims["exp"].(float64) < float64(time.Now().Unix()) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false, 
			"message": "Unauthorized", 
			"error": "Token expired", 
		})
	}

	// Store required data in locals
	c.Locals("userId", claims["userId"])
	c.Locals("role", claims["role"])

	result := c.Next()

	requestTime := time.Since(startTime)
	log.Println(requestTime)

	return result
}