package data

import (
	"context"
	"time"

	"github.com/benk-techworld/www-backend/internal/db"
	"github.com/benk-techworld/www-backend/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ArticleRepo struct {
	storage *db.DB
}

func (ar ArticleRepo) Create(article *models.Article) error {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	res, err := ar.storage.Client.Database(ar.storage.Name).Collection("articles").InsertOne(ctx, article)
	if err != nil {
		return err
	}
	article.ID = res.InsertedID.(primitive.ObjectID)

	return nil
}

func (ar ArticleRepo) GetByID(id primitive.ObjectID) (*models.Article, error) {

	var article models.Article

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	filter := bson.M{"_id": id}

	err := ar.storage.Client.Database(ar.storage.Name).Collection("articles").FindOne(ctx, filter).Decode(&article)
	if err != nil {
		return nil, err
	}

	return &article, nil
}

func (ar ArticleRepo) Delete(id primitive.ObjectID) (*mongo.DeleteResult, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	filter := bson.M{"_id": id}

	res, err := ar.storage.Client.Database(ar.storage.Name).Collection("articles").DeleteOne(ctx, filter)

	return res, err

}
