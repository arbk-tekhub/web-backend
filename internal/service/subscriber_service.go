package service

import (
	"errors"
	"time"

	"github.com/benk-techworld/www-backend/internal/models"
	"github.com/benk-techworld/www-backend/internal/validator"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CreateSubscriberInput struct {
	Email            string            `json:"email"`
	ValidationErrors map[string]string `json:"-"`
}

type FetchSubscribersInput struct {
	Email string
	models.Filters
	ValidationErrors map[string]string
}

func (svc Service) CreateSubscriber(input *CreateSubscriberInput) (*models.Subscriber, error) {

	v := validator.New()
	input.ValidationErrors = v.Errors

	v.Check(validator.Matchs(input.Email, *validator.EmailRX), "email", "must be a valid email address")

	if v.HasErrors() {
		return nil, ErrFailedValidation
	}

	sub := &models.Subscriber{
		Email:   input.Email,
		Created: time.Now(),
	}

	err := svc.Repo.Subscriber.Create(sub)
	if err != nil {
		return nil, err
	}

	return sub, nil
}

func (svc Service) GetSubscribers(input *FetchSubscribersInput) ([]models.Subscriber, error) {
	v := validator.New()
	input.ValidationErrors = v.Errors

	input.SortSafeList = []string{"email", "created", "-email", "-created"}
	if models.ValidateFilters(v, input.Filters); v.HasErrors() {
		return nil, ErrFailedValidation
	}

	subs, err := svc.Repo.Subscriber.Get(input.Email, input.Filters)
	if err != nil {
		return nil, err
	}

	return subs, nil
}

func (svc Service) DeleteSubscriber(id string) error {

	docID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		if errors.Is(err, primitive.ErrInvalidHex) {
			return ErrNotFound
		}
		return err
	}

	res, err := svc.Repo.Subscriber.Delete(docID)
	if err != nil {
		return err
	}

	if res.DeletedCount == 0 {
		return ErrNotFound
	}

	return nil
}
