package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Player struct {
	Id       primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Name     string             `json:"name,omitempty" validate:"required"`
	Region   string             `json:"region,omitempty" validate:"required"`
	Position string             `json:"position,omitempty" validate:"required"`
}
