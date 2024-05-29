package data

import (
	"context"
	"time"

	"github.com/benk-techworld/www-backend/internal/db"
	"github.com/benk-techworld/www-backend/internal/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SubscriberRepo struct {
	storage *db.DB
}

func (sr SubscriberRepo) Create(sub *models.Subscriber) error {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	res, err := sr.storage.Client.Database(sr.storage.Name).Collection("subscribers").InsertOne(ctx, sub)
	if err != nil {
		return err
	}

	sub.ID = res.InsertedID.(primitive.ObjectID)

	return nil
}
