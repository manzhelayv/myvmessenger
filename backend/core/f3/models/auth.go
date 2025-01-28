package models

import (
	"errors"
	"github.com/badoux/checkmail"
)

// Секретный ключ jwt - Todo перенетси в конфиги
var JwtSecretKey = []byte("secret-key")

const (
	ClientType    = "Client-Type"
	ClientVersion = "Client-Version"
	DeviceID      = "Device-id"
	Channel       = "Channel"
	Affiliation   = "Affiliation"
	Uuid          = "Uuid"
)

func AllowedHeaders() []string {
	return []string{
		ClientType,
		ClientVersion,
		DeviceID,
		Channel,
		Affiliation,
		Uuid,
	}
}

// RegisterRequest Структура для регистрации пользователя
type RegisterRequest struct {
	Email    string `json:"email"`    // Email пользователя
	Login    string `json:"login"`    // Логин пользователя
	Name     string `json:"name"`     // Имя пользователя
	Password string `json:"password"` // Пароль пользователя
}

// LoginResponse Структура пользователя для фронтенда
type LoginResponse struct {
	AccessToken string `json:"access_token"` // Токен доступа
	Login       string `json:"login"`        // Логин пользователя
	Email       string `json:"email"`        // Email пользователя
	Name        string `json:"name"`         // Имя пользователя
	Tdid        string `json:"tdid"`         // Уникальный идентификатор пользователя
}

// Validate Валидация полей при регистрации
func (u *RegisterRequest) Validate() error {
	if u == nil {
		return errors.New("Error RegisterRequest")
	}

	email := u.Email
	if len(email) != len([]rune(email)) {
		return errors.New("Адрес электронной почты должен содержать сиволы только латинского алфавита!")
	}

	login := u.Login
	if len(login) != len([]rune(login)) {
		return errors.New("Логин должен содержать сиволы только латинского алфавита!")
	} else if len([]rune(login)) < 4 {
		return errors.New("Логин должен содержать минимум 4 символа!")
	}

	password := u.Password
	if len(password) != len([]rune(password)) {
		return errors.New("Пароль должен содержать сиволы только латинского алфавита!")
	} else if len([]rune(password)) < 4 {
		return errors.New("Пароль должен содержать минимум 4 символа!")
	}

	return u.ValidEmail()
}

// ValidEmail Метод для проверки валидности email
func (u RegisterRequest) ValidEmail() error {
	err := checkmail.ValidateFormat(u.Email)

	return err
}
