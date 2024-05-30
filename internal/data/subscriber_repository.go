package data

import (
	"context"
	"fmt"
	"time"

	"github.com/benk-techworld/www-backend/internal/db"
	"github.com/benk-techworld/www-backend/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

func (sr SubscriberRepo) Get(email string, filters models.Filters) ([]models.Subscriber, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	regexPattern := fmt.Sprintf(".*%s.*", email)

	filter := bson.M{
		"email": bson.M{"$regex": regexPattern, "$options": "i"},
	}

	findOptions := options.Find()
	findOptions.SetLimit(int64(filters.Limit()))
	findOptions.SetSkip(int64(filters.Offset()))
	findOptions.SetSort(bson.D{{Key: filters.SortField(), Value: filters.SortDirection()}})

	cursor, err := sr.storage.Client.Database(sr.storage.Name).Collection("subscribers").Find(ctx, filter, findOptions)
	if err != nil {
		return nil, err
	}

	defer cursor.Close(context.Background())

	subs := []models.Subscriber{}

	for cursor.Next(context.Background()) {

		var sub models.Subscriber

		err := cursor.Decode(&sub)
		if err != nil {
			return nil, err
		}

		subs = append(subs, sub)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return subs, nil

}

func (sr SubscriberRepo) Delete(id primitive.ObjectID) (*mongo.DeleteResult, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	filter := bson.M{"_id": id}

	res, err := sr.storage.Client.Database(sr.storage.Name).Collection("subscribers").DeleteOne(ctx, filter)

	return res, err

}
