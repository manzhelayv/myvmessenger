package mongo

import (
	"chat/database/drivers"
	"chat/models"
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

// GetLastMessages Получение последних сообщений в чатах
func (m Mongo) GetLastMessages(ctx context.Context, userFromOrUserTo string) ([]*models.Messages, error) {
	if userFromOrUserTo == "" {
		return nil, drivers.ErrorID
	}

	match1 := bson.D{
		{
			Key:   "user_to",
			Value: userFromOrUserTo,
		},
	}

	match2 := bson.D{
		{
			Key:   "user_from",
			Value: userFromOrUserTo,
		},
	}

	pipeline2 := bson.A{
		bson.M{"$match": match2},
		bson.M{"$group": bson.M{
			"_id":        "$user_to",
			"user_from":  bson.M{"$last": "$user_from"},
			"user_to":    bson.M{"$last": "$user_to"},
			"message":    bson.M{"$last": "$message"},
			"file":       bson.M{"$last": "$file"},
			"created_at": bson.M{"$last": "$created_at"},
		}},
		bson.M{"$project": bson.M{
			"_id":        0,
			"user_from":  1,
			"user_to":    1,
			"message":    1,
			"file":       1,
			"created_at": 1,
		}},
	}

	pipeline := bson.A{
		bson.M{"$match": match1},
		bson.M{"$group": bson.M{
			"_id":        "$user_from",
			"user_from":  bson.M{"$last": "$user_to"},
			"user_to":    bson.M{"$last": "$user_from"},
			"message":    bson.M{"$last": "$message"},
			"file":       bson.M{"$last": "$file"},
			"created_at": bson.M{"$last": "$created_at"},
		}},
		bson.M{"$project": bson.M{
			"_id":        0,
			"user_from":  1,
			"user_to":    1,
			"message":    1,
			"file":       1,
			"created_at": 1,
		}},
		bson.M{"$unionWith": bson.M{"coll": "chat", "pipeline": pipeline2}},
		bson.M{"$sort": bson.M{"user_to": -1, "created_at": -1}},
		bson.M{"$group": bson.M{
			"_id":        "$user_to",
			"user_from":  bson.M{"$first": "$user_from"},
			"user_to":    bson.M{"$first": "$user_to"},
			"message":    bson.M{"$first": "$message"},
			"file":       bson.M{"$first": "$file"},
			"created_at": bson.M{"$first": "$created_at"},
		}},
	}

	res, err := m.db.Collection(CollectionChat).Aggregate(ctx, pipeline)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("[Error] get messages chat %v", err.Error()))
	}

	defer func() {
		_ = res.Close(ctx)
	}()

	if res.RemainingBatchLength() == 0 {
		return nil, drivers.ErrorID
	}

	result := make([]*models.Messages, 0, res.RemainingBatchLength())

	if err = res.All(ctx, &result); err != nil {
		return nil, err
	}

	return result, nil
}

// GetMessages Получение сообщений в чате
func (m Mongo) GetMessages(ctx context.Context, userFrom, userTo string) ([]*models.Messages, error) {
	if userFrom == "" || userTo == "" {
		return nil, drivers.ErrorID
	}

	doc1 := bson.M{
		"user_from": userFrom,
		"user_to":   userTo,
	}

	doc2 := bson.M{
		"user_from": userTo,
		"user_to":   userFrom,
	}

	pipeline := bson.M{
		"$or": []bson.M{
			doc1, doc2,
		},
	}

	findOptions := options.Find()
	findOptions.SetSort(bson.D{{"created_at", 1}})

	res, err := m.db.Collection(CollectionChat).Find(ctx, pipeline, findOptions)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("[Error] get messages chat %v", err.Error()))
	}

	defer func() {
		_ = res.Close(ctx)
	}()

	//if res.RemainingBatchLength() == 0 {
	//	return nil, models.ErrorFindMessages
	//}

	result := make([]*models.Messages, 0, res.RemainingBatchLength())

	if err = res.All(ctx, &result); err != nil {
		return nil, err
	}

	return result, nil
}

// SetMessage Добавление сообщения в чат
func (m Mongo) SetMessage(ctx context.Context, message *models.Messages) error {
	if message == nil {
		return drivers.ErrorID
	}

	doc := bson.D{
		{Key: "id", Value: message.UserFrom},
		{Key: "user_from", Value: message.UserFrom},
		{Key: "user_to", Value: message.UserTo},
		{Key: "message", Value: message.Message},
		{Key: "file", Value: message.File},
		{Key: "created_at", Value: message.CreatedAt},
		{Key: "updated_at", Value: message.UpdatedAt},
	}
	log.Println("doc", doc)
	_, err := m.db.Collection(CollectionChat).InsertOne(ctx, doc)
	if err != nil {
		log.Println("[Error] insert message", err)
	}

	return nil
}
