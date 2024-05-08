package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Subscriber struct {
	ID      primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Email   string             `bson:"email" json:"email"`
	Created time.Time          `bson:"created" json:"created"`
}
