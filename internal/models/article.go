package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Article struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Title     string             `bson:"title" json:"title"`
	Content   string             `bson:"content,omitempty" json:"content"`
	Tags      []string           `bson:"tags,omitempty" json:"tags"`
	Author    string             `bson:"author" json:"author"`
	Status    string             `bson:"status" json:"status"`
	Published time.Time          `bson:"published,omitempty" json:"published"`
	Created   time.Time          `bson:"created" json:"created"`
	Updated   time.Time          `bson:"updated,omitempty" json:"updated"`
}
