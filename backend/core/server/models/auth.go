package models

import (
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
	Phone    string `json:"phone"`    // Телефон пользователя
}

// UpdateRequest Структура для изменения пользователя
type UpdateRequest struct {
	Email string `json:"email"` // Email пользователя
	Login string `json:"login"` // Логин пользователя
	Name  string `json:"name"`  // Имя пользователя
	Tdid  string `json:"tdid"`  // Tdid пользователя
	Phone string `json:"phone"` // Телефон пользователя
}

// LoginResponse Структура пользователя для фронтенда
type LoginResponse struct {
	AccessToken string `json:"access_token"` // Токен доступа
	Login       string `json:"login"`        // Логин пользователя
	Email       string `json:"email"`        // Email пользователя
	Name        string `json:"name"`         // Имя пользователя
	Tdid        string `json:"tdid"`         // Уникальный идентификатор пользователя
	Phone       string `json:"phone"`        // Телефон
}

// Validate Валидация полей при регистрации
func (u *RegisterRequest) Validate() error {
	if u == nil {
		return ErrorRegisterRequest
	}

	err := ValidateUserData(u.Email, u.Login, u.Phone)
	if err != nil {
		return err
	}

	err = ValidatePassword(u.Password)
	if err != nil {
		return err
	}

	err = ValidateName(u.Name)
	if err != nil {
		return err
	}

	return nil
}

// UpdateUserValidate Валидация полей при обновлении пользователя
func (u *UpdateRequest) UpdateUserValidate() error {
	if u == nil {
		return ErrorRegisterRequest
	}

	err := ValidateUserData(u.Email, u.Login, u.Phone)
	if err != nil {
		return err
	}

	return nil
}

// ValidEmail Метод для проверки валидности email
func ValidEmail(email string) error {
	err := checkmail.ValidateFormat(email)

	return err
}

// ValidateUserData Валидация email, login, phone пользователя
func ValidateUserData(email, login, phone string) error {
	if len(email) != len([]rune(email)) {
		return ErrorEmailRegistration
	}

	err := ValidEmail(email)
	if err != nil {
		return ErrorValidEmail
	}

	if len(login) != len([]rune(login)) {
		return ErrorLoginSimbolRegistration
	} else if len([]rune(login)) < 4 {
		return ErrorLoginCountRegistration
	}

	if len([]rune(phone)) != 15 {
		return ErrorPhoneRegistration
	}

	return nil
}

// ValidatePassword Валидация логина пользователя
func ValidatePassword(password string) error {
	if len(password) != len([]rune(password)) {
		return ErrorPasswordSimbolRegistration
	} else if len([]rune(password)) < 4 {
		return ErrorPasswordCountRegistration
	}

	return nil
}

// ValidateName Валидация имени пользователя
func ValidateName(name string) error {
	if len(name) == 0 {
		return ErrorNameUser
	}

	return nil
}
