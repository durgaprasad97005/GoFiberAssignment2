package middleware

import (
	"context"
	"time"

	"github.com/durgaprasad97005/GoFiberAssignment2/src/db"
	"github.com/durgaprasad97005/GoFiberAssignment2/src/models"
	"github.com/durgaprasad97005/GoFiberAssignment2/src/utils"
	"github.com/gofiber/fiber/v3"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

// Admin middleware
func Admin(c fiber.Ctx) error {
	userId := c.Locals("userId").(string)
	objId, err := bson.ObjectIDFromHex(userId)
	if err != nil {
		return utils.Error(
			c, 
			fiber.StatusInternalServerError, 
			"Error parsing document Id", 
			err.Error(),
		)
	}

	ctx, cancel := context.WithTimeout(c.Context(), 5 * time.Second)
	defer cancel()

	collection := db.GetCollection("users")
	if collection == nil {
		return utils.Error(
			c, 
			fiber.StatusInternalServerError, 
			"Error getting collection from db", 
			"Cannot access database collection", 
		)
	}

	var user models.User
	err = collection.FindOne(ctx, bson.M{"_id": objId}).Decode(&user)

	if err == mongo.ErrNoDocuments {
		return utils.Error(
			c, 
			fiber.StatusNotFound, 
			"Document not found", 
			err.Error(), 
		)
	}

	if err != nil {
		return utils.Error(
			c, 
			fiber.StatusInternalServerError, 
			"Error finding user document", 
			err.Error(), 
		)
	}

	// check for the role
	if user.Role != "admin" {
		return utils.Error(
			c, 
			fiber.StatusUnauthorized, 
			"Unauthorized", 
			"Found invalid role. This route requires admin role", 
		)
	}

	return c.Next()
}