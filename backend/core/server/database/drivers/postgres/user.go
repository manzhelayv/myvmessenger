package postgres

import (
	"context"
	"errors"
	"github.com/go-pg/pg/v10"
	"server/database/drivers"
	"server/models"
)

// GetUsersFromTdid Поиск пользователей по массиву tdid
func (p *Postgres) GetUsersFromTdid(ctx context.Context, tdids []string) ([]*models.User, error) {
	if len(tdids) == 0 {
		return nil, drivers.ErrorID
	}

	users := []*models.User{}

	err := p.client.ModelContext(ctx, &users).
		Column("id", "email", "tdid", "name", "login", "profiles.image").
		Where("tdid in (?)", pg.In(tdids)).
		Join("JOIN profiles ON profiles.user_id = id").
		Select()
	if err != nil {
		return nil, err
	}

	return users, nil
}

// GetUserEmailOrLogin Поиск пользователя по email или login
func (p *Postgres) GetUserEmailOrLogin(ctx context.Context, emailOrLogin string) (*models.Users, error) {
	if emailOrLogin == "" {
		return nil, drivers.ErrorID
	}

	user := &models.Users{}

	err := p.client.ModelContext(ctx, user).
		Where("email = ?", emailOrLogin).
		WhereOr("login = ?", emailOrLogin).
		Select()
	if err != nil {
		return nil, err
	}

	return user, nil
}

// GetUserTdidForEmailOrLogin Получение пользователя по email или login
func (p *Postgres) GetUserTdidForEmailOrLogin(ctx context.Context, emailOrLogin string) (*models.Users, error) {
	if emailOrLogin == "" {
		return nil, drivers.ErrorID
	}

	user := &models.Users{}

	err := p.client.ModelContext(ctx, user).
		Where("email = ?", emailOrLogin).
		WhereOr("login = ?", emailOrLogin).
		Select()
	if err != nil {
		return nil, err
	}

	return user, nil
}

// InserUser Добавление пользователя
func (p *Postgres) InserUser(ctx context.Context, user *models.Users) error {
	if user == nil {
		return drivers.ErrorID
	}

	err := p.SearchTdidNotExists(ctx, user.Tdid)
	if err != nil {
		return err
	}

	err = user.ValidEmail()
	if err != nil {
		return errors.New("Неправильный формат email адреса")
	}

	_, err = p.client.ModelContext(ctx, user).Insert()
	if err != nil {
		return err
	}

	err = p.InsertProfile(ctx, user.Id)
	if err != nil {
		return err
	}

	return nil
}

func (p *Postgres) SearchTdidNotExists(ctx context.Context, tdid string) error {
	if tdid == "" {
		return drivers.ErrorID
	}

	user := &models.Users{}

	err := p.client.ModelContext(ctx, user).Where("tdid = ?", tdid).Select()
	if err != nil {
		return nil
	}

	return errors.New("Пользователь с таким идентификатором уже существует, пожалуйста, повторите попытку еще раз")
}

// UpdateUser Изменение пользователя
func (p *Postgres) UpdateUser(ctx context.Context, user *models.Users) error {
	if user == nil {
		return drivers.ErrorID
	}

	userSet := "email = ?email,login = ?login,name = ?name,phone = ?phone,updated_at = ?updated_at"

	_, err := p.client.ModelContext(ctx, user).Set(userSet).Where("tdid = ?", user.Tdid).Update()
	if err != nil {
		return err
	}

	return nil
}

// UpdatePassword Изменение пароля пользователя
func (p *Postgres) UpdatePassword(ctx context.Context, user *models.Users) error {
	if user == nil {
		return drivers.ErrorID
	}

	userSet := "password = ?password"

	_, err := p.client.ModelContext(ctx, user).Set(userSet).Where("tdid = ?", user.Tdid).Update()
	if err != nil {
		return err
	}

	return nil
}

// GetUsersByPhones Поиск пользователей по массиву телефонов
func (p *Postgres) GetUsersByPhones(ctx context.Context, phones []string) ([]models.User, error) {
	if len(phones) == 0 {
		return nil, drivers.ErrorID
	}

	users := []models.User{}

	err := p.client.ModelContext(ctx, &users).
		Column("tdid", "phone").
		Where("phone in (?)", pg.In(phones)).
		Select()
	if err != nil {
		return nil, err
	}

	return users, nil
}
