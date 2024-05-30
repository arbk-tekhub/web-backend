package data

import (
	"github.com/benk-techworld/www-backend/internal/db"
	"github.com/benk-techworld/www-backend/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Repository struct {
	Article interface {
		Create(article *models.Article) error
		GetByID(id primitive.ObjectID) (*models.Article, error)
		Delete(id primitive.ObjectID) (*mongo.DeleteResult, error)
		Get(title string, tags []string, filters models.Filters) ([]models.Article, error)
		Update(id primitive.ObjectID, updateDoc bson.M) *mongo.SingleResult
	}
	Subscriber interface {
		Create(sub *models.Subscriber) error
		Get(email string, filters models.Filters) ([]models.Subscriber, error)
		Delete(id primitive.ObjectID) (*mongo.DeleteResult, error)
	}
}

func NewRepo(db *db.DB) *Repository {
	return &Repository{
		Article:    ArticleRepo{db},
		Subscriber: SubscriberRepo{db},
	}
}
