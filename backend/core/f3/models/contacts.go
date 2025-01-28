package models

import (
	"errors"
	"time"
)

var ErrJsonRequestContact = errors.New("Введите логин или email пользователя")

// Contacts Структура таблицы контактов, с добавлением email пользователя с таблицы пользователей,
// для получения из базы данных, таблицы контактов
type Contacts struct {
	UserFrom   string    `json:"user_from" bson:"user_from"`       // Уникальный идентификатор текущего пользователя
	UserTo     string    `json:"user_to" bson:"user_to"`           // Уникальный идентификатор пользователя добавленного в контакты
	UserToName string    `json:"user_to_name" bson:"user_to_name"` // Имя пользователя добавленного в контакты
	Email      string    `json:"email" bson:"email"`               // Email пользователя, которому отпраляются сообщения
	CreatedAt  time.Time `json:"created_at" bson:"created_at"`     // Дата создания контакта
	UpdatedAt  time.Time `json:"updated_at" bson:"updated_at"`     // Дата редактирования контакта
}

// Contact Структура таблицы контактов, для добавления в базу данных, в таблицу контактов
type Contact struct {
	UserFrom   string    `bson:"user_from"`    // Уникальный идентификатор текущего пользователя
	UserTo     string    `bson:"user_to"`      // Уникальный идентификатор пользователя добавляемого в контакты
	UserToName string    `bson:"user_to_name"` // Имя пользователя добавляемого в контакты
	CreatedAt  time.Time `bson:"created_at"`   // Дата создания контакта
	UpdatedAt  time.Time `bson:"updated_at"`   // Дата редактирования контакта
}

// RequestContact Структура HTTP-запроса для добавления контакта
type RequestContact struct {
	EmailOrLogin string `json:"email_or_login"` // Email или логин пользователя с HTTP-запроса
}
