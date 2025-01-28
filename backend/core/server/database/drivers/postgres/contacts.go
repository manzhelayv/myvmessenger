package postgres

import (
	"context"
	"errors"
	"server/database/drivers"
	"server/models"
)

// GetContacts Получение контактов по tdid пользователя
func (p *Postgres) GetContacts(ctx context.Context, tdid string) ([]*models.Contacts, error) {
	if tdid == "" {
		return nil, drivers.ErrorID
	}

	var contacts []*models.Contacts

	err := p.client.ModelContext(ctx, &contacts).
		Column("user_from", "user_to", "user_to_name", "users.email", "profiles.image").
		Where("user_from = ?", tdid).
		Join("JOIN users ON user_to = users.tdid").
		Join("JOIN profiles ON profiles.user_id = users.id").
		Select()
	if err != nil {
		return nil, err
	}

	return contacts, nil
}

// AddContact Добавление контакта
func (p *Postgres) AddContact(ctx context.Context, contact *models.Contact) error {
	if contact == nil {
		return drivers.ErrorID
	}

	_, err := p.client.ModelContext(ctx, contact).Insert()
	if err != nil {
		return err
	}

	return nil
}

// AddContacts Добавление контактов
func (p *Postgres) AddContacts(ctx context.Context, contacts []models.Contact) error {
	if contacts == nil {
		return drivers.ErrorID
	}

	_, err := p.client.ModelContext(ctx, &contacts).Insert()
	if err != nil {
		return err
	}

	return nil
}

// FindContact Поиск существования контакта по user_from и user_to
func (p *Postgres) FindContact(ctx context.Context, userFrom, UserTo string) error {
	if userFrom == "" || UserTo == "" {
		return drivers.ErrorID
	}

	var contacts models.Contact

	err := p.client.ModelContext(ctx, &contacts).
		Where("user_from = ?", userFrom).
		Where("user_to = ?", UserTo).
		Select()
	if err != nil && err.Error() == "pg: no rows in result set" {
		return nil
	}

	if contacts.UserFrom == userFrom && contacts.UserTo == UserTo {
		return errors.New("Такой контакт уже существует")
	}

	return nil
}

// FindContacts Поиск существования контактов с пользователемя
func (p *Postgres) FindContacts(ctx context.Context, userFrom string, UserTo []string) ([]models.Contact, error) {
	if userFrom == "" || UserTo == nil {
		return nil, drivers.ErrorID
	}

	var contacts []models.Contact

	err := p.client.ModelContext(ctx, &contacts).
		Column("user_from", "user_to").
		Where("user_from = ?", userFrom).
		//WhereIn("user_from in (?)", UserTo).
		//Where("user_to = ?", userFrom).
		//WhereOr("user_from = ?", userFrom).
		WhereIn("user_to in (?)", UserTo).
		Select()
	if err != nil {
		return nil, err
	}

	return contacts, nil
}
