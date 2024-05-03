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

func (ar ArticleRepo) Get(title string, tags []string, filters models.Filters) ([]models.Article, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	regexPattern := fmt.Sprintf(".*%s.*", title)
	filter := bson.M{
		"title": bson.M{"$regex": regexPattern, "$options": "i"},
	}

	if len(tags) > 0 {
		filter["tags"] = bson.M{"$all": tags}
	}

	findOptions := options.Find()
	findOptions.SetLimit(int64(filters.Limit()))
	findOptions.SetSkip(int64(filters.Offset()))
	findOptions.SetSort(bson.D{{Key: filters.SortField(), Value: filters.SortDirection()}})

	cursor, err := ar.storage.Client.Database(ar.storage.Name).Collection("articles").Find(ctx, filter, findOptions)
	if err != nil {
		return nil, err
	}

	defer cursor.Close(context.Background())

	articles := []models.Article{}

	for cursor.Next(context.Background()) {

		var article models.Article

		err := cursor.Decode(&article)
		if err != nil {
			return nil, err
		}

		articles = append(articles, article)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return articles, nil
}
