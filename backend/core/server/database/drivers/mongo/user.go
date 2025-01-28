package mongo

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"log"
	"server/database/drivers"
	"server/models"
)

func (m Mongo) GetUserEmailOrLogin(ctx context.Context, email string) (*models.Users, error) {
	return nil, nil
}

func (m Mongo) GetUsersFromTdid(ctx context.Context, tdids []string) ([]*models.User, error) {
	return nil, nil
}

func (m Mongo) GetUserTdidForEmailOrLogin(ctx context.Context, email string) (*models.Users, error) {
	if email == "" || ctx.Err() != nil {
		return nil, drivers.ErrorID
	}

	getUser := &models.Users{}

	doc := bson.D{
		{Key: "email", Value: email},
	}

	err := m.db.Collection(CollectionUsers).FindOne(ctx, doc).Decode(getUser)
	if err != nil {
		log.Println("[Error] get user", err)
	}

	return getUser, nil
}

func (m Mongo) InserUser(ctx context.Context, user *models.Users) error {
	if user == nil {
		return drivers.ErrorID
	}

	doc := bson.D{
		{Key: "email", Value: user.Email},
		{Key: "name", Value: user.Name},
	}

	_, err := m.db.Collection(CollectionUsers).InsertOne(ctx, doc)
	if err != nil {
		log.Println("[Error] insert user", err)
	}

	return nil
}

func (m Mongo) UpdateUser(ctx context.Context, user *models.Users) error {
	return nil
}

func (m Mongo) UpdatePassword(ctx context.Context, user *models.Users) error {
	return nil
}

func (p *Mongo) GetUsersByPhones(ctx context.Context, phones []string) ([]models.User, error) {
	return nil, nil
}
