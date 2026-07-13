package controllers

import (
	"context"
	"strconv"
	"time"

	"github.com/durgaprasad97005/GoFiberAssignment2/src/db"
	"github.com/durgaprasad97005/GoFiberAssignment2/src/models"
	"github.com/durgaprasad97005/GoFiberAssignment2/src/utils"
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
		return utils.Error(
			c, 
			fiber.StatusInternalServerError, 
			"Unable to find collection", 
			"Internal server error", 
		)
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
		return utils.Error(
			c, 
			fiber.StatusInternalServerError, 
			"Internal server error", 
			err.Error(), 
		)
	}

	// Parse the result to students array
	var students []models.Student

	err = cursor.All(ctx, &students)
	if err != nil {
		return utils.Error(
			c, 
			fiber.StatusInternalServerError, 
			"Internal server error", 
			err.Error(), 
		)
	}

	// Success response
	return utils.Success(
		c, 
		fiber.StatusOK, 
		"Successfully returned students data", 
		students, 
	)
}

// Get a student by Id
func Get(c fiber.Ctx) error {
	// Get student id from route/path parameters
	id := c.Params("id")
	if id == "" {
		return utils.Error(
			c, 
			fiber.StatusBadRequest, 
			"Student's object Id is required", 
			"Object Id not found", 
		)
	}

	// Parsing the string id to objId
	objId, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return utils.Error(
			c, 
			fiber.StatusBadRequest, 
			"Invalid object Id", 
			err.Error(), 
		)
	}

	// Get collection
	collection := db.GetCollection("students")
	if collection == nil {
		return utils.Error(
			c, 
			fiber.StatusInternalServerError, 
			"Unable to find collection", 
			"Internal server error", 
		)
	}

	// Get the student from db
	ctx, cancel := context.WithTimeout(c.Context(), 5*time.Second)
	defer cancel()

	var student models.Student
	err = collection.FindOne(ctx, bson.M{"_id": objId}).Decode(&student)

	// If document not found
	if err == mongo.ErrNoDocuments {
		return utils.Error(
			c, 
			fiber.StatusNotFound, 
			"Student not found", 
			err.Error(), 
		)
	}

	// If some other error occurred
	if err != nil {
		return utils.Success(
			c, 
			fiber.StatusInternalServerError, 
			"Internal server error", 
			err.Error(), 
		)
	}

	// Success response
	return utils.Success(
		c, 
		fiber.StatusOK, 
		"Student found", 
		student, 
	)
}

// Create a student
func Create(c fiber.Ctx) error {
	// Get collection
	collection := db.GetCollection("students")
	if collection == nil {
		return utils.Error(
			c, 
			fiber.StatusInternalServerError, 
			"Unable to find collection", 
			"Internal server error", 
		)
	}

	// Parse the body
	var body models.Student
	if err := c.Bind().Body(&body); err != nil {
		return utils.Error(
			c, 
			fiber.StatusBadRequest, 
			"Error parsing body", 
			err.Error(), 
		)
	}

	// validation check
	var validate = validator.New()
	if err := validate.Struct(&body); err != nil {
		return utils.Error(
			c, 
			fiber.StatusBadRequest, 
			"Validation error for invalid data", 
			err.Error(), 
		)
	}

	// Get createdBy user id
	userId, err := bson.ObjectIDFromHex(c.Locals("userId").(string))
	if err != nil {
		return utils.Error(
			c, 
			fiber.StatusInternalServerError, 
			"Error getting CreatedBy user id", 
			err.Error(), 
		)
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
		return utils.Error(
			c, 
			fiber.StatusInternalServerError, 
			"Internal server error", 
			err.Error(), 
		)
	}

	// Success response
	body.ID = result.InsertedID.(bson.ObjectID)
	return utils.Success(
		c, 
		fiber.StatusCreated, 
		"Successfully inserted student", 
		body, 
	)
}

// Update a student by Id
func Update(c fiber.Ctx) error {
	// Get the id from request
	id := c.Params("id")
	if id == "" {
		return utils.Error(
			c, 
			fiber.StatusBadRequest, 
			"Student's object Id is required", 
			"Object Id not found", 
		)
	}

	// Parse the id to get objId
	objId, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return utils.Error(
			c, 
			fiber.StatusInternalServerError, 
			"Invalid Id", 
			err.Error(), 
		)
	}

	// Get collection
	collection := db.GetCollection("students")
	if collection == nil {
		return utils.Error(
			c, 
			fiber.StatusInternalServerError, 
			"Collection not found", 
			"Internal server error", 
		)
	}

	// Parse the body
	var body models.Student
	if err := c.Bind().Body(&body); err != nil {
		return utils.Error(
			c, 
			fiber.StatusBadRequest, 
			"Error parsing body", 
			err.Error(), 
		)
	}

	// validation check
	validate := validator.New()
	if err := validate.Struct(&body); err != nil {
		return utils.Error(
			c, 
			fiber.StatusBadRequest, 
			"Unformatted data", 
			err.Error(), 
		)
	}

	// Get createdBy user id
	userId, err := bson.ObjectIDFromHex(c.Locals("userId").(string))
	if err != nil {
		return utils.Error(
			c, 
			fiber.StatusInternalServerError, 
			"Error getting CreatedBy user id", 
			err.Error(), 
		)
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
		return utils.Error(
			c, 
			fiber.StatusInternalServerError, 
			"Internal server error", 
			err.Error(), 
		)
	}

	if result.MatchedCount == 0 {
		return utils.Error(
			c, 
			fiber.StatusNotFound, 
			"No matching Id found", 
			"Object Id not found", 
		)
	}

	// Get the updated student
	var updatedStd models.Student
	err = collection.FindOne(ctx, bson.M{"_id": objId}).Decode(&updatedStd)
	if err != nil {
		return utils.Success(
			c, 
			fiber.StatusOK, 
			"Updated successfully, but unable to return updated data", 
			nil, 
		)
	}

	// Success response
	return utils.Success(
		c, 
		fiber.StatusOK, 
		"Updated student successfully", 
		updatedStd, 
	)
}

// Delete a student by Id
func Delete(c fiber.Ctx) error {
	// Get the id from request
	id := c.Params("id")
	if id == "" {
		return utils.Error(
			c, 
			fiber.StatusBadRequest, 
			"Id is required", 
			"Cannot find Id", 
		)
	}

	// Parse the id to get objId
	objId, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return utils.Error(
			c, 
			fiber.StatusInternalServerError, 
			"Invalid Id", 
			err.Error(), 
		)
	}

	// Get collection
	collection := db.GetCollection("students")
	if collection == nil {
		return utils.Error(
			c, 
			fiber.StatusInternalServerError, 
			"Internal server error", 
			"Cannot find the collection in database", 
		)
	}

	// Delete student
	ctx, cancel := context.WithTimeout(c.Context(), 5*time.Second)
	defer cancel()

	result, err := collection.DeleteOne(ctx, bson.M{"_id": objId})
	if err != nil {
		return utils.Error(
			c, 
			fiber.StatusInternalServerError, 
			"Internal server error", 
			err.Error(), 
		)
	}

	if result.DeletedCount == 0 {
		return utils.Error(
			c, 
			fiber.StatusNotFound, 
			"No documents found to delete", 
			"Id not found", 
		)
	}

	// Success response
	return utils.Success(
		c, 
		fiber.StatusOK, 
		"Successfully deleted student", 
		result, 
	)
}
