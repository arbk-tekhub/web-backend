package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Article struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	Title     string             `bson:"title"`
	Content   string             `bson:"content,omitempty"`
	Tags      []string           `bson:"tags,omitempty"`
	Author    string             `bson:"author"`
	Status    string             `bson:"status"`
	Published time.Time          `bson:"published,omitempty"`
	Created   time.Time          `bson:"created"`
	Updated   time.Time          `bson:"updated,omitempty"`
}
