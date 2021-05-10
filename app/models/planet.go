package models

import (
	"time"

	"github.com/asaskevich/govalidator"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Planet struct {
	ID          primitive.ObjectID `json:"_id" valid:"-" bson:"_id, omitempty"`
	Name        string             `bson:"name, omitempty" valid:"notnull" json:"name"`
	Climate     string             `bson:"climate, omitempty" valid:"notnull" json:"climate"`
	Terrain     string             `bson:"terrain, omitempty" valid:"notnull" json:"terrain"`
	Appearances int                `bson:"appearances, omitempty" valid:"-" json:"appearances"`
	CreatedAt   time.Time          `bson:"createdAt, omitempty" valid:"-" json:"createdAt"`
}

func init() {
	govalidator.SetFieldsRequiredByDefault(true)
}

func (planet *Planet) Validate() error {
	_, err := govalidator.ValidateStruct(planet)

	if err != nil {
		return err
	}

	return nil
}
