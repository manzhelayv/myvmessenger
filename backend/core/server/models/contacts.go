package models

import (
	"errors"
	"gitlab.com/myvmessenger/client/f3-client/protobuf"
	"time"
)

var ErrJsonRequestContact = errors.New("Введите логин или email пользователя")

const EMPTY_AVATAR = "profile_0b5d723c6e8fcf7dd8f301f93a19a518"

// Contacts Структура таблицы контактов, с добавлением email пользователя с таблицы пользователей,
// для получения из базы данных, таблицы контактов
type Contacts struct {
	UserFrom   string    `json:"user_from" bson:"user_from"`             // Уникальный идентификатор текущего пользователя
	UserTo     string    `json:"user_to" bson:"user_to"`                 // Уникальный идентификатор пользователя добавленного в контакты
	UserToName string    `json:"user_to_name" bson:"user_to_name"`       // Имя пользователя добавленного в контакты
	Email      string    `json:"email" bson:"email"`                     // Email пользователя, которому отпраляются сообщения
	Image      string    `json:"image,omitempty" bson:"image,omitempty"` // Название фото пользователя, которому отпраляются сообщения в f3
	Avatar     []byte    `json:"avatar"`                                 // Аватарка пользователя, которому отпраляются сообщения из f3
	CreatedAt  time.Time `json:"created_at" bson:"created_at"`           // Дата создания контакта
	UpdatedAt  time.Time `json:"updated_at" bson:"updated_at"`           // Дата редактирования контакта
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

// RequestContactsPhones Структура HTTP-запроса для добавления контактов по номерам телефонов
type RequestContactsPhones struct {
	Contacts []ContactByPhones `json:"contacts"` // Контакты пользователя с HTTP-запроса
}

// ContactByPhones Структура HTTP-запроса для добавления контактов по номерам телефонов
type ContactByPhones struct {
	Phone string `json:"phone"` // Телефон пользователя
	Name  string `json:"name"`  // Имя пользователя
}

// AvatarProto Функция для преобразования структуры Contacts в f3 струтуру Files
func AvatarProto(contactsReturn []*Contacts) *protobuf.Files {
	if contactsReturn == nil || len(contactsReturn) == 0 {
		return nil
	}

	file := make([]*protobuf.File, 0, len(contactsReturn))
	for _, contact := range contactsReturn {
		img := contact.Image
		if img == "" {
			img = EMPTY_AVATAR
		}

		file = append(file, &protobuf.File{
			Tdid: contact.UserTo,
			Name: img,
		})
	}

	files := &protobuf.Files{
		Files: file,
	}

	return files
}

// ContactsAvatar Функция для преобразования f3 струтуры Files в структуру Contacts
func ContactsAvatar(contactsReturn []*Contacts, images *protobuf.Files) {
	if images == nil || len(images.Files) == 0 || contactsReturn == nil || len(contactsReturn) == 0 {
		return
	}

	for _, img := range images.Files {
		for _, contact := range contactsReturn {
			if contact.UserTo == img.Tdid {
				contact.Avatar = img.Content
			}
		}
	}
}

// GetPhones Функция для обработки полученных с фронта телефонов и запись их в слайс
func GetPhones(contacts RequestContactsPhones) []string {
	if len(contacts.Contacts) == 0 {
		return nil
	}

	phones := make([]string, 0, len(contacts.Contacts))
	for _, c := range contacts.Contacts {
		phones = append(phones, c.Phone)
	}

	return phones
}

// GetTdids Функция для обработки tdid пользователей и запись их в слайс
func GetTdids(users []User) []string {
	if len(users) == 0 {
		return nil
	}

	tdids := make([]string, 0, len(users))
	for _, u := range users {
		tdids = append(tdids, u.Tdid)
	}

	return tdids
}

// GetContacts Функция для получения контактов пользователей, для записи их в базу данных
func GetContacts(contactsIsset []Contact, users []User, tdid string, contacts RequestContactsPhones) []Contact {
	if users == nil || len(users) == 0 {
		return nil
	}

	isset := false

	usersIn := []User{}
	for _, u := range users {
		isset = false

		for _, c := range contactsIsset {
			if u.Tdid == c.UserFrom || u.Tdid == c.UserTo {
				isset = true

				break
			}
		}

		if isset == false {
			usersIn = append(usersIn, u)
		}
	}

	now := time.Now().Format(time.RFC3339)
	date, _ := time.Parse(time.RFC3339, now)

	contactsIn := make([]Contact, 0, len(usersIn))
	for _, c := range contacts.Contacts {
		for _, u := range usersIn {
			if u.Phone == c.Phone {
				contactsIn = append(contactsIn, Contact{
					UserFrom:   tdid,
					UserTo:     u.Tdid,
					UserToName: c.Name,
					CreatedAt:  date,
					UpdatedAt:  date,
				})
			}
		}
	}

	return contactsIn
}
