package postgres

import (
	"context"
	"server/database/drivers"
	"server/models"
)

// InsertProfile Добавление профайла пользователя
func (p *Postgres) InsertProfile(ctx context.Context, userId int) error {
	if userId == 0 {
		return drivers.ErrorID
	}

	profile := &models.Profile{
		UserId: userId,
	}

	_, err := p.client.ModelContext(ctx, profile).Insert()
	if err != nil {
		return err
	}

	return nil
}

// UpdateProfile Изменение профайла пользователя
func (p *Postgres) UpdateProfile(ctx context.Context, userId int, imagef3 string) error {
	if userId == 0 {
		return drivers.ErrorID
	}

	profile := &models.Profile{
		UserId: userId,
		Image:  imagef3,
	}

	_, err := p.client.ModelContext(ctx, profile).
		Where("user_id = ?", userId).
		Update()
	if err != nil {
		return err
	}

	return nil
}

func (p Postgres) GetProfile(ctx context.Context, tdid string) (*models.Profile, error) {
	if tdid == "" {
		return nil, drivers.ErrorID
	}

	var profile models.Profile

	err := p.client.ModelContext(ctx, &profile).
		//Column("image").
		Where("users.tdid = ?", tdid).
		Join("JOIN users ON user_id = users.id").
		Select()
	if err != nil {
		return nil, err
	}

	return &profile, nil
}
