package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Planet struct {
	ID          primitive.ObjectID `json:"_id" bson:"_id"`
	Name        string             `bson:"name, omitempty" json:"name"`
	Climate     string             `bson:"climate, omitempty" json:"climate"`
	Terrain     string             `bson:"terrain, omitempty" json:"terrain"`
	Appearances int                `bson:"appearances, omitempty" json:"appearances"`
	CreatedAt   time.Time          `bson:"createdAt, omitempty" json:"createdAt"`
}
