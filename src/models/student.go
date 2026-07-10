package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

// struct used for request parsing
type Student struct {
	ID      bson.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name    string        `json:"name" bson:"name" validate:"required"`
	Email   string        `json:"email" bson:"email" validate:"omitempty,email"`
	Branch  string        `json:"branch" bson:"branch" validate:"required"`
	Age     int           `json:"age" bson:"age" validate:"required,min=0"`
	Phone   string        `json:"phone" bson:"phone" validate:"omitempty,len=10,numeric"`
	Address string        `json:"address" bson:"address"`
}

// Audit fields
type Audit struct {
	CreatedAt      time.Time     `bson:"createdAt,omitempty"`
	CreatedBy      bson.ObjectID `bson:"createdBy,omitempty"`
	LastModifiedAt time.Time     `bson:"lastModifiedAt"`
	LastModifiedBy bson.ObjectID `bson:"lastModifiedBy"`
}

// DTO object for student - represents database model
type StudentDTO struct {
	Student `bson:",inline"`
	Audit   `bson:",inline"`
}
