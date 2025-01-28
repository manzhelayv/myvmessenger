package models

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	password "github.com/vzglad-smerti/password_hash"
	"math/rand/v2"
	"net/mail"
	"strconv"
	"time"
)

var (
	ErrNoResultResponse      = errors.New("Вы ввели неверные данные")
	ErrFindUserEmailValidate = errors.New("Введите логин или email и пароль")
	ErrUsersNil              = errors.New("Структура пользовател пустая") // Изменить текст
	ErrTokenUser             = errors.New("[ERROR] GetToken Users = nil")
	ErrPasswordNil           = errors.New("Введите пароль")
)

// Структура данных с информацией о пользователе для аунтификации
type User struct {
	Tdid      string    // Уникальный идентификатор пользователя
	Email     string    // Email пользователя
	Name      string    // Имя пользователя
	Login     string    // Логин пользователя
	Password  string    // Пароль пользователя
	CreatedAt time.Time // Дата создания пользователя
	UpdatedAt time.Time // Дата редактирования пользователя
	jwt.RegisteredClaims
}

// Структура данных с информацией о пользователе для БД
type Users struct {
	Tdid      string    `bson:"tdid"`       // Уникальный идентификатор пользователя
	Email     string    `bson:"email"`      // Email пользователя
	Name      string    `bson:"name"`       // Имя пользователя
	Login     string    `bson:"login"`      // Логин пользователя
	Password  string    `bson:"password"`   // Пароль пользователя
	CreatedAt time.Time `bson:"created_at"` // Дата создания пользователя
	UpdatedAt time.Time `bson:"updated_at"` // Дата редактирования пользователя
}

// FindUser Структура для авторизации пользователя
type FindUser struct {
	EmailOrLogin string `json:"email_or_login"`
	Password     string `json:"password"`
}

// GeneratedTdid Генерация случайного tdid пользователя
func (u *Users) GeneratedTdid(max, min int) error {
	if u == nil {
		return errors.New("[ERROR] Структура Users = nil")
	}

	tdid := rand.IntN(max-min) + min

	strTdid := strconv.Itoa(tdid)

	u.Tdid = strTdid

	return nil
}

// GetDataUserForInsert Метод для создания данных о пользователе для redis
func (u *Users) GetDataUserForInsert() (map[string]interface{}, error) {
	if u == nil {
		return nil, ErrUsersNil
	}

	data := make(map[string]interface{}, 7)
	data["Email"] = u.Email
	data["Name"] = u.Name
	data["Login"] = u.Login
	data["Password"] = u.Password
	data["Tdid"] = u.Tdid
	data["CreatedAt"] = u.CreatedAt
	data["UpdatedAt"] = u.UpdatedAt

	return data, nil
}

// GetToken создание токена пользователя
func (u *Users) GetToken() (string, error) {
	if u == nil {
		return "", ErrTokenUser
	}

	payload := User{
		Email: u.Email,
		Name:  u.Name,
		Tdid:  u.Tdid,
		Login: u.Login,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)

	t, err := token.SignedString([]byte("secret-key")) // Для тестов - []byte("secret-key") // JwtSecretKey // secret-key
	if err != nil {
		return "", err
	}

	return t, err
}

// ValidEmail Метод для проверки валидности email
func (u Users) ValidEmail() error {
	_, err := mail.ParseAddress(u.Email)

	return err
}

// Validate Валидация полученных данных для авторизации пользователя
func (f FindUser) Validate() error {
	if f.EmailOrLogin == "" && f.Password == "" {
		return ErrFindUserEmailValidate
	}

	return nil
}

// HashPassword получение хэша пароля
func HashPassword(psw string) (string, error) {
	if psw == "" {
		return "", ErrPasswordNil
	}

	hashPassword, err := password.Hash(psw)
	if err != nil {
		return "", err
	}

	return hashPassword, nil
}

// VerifyHashPassword проверка хэша пароля с сохраненным в базы данных
func VerifyHashPassword(hash, psw string) error {
	if hash == "" || psw == "" {
		return ErrNoResultResponse
	}

	hash_veriry, err := password.Verify(hash, psw)
	if err != nil {
		return err
	}

	if hash_veriry {
		return nil
	}

	return ErrNoResultResponse
}
