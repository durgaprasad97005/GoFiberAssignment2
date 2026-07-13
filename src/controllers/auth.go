package controllers

import (
	"context"
	"os"
	"time"

	"github.com/durgaprasad97005/GoFiberAssignment2/src/db"
	"github.com/durgaprasad97005/GoFiberAssignment2/src/models"
	"github.com/durgaprasad97005/GoFiberAssignment2/src/utils"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/v2/bson"
	"golang.org/x/crypto/bcrypt"
)

var jwtSigningKey = []byte(os.Getenv("JWT_SIGNING_KEY"))

// function to get jwt token string
func getJwtTokenString(user models.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId": user.ID.Hex(),
		"exp":    time.Now().Add(3 * time.Hour).Unix(),
	})

	return token.SignedString(jwtSigningKey)
}

// Sign up controller
func SignUp(c fiber.Ctx) error {
	// Parsing request body
	var body models.User
	if err := c.Bind().Body(&body); err != nil {
		return utils.Error(
			c, 
			fiber.StatusBadRequest, 
			"Error parsing body", 
			err.Error(), 
		)
	}

	// Validate the request body
	validate := validator.New()
	if err := validate.Struct(&body); err != nil {
		return utils.Error(
			c, 
			fiber.StatusBadRequest, 
			"Data validation error", 
			err.Error(), 
		)
	}

	// Get the collection
	collection := db.GetCollection("users")
	if collection == nil {
		return utils.Error(
			c, 
			fiber.StatusInternalServerError, 
			"Unable to find collection", 
			"Internal server error", 
		)
	}

	// Context creation
	ctx, cancel := context.WithTimeout(c.Context(), 5*time.Second)
	defer cancel()

	// Check whether user exists already
	var existingUser models.User
	err := collection.FindOne(ctx, bson.M{"email": body.Email}).Decode(&existingUser)
	if err == nil {
		return utils.Error(
			c, 
			fiber.StatusConflict, 
			"There exists another user with the given email", 
			"User creation failed", 
		)
	}

	// Create password hash
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(body.Password), 12)
	if err != nil {
		return utils.Error(
			c, 
			fiber.StatusInternalServerError, 
			"Error while hashing password", 
			err.Error(), 
		)
	}

	body.PasswordHash = string(hashedPassword)

	// create UserDTO object
	userDto := models.UserDTO{
		User: body,
		Audit: models.Audit{
			CreatedAt:      time.Now(),
			LastModifiedAt: time.Now(), // Might need to add CreatedBy and LastModifiedBy fields
		},
	}

	// Insert user
	result, err := collection.InsertOne(ctx, userDto)
	if err != nil {
		return utils.Error(
			c, 
			fiber.StatusInternalServerError, 
			"Error while trying to signup", 
			err.Error(), 
		)
	}

	// Updating fields for response
	body.ID = result.InsertedID.(bson.ObjectID)
	body.Password = ""

	// Creating jwt token string
	jwtTokenString, err := getJwtTokenString(body)
	if err != nil {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"success": true,
			"message": "Signup successful",
			"data":    body,
			"jwtToken": "", 
		})
	}

	// Success response
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Signup successful",
		"data":    body,
		"jwtToken": jwtTokenString, 
	})
}

// Login controller
func Login(c fiber.Ctx) error {
	// Parse body of the reques
	var body models.User
	if err := c.Bind().Body(&body); err != nil {
		return utils.Error(
			c, 
			fiber.StatusBadRequest, 
			"Error while parsing body", 
			err.Error(), 
		)
	}

	// Check whether required data exists or not
	if body.Email == "" || body.Password == "" {
		return utils.Error(
			c, 
			fiber.StatusBadRequest, 
			"Email and Password fields are required", 
			"Required fields are missing", 
		)
	}

	// Get collection
	collection := db.GetCollection("users")
	if collection == nil {
		return utils.Error(
			c, 
			fiber.StatusInternalServerError, 
			"Unabe to get the collection", 
			"Internal server error", 
		)
	}

	// Create context variable
	ctx, cancel := context.WithTimeout(c.Context(), 5 * time.Second)
	defer cancel()

	// Check whether user record exists
	var existingUser models.User
	if err := collection.FindOne(ctx, bson.M{"email": body.Email}).Decode(&existingUser); err != nil {
		return utils.Error(
			c, 
			fiber.StatusNotFound, 
			"There exists no user with the given email", 
			err.Error(), 
		)
	}

	// Compare password hash
	if err := bcrypt.CompareHashAndPassword([]byte(existingUser.PasswordHash), []byte(body.Password)); err != nil {
		return utils.Error(
			c, 
			fiber.StatusUnauthorized, 
			"Wrong password", 
			err.Error(), 
		)
	}

	// Generate jwt token string
	jwtTokenString, err := getJwtTokenString(existingUser)
	if err != nil {
		return utils.Error(
			c, 
			fiber.StatusInternalServerError, 
			"Unable to generate token string", 
			err.Error(), 
		)
	}

	// Success response
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true, 
		"message": "Logged in successfully", 
		"jwtToken": jwtTokenString, 
	})
}
