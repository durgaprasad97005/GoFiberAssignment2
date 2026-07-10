package models

import (
	"go.mongodb.org/mongo-driver/v2/bson"
)

// User struct
type User struct {
	ID           bson.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	UserName     string        `json:"userName" bson:"userName" validate:"required"`
	Email        string        `json:"email" bson:"email" validate:"required"`
	Password     string        `json:"password,omitempty" bson:"-" validate:"required"`
	PasswordHash string        `json:"-" bson:"passwordHash"`
	Role         string        `json:"role" bson:"role" validate:"required"`
}

type UserDTO struct {
	User  `bson:",inline"`
	Audit `bson:",inline"`
}
