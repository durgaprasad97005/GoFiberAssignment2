package controllers

import (
	"context"
	"strconv"
	"time"

	"github.com/durgaprasad97005/GoFiberAssignment2/src/db"
	"github.com/durgaprasad97005/GoFiberAssignment2/src/models"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// Get all students with filter and pagination
func Find(c fiber.Ctx) error {
	// Collection
	collection := db.GetCollection("students")
	if collection == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Unable to find collection",
			"error":   "Internal server error",
		})
	}

	// Context
	ctx, cancel := context.WithTimeout(c.Context(), 5*time.Second)
	defer cancel()

	// Filter
	filter := bson.M{}
	branch := c.Query("branch")
	if branch != "" {
		filter["branch"] = bson.M{
			"$regex":   branch,
			"$options": "i",
		}
	}

	// findOptions - getting query parameters
	page, err := strconv.Atoi(c.Query("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(c.Query("limit", "10"))
	if err != nil || limit < 1 {
		limit = 10
	}

	skip := (page - 1) * limit

	findOptions := options.Find()
	findOptions.SetSort(bson.D{
		{Key: "name", Value: 1},
	})
	findOptions.SetSkip(int64(skip))
	findOptions.SetLimit(int64(limit))

	// Call db to get data
	cursor, err := collection.Find(ctx, filter, findOptions)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Internal server error",
			"error":   err.Error(),
		})
	}

	// Parse the result to students array
	var students []models.Student

	err = cursor.All(ctx, &students)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Internal server error",
			"error":   err.Error(),
		})
	}

	// Success response
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Successfully returned students data",
		"data":    students,
	})
}

// Get a student by Id
func Get(c fiber.Ctx) error {
	// Get student id from route/path parameters
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Student's object Id is required",
			"error":   "Object Id not found",
		})
	}

	// Parsing the string id to objId
	objId, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Invalid object Id",
			"error":   err.Error(),
		})
	}

	// Get collection
	collection := db.GetCollection("students")
	if collection == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Unable to find collection",
			"error":   "Internal server error",
		})
	}

	// Get the student from db
	ctx, cancel := context.WithTimeout(c.Context(), 5*time.Second)
	defer cancel()

	var student models.Student
	err = collection.FindOne(ctx, bson.M{"_id": objId}).Decode(&student)

	// If document not found
	if err == mongo.ErrNoDocuments {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "Student not found",
			"error":   err.Error(),
		})
	}

	// If some other error occurred
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Internal server error",
			"error":   err.Error(),
		})
	}

	// Success response
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Student found",
		"data":    student,
	})
}

// Create a student
func Create(c fiber.Ctx) error {
	// Get collection
	collection := db.GetCollection("students")
	if collection == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Unable to find collection",
			"error":   "Internal server error",
		})
	}

	// Parse the body
	var body models.Student
	if err := c.Bind().Body(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Error parsing body",
			"error":   err.Error(),
		})
	}

	// validation check
	var validate = validator.New()
	if err := validate.Struct(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Validation errors for invalid data",
			"error":   err.Error(),
		})
	}

	// Get createdBy user id
	userId, err := bson.ObjectIDFromHex(c.Locals("userId").(string))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false, 
			"message": "Error getting CreatedBy user id", 
			"error": err.Error(), 
		})
	}

	// converting body to DTO object
	studentDto := models.StudentDTO{
		Student: body,
		Audit: models.Audit{
			CreatedAt:      time.Now(),
			CreatedBy: userId,
			LastModifiedAt: time.Now(),
			LastModifiedBy: userId, 
		},
	}

	// Insert new student to database
	ctx, cancel := context.WithTimeout(c.Context(), 5*time.Second)
	defer cancel()

	result, err := collection.InsertOne(ctx, studentDto)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Internal server error",
			"error":   err.Error(),
		})
	}

	// Success response
	body.ID = result.InsertedID.(bson.ObjectID)
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"message": "Successfully inserted student",
		"data":    body,
	})
}

// Update a student by Id
func Update(c fiber.Ctx) error {
	// Get the id from request
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Student's object Id is required.",
			"error":   "Object Id not found",
		})
	}

	// Parse the id to get objId
	objId, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Invalid Id",
			"error":   err.Error(),
		})
	}

	// Get collection
	collection := db.GetCollection("students")
	if collection == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Collection not found",
			"error":   "Internal server error",
		})
	}

	// Parse the body
	var body models.Student
	if err := c.Bind().Body(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Error parsing body",
			"error":   err.Error(),
		})
	}

	// validation check
	validate := validator.New()
	if err := validate.Struct(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Unformatted data",
			"error":   err.Error(),
		})
	}

	// Get createdBy user id
	userId, err := bson.ObjectIDFromHex(c.Locals("userId").(string))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false, 
			"message": "Error getting CreatedBy user id", 
			"error": err.Error(), 
		})
	}

	// converting body to DTO object
	studentDto := models.StudentDTO{
		Student: body,
		Audit: models.Audit{
			LastModifiedAt: time.Now(), 
			LastModifiedBy: userId,
		},
	}

	// Updating document
	ctx, cancel := context.WithTimeout(c.Context(), 5*time.Second)
	defer cancel()

	result, err := collection.UpdateOne(ctx, bson.M{"_id": objId}, bson.M{"$set": studentDto})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Internal server error",
			"error":   err.Error(),
		})
	}

	if result.MatchedCount == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "No matching Id found",
			"error":   "Object Id not found",
		})
	}

	// Get the updated student
	var updatedStd models.Student
	err = collection.FindOne(ctx, bson.M{"_id": objId}).Decode(&updatedStd)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": true,
			"message": "Updated successfully, but unable return updated data",
			"error":   err.Error(),
		})
	}

	// Success response
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Updated student successfully",
		"data":    updatedStd,
	})
}

// Delete a student by Id
func Delete(c fiber.Ctx) error {
	// Get the id from request
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Id is required.",
		})
	}

	// Parse the id to get objId
	objId, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Invalid Id: " + err.Error(),
		})
	}

	// Get collection
	collection := db.GetCollection("students")
	if collection == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Internal server error",
		})
	}

	// Delete student
	ctx, cancel := context.WithTimeout(c.Context(), 5*time.Second)
	defer cancel()

	result, err := collection.DeleteOne(ctx, bson.M{"_id": objId})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Internal server error: " + err.Error(),
		})
	}

	if result.DeletedCount == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "Id not found",
		})
	}

	// Success response
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Successfully deleted student",
		"data":    result,
	})
}
