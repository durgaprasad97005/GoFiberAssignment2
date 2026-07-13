package middleware

import (
	"log"
	"os"
	"strings"
	"time"

	"github.com/durgaprasad97005/GoFiberAssignment2/src/utils"
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
		return utils.Error(
			c, 
			fiber.StatusUnauthorized, 
			"Unauthorized", 
			"Missing Authorization header", 
		)
	}

	// Get token from tokenString
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return jwtSigningKey, nil
	})

	if err != nil {
		return utils.Error(
			c, 
			fiber.StatusInternalServerError, 
			"Error parsing token", 
			err.Error(), 
		)
	}

	// Check token expiration 
	claims := token.Claims.(jwt.MapClaims)

	if claims["exp"].(float64) < float64(time.Now().Unix()) {
		return utils.Error(
			c, 
			fiber.StatusUnauthorized, 
			"Unauthorized", 
			"Token expired", 
		)
	}

	// Store required data in locals
	c.Locals("userId", claims["userId"])

	result := c.Next()

	requestTime := time.Since(startTime)
	log.Println(requestTime)

	return result
}