package models

// Messages Структура сообщения чата, для таблицы чатов
type ChatsChanges struct {
	Before AfterOrBefore `json:"before" bson:"before"`
	After  AfterOrBefore `json:"after" bson:"after"`
	Op     string        `json:"op" bson:"op"` // Версия
}

type AfterOrBefore struct {
	Id        Oid           `json:"_id"`                        // Первичный ключ
	UserFrom  string        `json:"user_from" bson:"user_from"` // Уникальный идентификатор текущего пользователя
	UserTo    string        `json:"user_to" bson:"user_to"`     // Уникальный идентификатор пользователя которому отправили сообщение
	Message   string        `json:"message" bson:"message"`     // Сообщение
	File      string        `json:"file" bson:"file"`           // Если передан файл
	CreatedAt DateCreatedAt `json:"created_at"`                 // Дата создания сообщения
	UpdatedAt DateUpdatedAt `json:"updated_at"`                 // Дата редактирования сообщения
}

type Oid struct {
	Oid string `json:"$oid" bson:"_id"`
}

type DateCreatedAt struct {
	Date int64 `json:"$date" bson:"created_at"`
}

type DateUpdatedAt struct {
	Date int64 `json:"$date" bson:"updated_at"`
}

type AfterAndBefore struct {
	After  AfterString  `json:""`
	Before BeforeString `json:""`
	Op     Op           `json:""`
}

type Op struct {
	Op string `json:"op" bson:"op"`
}

type AfterString struct {
	After string `json:"after" bson:"after"`
}

type BeforeString struct {
	Before string `json:"before" bson:"before"`
}
