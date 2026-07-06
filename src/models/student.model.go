package models

import "go.mongodb.org/mongo-driver/v2/bson"

// struct used for both creation and representing database model
type Student struct {
	ID      bson.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name    string        `json:"name" bson:"name" validate:"required"`
	Email   string        `json:"email" bson:"email" validate:"omitempty,email"`
	Branch  string        `json:"branch" bson:"branch" validate:"required"`
	Age     int           `json:"age" bson:"age" validate:"required,min=0"`
	Phone   string        `json:"phone" bson:"phone" validate:"omitempty,len=10,numeric"`
	Address string        `json:"address" bson:"address"`
}

// struct used for update
type UpdateStudentDTO struct {
	Name    *string `json:"name" bson:"name,omitempty"`
	Email   *string `json:"email" bson:"email,omitempty"`
	Branch  *string `json:"branch" bson:"branch,omitempty"`
	Age     *int    `json:"age" bson:"age,omitempty" validate:"omitempty,min=0"`
	Phone   *string `json:"phone" bson:"phone,omitempty"`
	Address *string `json:"address" bson:"address,omitempty"`
}
