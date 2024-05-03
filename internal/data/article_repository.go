package data

import (
	"context"
	"time"

	"github.com/benk-techworld/www-backend/internal/db"
	"github.com/benk-techworld/www-backend/internal/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
