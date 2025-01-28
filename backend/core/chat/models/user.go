package models

import (
	"github.com/golang-jwt/jwt/v5"
	"time"
)

// Структура данных с информацией о пользователе
type User struct {
	Tdid      string    `json:"tdid" bson:"tdid"`             // Уникальный идентификатор пользователя
	Email     string    `json:"email" bson:"email"`           // Email пользователя
	Name      string    `json:"name" bson:"name"`             // Имя пользователя
	Login     string    `json:"login" bson:"login"`           // Логин пользователя
	Password  string    `json:"password" bson:"password"`     // Пароль пользователя
	CreatedAt time.Time `json:"created_at" bson:"created_at"` // Дата создания пользователя
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"` // Дата редактирования пользователя
	jwt.RegisteredClaims
}
