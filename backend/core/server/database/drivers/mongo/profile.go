package mongo

import (
	"context"
	"server/models"
)

func (m Mongo) InsertProfile(ctx context.Context, userId int) error {
	return nil
}

func (m Mongo) UpdateProfile(ctx context.Context, userId int, imagef3 string) error {
	return nil
}

func (m Mongo) GetProfile(ctx context.Context, tdid string) (*models.Profile, error) {
	return nil, nil
}
