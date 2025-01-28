package drivers

import (
	"context"
	"server/models"
)

type DbInterfase interface {
	Base
	User
	Contacts
	Profile
}

type Base interface {
	Connect() error
	Close()
}

type User interface {
	// GetUserEmailOrLogin Поиск пользователя по email или login
	GetUserEmailOrLogin(ctx context.Context, email string) (*models.Users, error)

	// GetUserTdidForEmailOrLogin Получение пользователя по email или login
	GetUserTdidForEmailOrLogin(ctx context.Context, email string) (*models.Users, error)

	// InserUser Добавление пользователя
	InserUser(ctx context.Context, user *models.Users) error

	// UpdateUser Изменение пользователя
	UpdateUser(ctx context.Context, user *models.Users) error

	// UpdatePassword Изменение пароля пользователя
	UpdatePassword(ctx context.Context, user *models.Users) error

	// GetUsersFromTdid Поиск пользователей по массиву tdid
	GetUsersFromTdid(ctx context.Context, tdids []string) ([]*models.User, error)

	// GetUsersByPhones Поиск пользователей по массиву телефонов
	GetUsersByPhones(ctx context.Context, phones []string) ([]models.User, error)
}

type Profile interface {
	// InsertProfile Добавление профайла пользователя
	InsertProfile(ctx context.Context, userId int) error

	// UpdateProfile Изменение профайла пользователя
	UpdateProfile(ctx context.Context, userId int, imagef3 string) error

	// GetProfile Получение профайла пользователя
	GetProfile(ctx context.Context, tdid string) (*models.Profile, error)
}

type Contacts interface {
	// GetContacts Получение контактов по tdid пользователя
	GetContacts(ctx context.Context, tdid string) ([]*models.Contacts, error)

	// AddContact Добавление контакта
	AddContact(ctx context.Context, contact *models.Contact) error

	// AddContacts Добавление контактов
	AddContacts(ctx context.Context, contacts []models.Contact) error

	// FindContact Поиск существования контакта по user_from и user_to
	FindContact(ctx context.Context, userFrom, UserTo string) error

	// FindContacts Поиск существования контактов с пользователемя
	FindContacts(ctx context.Context, userFrom string, UserTo []string) ([]models.Contact, error)
}
