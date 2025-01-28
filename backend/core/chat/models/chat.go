package models

import (
	"errors"
	protobufF3 "gitlab.com/myvmessenger/client/f3-client/protobuf"
	"gitlab.com/myvmessenger/client/server-client/user/protobuf"
	"log"
	"strconv"
	"time"
)

var ErrorFindMessages = errors.New("В данном чате нет сообщений")

const EMPTY_AVATAR = "profile_0b5d723c6e8fcf7dd8f301f93a19a518"

// MessageWS Структура websocket
type MessageWS struct {
	UserTo  string `json:"userto"`  // Собеседник
	Message string `json:"message"` // Сообщение
	File    string `json:"file"`
	Status  string `json:"status"` // Статус
}

// StatusUser Структура статуса пользователя в сети
type StatusUser struct {
	Message string `json:"message"` // Сообщение
	Status  string `json:"status"`  // Статус
}

// MessagesUser Структура HTTP ответа для получения переписки с определенным пользователем
type MessagesUser struct {
	Messages []*Messages `json:"messages"` // Сообщения
}

// Message Структура для HTTP запроса, для добавления сообщения в чат
type Message struct {
	UserTo  string `json:"user_to"` // Уникальный идентификатор пользователя которому отправили сообщение
	Message string `json:"message"` // Сообщение
}

// Chat Структура для получения сообщений для чата с бд
type Chat struct {
	UserFrom string // Уникальный идентификатор текущего пользователя
	UserTo   string // Уникальный идентификатор собеседника
}

// Messages Структура сообщения чата, для таблицы чатов
type Messages struct {
	Id            string    `json:"_id" bson:"_id"`
	UserFrom      string    `json:"user_from" bson:"user_from"`           // Уникальный идентификатор текущего пользователя
	UserTo        string    `json:"user_to" bson:"user_to"`               // Уникальный идентификатор пользователя которому отправили сообщение
	Message       string    `json:"message" bson:"message"`               // Сообщение
	File          string    `json:"file" bson:"file"`                     // Если передан файл
	CreatedAt     time.Time `json:"created_at" bson:"created_at"`         // Дата создания сообщения
	UpdatedAt     time.Time `json:"updated_at" bson:"updated_at"`         // Дата редактирования сообщения
	DateTimestamp int64     `json:"date_timestamp" bson:"date_timestamp"` // Дата в timestamp
	Date          string    `json:"date" bson:"date"`                     // День, месяц, год отправленного сообщения
	DateTime      string    `json:"date_time" bson:"date_time"`           // Часы, минуты отправленного сообщения
}

// Chats Структура открытых чатов текущего пользователя
type Chats struct {
	UserFrom      string    `json:"user_from" bson:"user_from"`             // Уникальный идентификатор текущего пользователя
	UserTo        string    `json:"user_to" bson:"user_to"`                 // Уникальный идентификатор собеседника
	Message       string    `json:"message" bson:"message"`                 // Последнее сообщение в текущем чате
	File          string    `json:"file" bson:"file"`                       // Если передан файл
	Name          string    `json:"name" bson:"name"`                       // Имя собеседника
	Email         string    `json:"email" bson:"email"`                     // Email собеседника
	Image         string    `json:"image,omitempty" bson:"image,omitempty"` // Название фото пользователя, которому отпраляются сообщения в f3
	Avatar        []byte    `json:"avatar"`                                 // Аватарка пользователя, которому отпраляются сообщения из f3
	Date          time.Time `json:"date" bson:"date"`                       // Дата и время последнего сообщения
	DateTimestamp int64     `json:"date_timestamp" bson:"date_timestamp"`   // Дата в timestamp
	DateDay       string    `json:"date_day" bson:"date_day"`               // День, месяц, год или "Вчера" последнего сообщения
	DateTime      string    `json:"date_time" bson:"date_time"`             // Время последнего сообщения
}

// GetChatsData Преобразование массива структур Messages для
// usersTo - для структуры grpc запроса, messMessage для последнего сообщения, messDate - для даты последнего сообщения
func GetChatsData(messages []*Messages) (*protobuf.UserTdids, map[string]string, map[string]time.Time, map[string]string) {
	if messages == nil {
		return nil, nil, nil, nil
	}

	usersTo := &protobuf.UserTdids{}
	messMessage := make(map[string]string, len(messages))
	messDate := make(map[string]time.Time, len(messages))
	messFiles := make(map[string]string, len(messages))

	for _, message := range messages {
		usersTo.Tdid = append(usersTo.Tdid, message.UserTo)
		messMessage[message.UserTo] = message.Message
		messDate[message.UserTo] = message.CreatedAt
		messFiles[message.UserTo] = message.File
	}

	return usersTo, messMessage, messDate, messFiles
}

// GetChats Преобразование protobuf, хэш таблиц message, date в структуру Chats для HTTP ответа
func GetChats(users *protobuf.UsersItem, message map[string]string, date map[string]time.Time, files map[string]string) []*Chats {
	if users == nil || len(message) != len(users.Users) || len(date) != len(users.Users) {
		return nil
	}

	chats := make([]*Chats, 0, len(users.Users))
	for _, user := range users.Users {
		messageUser := message[user.Tdid]
		if len([]rune(messageUser)) >= 60 {
			runeStr := []rune(messageUser)
			messageUser = string(runeStr[:60]) + "..."
		}

		file := files[user.Tdid]

		dateDay, dateTime := GetdateToChats(user.Tdid, date)

		chats = append(chats, &Chats{
			UserFrom:      user.Tdid,
			UserTo:        user.Tdid,
			Message:       messageUser,
			File:          file,
			Name:          user.Name,
			Image:         user.Avatar,
			Email:         user.Email,
			Date:          date[user.Tdid],
			DateTimestamp: date[user.Tdid].UnixMilli(),
			DateDay:       dateDay,
			DateTime:      dateTime,
		})
	}

	return chats
}

// GetdateToChats Преобразует дату для вывода на фронт
func GetdateToChats(tdid string, date map[string]time.Time) (string, string) {
	if date == nil {
		return "", ""
	}

	loc, _ := time.LoadLocation("Asia/Almaty")

	now := time.Now().In(loc)
	date[tdid] = date[tdid].In(loc)

	dateDayFormat := "2"
	timeNow := now.Format(dateDayFormat)
	dateNow := date[tdid].Format(dateDayFormat)

	dateDayFormatDMY := "02.01.2006"
	dateDay := date[tdid].Format(dateDayFormatDMY)

	dateTimeFormatHM := "15:04"
	dateTime := date[tdid].Format(dateTimeFormatHM)

	timeNowInt, err := strconv.Atoi(timeNow)
	if err != nil {
		log.Println(err)
	}

	dateDayInt, err := strconv.Atoi(dateNow)
	if err != nil {
		log.Println(err)
	}

	if timeNowInt-dateDayInt == 1 {
		dateDay = "Вчера"
		dateTime = "0"
	} else if timeNowInt-dateDayInt == 0 {
		dateDay = ""
	}

	return dateDay, dateTime
}

func GetStringMonth(month string) string {
	switch month {
	case "1":
		return "Января"
	case "2":
		return "Феврлаля"
	case "3":
		return "Марта"
	case "4":
		return "Апреля"
	case "5":
		return "Мая"
	case "6":
		return "Июня"
	case "7":
		return "Июля"
	case "8":
		return "Августа"
	case "9":
		return "Сентября"
	case "10":
		return "Октября"
	case "11":
		return "Ноября"
	case "12":
		return "Декабря"
	}

	return ""
}

// AvatarProto Функция для преобразования структуры Contacts в f3 струтуру Files
func AvatarProto(chats []*Chats) *protobufF3.Files {
	if chats == nil || len(chats) == 0 {
		return nil
	}

	file := make([]*protobufF3.File, 0, len(chats))
	for _, chat := range chats {
		img := chat.Image
		if img == "" {
			img = EMPTY_AVATAR
		}

		file = append(file, &protobufF3.File{
			Tdid: chat.UserTo,
			Name: img,
		})
	}

	files := &protobufF3.Files{
		Files: file,
	}

	return files
}

// GetAvatar Получение аватарки пользователя с proto в структуру Chats
func GetAvatar(images *protobufF3.Files, chats []*Chats) {
	if images == nil {
		return
	}

	for _, img := range images.Files {
		for _, chat := range chats {
			if chat.UserTo == img.Tdid {
				chat.Avatar = img.Content
			}
		}
	}
}

// GetMessagesDate Преобразование дат сообщений
func GetMessagesDate(messages []*Messages) {
	if messages == nil || len(messages) == 0 {
		return
	}

	for _, message := range messages {
		loc, _ := time.LoadLocation("Asia/Almaty")

		message.DateTimestamp = message.UpdatedAt.UnixMilli()

		day := message.UpdatedAt.In(loc).Format("2")
		month := GetStringMonth(message.UpdatedAt.Format("1"))
		year := message.UpdatedAt.In(loc).Format("2006")

		message.Date = day + " " + month + " " + year + " г."
		message.DateTime = message.UpdatedAt.In(loc).Format("15:04")
	}
}

// GetMessagesForAdd Преобразование в структуру для добавления сообщения
func GetMessagesForAdd(userFrom any, message string) *Messages {
	if userFrom == nil {
		return nil
	}

	loc, _ := time.LoadLocation("Asia/Almaty")

	now := time.Now().In(loc).Format(time.RFC3339)
	date, err := time.Parse(time.RFC3339, now)
	if err != nil {
		log.Println("[ERROR] Time parse GetMessagesForAdd", err.Error())
	}

	chat := &Messages{
		UserFrom:  userFrom.(string),
		Message:   message,
		CreatedAt: date,
		UpdatedAt: date,
	}

	return chat
}

// AddMessagesForArray Добавление сообщения в массив сообщений
func AddMessagesForArray(chat *Messages) []*Messages {
	messages := []*Messages{
		chat,
	}

	return messages
}
