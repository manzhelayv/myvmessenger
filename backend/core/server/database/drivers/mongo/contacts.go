package mongo

import (
	"context"
	"server/models"
)

func (m Mongo) GetContacts(ctx context.Context, tdid string) ([]*models.Contacts, error) {
	return nil, nil
}

func (m Mongo) AddContact(ctx context.Context, contact *models.Contact) error {
	return nil
}

func (m Mongo) AddContacts(ctx context.Context, contacts []models.Contact) error {
	return nil
}

func (m Mongo) FindContact(ctx context.Context, userFrom, UserTo string) error {
	return nil
}

func (p *Mongo) FindContacts(ctx context.Context, userFrom string, UserTo []string) ([]models.Contact, error) {
	return nil, nil
}
