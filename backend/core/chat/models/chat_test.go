package models

import (
	"github.com/stretchr/testify/assert"
	protobufF3 "gitlab.com/myvmessenger/client/f3-client/protobuf"
	"gitlab.com/myvmessenger/client/server-client/user/protobuf"
	"log"
	"testing"
	"time"
)

func TestGetChatsData(t *testing.T) {
	date, _ := time.Parse(time.RFC3339, "2014-04-15T18:00:15-00:00")
	type args struct {
		messages []*Messages
	}
	tests := map[string]struct {
		args  args
		want  *protobuf.UserTdids
		want1 map[string]string
		want2 map[string]time.Time
		want3 map[string]string
	}{
		"Пустые данные": {
			args:  args{},
			want:  nil,
			want1: nil,
			want2: nil,
			want3: nil,
		},
		"Заполненные данные 1": {
			args: args{
				[]*Messages{
					{
						UserFrom:  "11111",
						UserTo:    "22222",
						Message:   "Привет!",
						CreatedAt: date,
						UpdatedAt: date,
						File:      "111.png",
					},
				},
			},
			want: &protobuf.UserTdids{
				Tdid: []string{"22222"},
			},
			want1: map[string]string{
				"22222": "Привет!",
			},
			want2: map[string]time.Time{
				"22222": date,
			},
			want3: map[string]string{
				"22222": "111.png",
			},
		},
		"Заполненные данные 2": {
			args: args{
				[]*Messages{
					{
						UserFrom:  "11111",
						UserTo:    "22222",
						Message:   "Привет1!",
						CreatedAt: date,
						UpdatedAt: date,
						File:      "111.png",
					},
					{
						UserFrom:  "33333",
						UserTo:    "44444",
						Message:   "Привет2!",
						CreatedAt: date,
						UpdatedAt: date,
						File:      "222.png",
					},
				},
			},
			want: &protobuf.UserTdids{
				Tdid: []string{"22222", "44444"},
			},
			want1: map[string]string{
				"22222": "Привет1!",
				"44444": "Привет2!",
			},
			want2: map[string]time.Time{
				"22222": date,
				"44444": date,
			},
			want3: map[string]string{
				"22222": "111.png",
				"44444": "222.png",
			},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			got, got1, got2, got3 := GetChatsData(tt.args.messages)
			assert.Equalf(t, tt.want, got, "GetChatsData(%v)", tt.args.messages)
			assert.Equalf(t, tt.want1, got1, "GetChatsData(%v)", tt.args.messages)
			assert.Equalf(t, tt.want2, got2, "GetChatsData(%v)", tt.args.messages)
			assert.Equalf(t, tt.want3, got3, "GetChatsData(%v)", tt.args.messages)
		})
	}
}

func TestGetChats(t *testing.T) {
	date1, _ := time.Parse(time.RFC3339, "2014-04-15T18:00:15-00:00")
	date2, _ := time.Parse(time.RFC3339, "2014-04-15T19:00:15-00:00")

	dateUser := make(map[string]time.Time, 2)
	dateUser["11111"] = date1
	dateUser["22222"] = date2

	dateDayUser1, dateTimeUser1 := GetdateToChats("11111", dateUser)
	dateDayUser2, dateTimeUser2 := GetdateToChats("22222", dateUser)

	type args struct {
		users   *protobuf.UsersItem
		message map[string]string
		date    map[string]time.Time
		files   map[string]string
	}
	tests := map[string]struct {
		args args
		want []*Chats
	}{
		"Пустые данные": {
			args: args{},
			want: nil,
		},
		"Заполненные данные 1": {
			args: args{
				users: &protobuf.UsersItem{
					Users: []*protobuf.UserItem{
						{
							Email:  "test@mail.ru",
							Tdid:   "11111",
							Name:   "Иван",
							Login:  "test",
							Avatar: "111.png",
						},
					},
				},
				message: map[string]string{
					"11111": "Привет!",
				},
				date: map[string]time.Time{
					"11111": date1,
				},
				files: map[string]string{
					"11111": "111.png",
				},
			},
			want: []*Chats{
				{
					UserFrom:      "11111",
					UserTo:        "11111",
					Message:       "Привет!",
					Name:          "Иван",
					Email:         "test@mail.ru",
					Date:          dateUser["11111"],
					File:          "111.png",
					Image:         "111.png",
					DateTimestamp: dateUser["11111"].UnixMilli(),
					DateDay:       dateDayUser1,
					DateTime:      dateTimeUser1,
				},
			},
		},
		"Заполненные данные 2": {
			args: args{
				users: &protobuf.UsersItem{
					Users: []*protobuf.UserItem{
						{
							Email:  "test1@mail.ru",
							Tdid:   "11111",
							Name:   "Иван 1",
							Login:  "test1",
							Avatar: "111",
						},
						{
							Email:  "test2@mail.ru",
							Tdid:   "22222",
							Name:   "Иван 2",
							Login:  "test2",
							Avatar: "222",
						},
					},
				},
				message: map[string]string{
					"11111": "Привет1!",
					"22222": "Привет2!",
				},
				date: map[string]time.Time{
					"11111": date1,
					"22222": date2,
				},
				files: map[string]string{
					"11111": "111.png",
					"22222": "222.png",
				},
			},
			want: []*Chats{
				{
					UserFrom:      "11111",
					UserTo:        "11111",
					Message:       "Привет1!",
					Name:          "Иван 1",
					Email:         "test1@mail.ru",
					Date:          dateUser["11111"],
					File:          "111.png",
					Image:         "111",
					DateTimestamp: dateUser["11111"].UnixMilli(),
					DateDay:       dateDayUser1,
					DateTime:      dateTimeUser1,
				},
				{
					UserFrom:      "22222",
					UserTo:        "22222",
					Message:       "Привет2!",
					Name:          "Иван 2",
					Email:         "test2@mail.ru",
					Date:          dateUser["22222"],
					File:          "222.png",
					Image:         "222",
					DateTimestamp: dateUser["22222"].UnixMilli(),
					DateDay:       dateDayUser2,
					DateTime:      dateTimeUser2,
				},
			},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			assert.Equalf(t, tt.want, GetChats(tt.args.users, tt.args.message, tt.args.date, tt.args.files), "GetChats(%v, %v, %v)", tt.args.users, tt.args.message, tt.args.date)
		})
	}
}

func TestGetdateToChats(t *testing.T) {
	date1, _ := time.Parse(time.RFC3339, "2014-04-15T13:00:15-00:00")
	date2, _ := time.Parse(time.RFC3339, "2014-04-15T18:00:15-00:00")

	dateUser1 := make(map[string]time.Time, 2)
	dateUser1["11111"] = date1

	dateUser2 := make(map[string]time.Time, 2)
	dateUser2["22222"] = date2

	type args struct {
		tdid string
		date map[string]time.Time
	}
	tests := map[string]struct {
		args  args
		want  string
		want1 string
	}{
		"Пустые данные": {
			args:  args{},
			want:  "",
			want1: "",
		},
		"Заполненные данные": {
			args: args{
				tdid: "11111",
				date: dateUser1,
			},
			want:  "15.04.2014",
			want1: "19:00",
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			got, got1 := GetdateToChats(tt.args.tdid, tt.args.date)
			assert.Equalf(t, tt.want, got, "GetdateToChats(%v, %v)", tt.args.tdid, tt.args.date)
			assert.Equalf(t, tt.want1, got1, "GetdateToChats(%v, %v)", tt.args.tdid, tt.args.date)
		})
	}
}

func TestAvatarProto(t *testing.T) {
	type args struct {
		chats []*Chats
	}
	tests := map[string]struct {
		args args
		want *protobufF3.Files
	}{
		"Пустые данные": {
			args: args{},
			want: nil,
		},
		"Заполненные данные без картинки": {
			args: args{
				chats: []*Chats{
					{
						Image:  "",
						UserTo: "11111",
					},
				},
			},
			want: &protobufF3.Files{
				Files: []*protobufF3.File{
					{
						Tdid: "11111",
						Name: EMPTY_AVATAR,
					},
				},
			},
		},
		"Заполненные данные 1": {
			args: args{
				chats: []*Chats{
					{
						Image:  "111.png",
						UserTo: "11111",
					},
				},
			},
			want: &protobufF3.Files{
				Files: []*protobufF3.File{
					{
						Tdid: "11111",
						Name: "111.png",
					},
				},
			},
		},
		"Заполненные данные 2": {
			args: args{
				chats: []*Chats{
					{
						Image:  "111.png",
						UserTo: "11111",
					},
					{
						Image:  "222.png",
						UserTo: "22222",
					},
				},
			},
			want: &protobufF3.Files{
				Files: []*protobufF3.File{
					{
						Tdid: "11111",
						Name: "111.png",
					},
					{
						Tdid: "22222",
						Name: "222.png",
					},
				},
			},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			assert.Equalf(t, tt.want, AvatarProto(tt.args.chats), "AvatarProto(%v)", tt.args.chats)
		})
	}
}

func TestGetAvatar(t *testing.T) {
	type args struct {
		images *protobufF3.Files
		chats  []*Chats
	}
	tests := map[string]struct {
		args args
		want []*Chats
	}{
		"Пустые данные": {
			args: args{},
			want: nil,
		},
		"Заполненные данные 1": {
			args: args{
				images: &protobufF3.Files{
					Files: []*protobufF3.File{
						{
							Tdid:    "11111",
							Content: []byte("aaa"),
						},
					},
				},
				chats: []*Chats{
					{
						UserTo: "11111",
					},
				},
			},
			want: []*Chats{
				{
					UserTo: "11111",
					Avatar: []byte("aaa"),
				},
			},
		},
		"Заполненные данные 2": {
			args: args{
				images: &protobufF3.Files{
					Files: []*protobufF3.File{
						{
							Tdid:    "11111",
							Content: []byte("aaa"),
						},
						{
							Tdid:    "22222",
							Content: []byte("bbb"),
						},
					},
				},
				chats: []*Chats{
					{
						UserTo: "11111",
					},
					{
						UserTo: "22222",
					},
				},
			},
			want: []*Chats{
				{
					UserTo: "11111",
					Avatar: []byte("aaa"),
				},
				{
					UserTo: "22222",
					Avatar: []byte("bbb"),
				},
			},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			GetAvatar(tt.args.images, tt.args.chats)
			assert.Equalf(t, tt.want, tt.args.chats, "GetAvatar(%v)", tt.args.chats)
		})
	}
}

func TestGetMessagesDate(t *testing.T) {
	date1, _ := time.Parse(time.RFC3339, "2014-04-15T13:00:15-00:00")
	date2, _ := time.Parse(time.RFC3339, "2014-04-15T18:00:15-00:00")

	type args struct {
		messages []*Messages
	}
	tests := map[string]struct {
		args args
		want []*Messages
	}{
		"Пустые данные": {
			args: args{},
			want: nil,
		},
		"Заполненные данные 1": {
			args: args{
				messages: []*Messages{
					{
						CreatedAt: date1,
						UpdatedAt: date1,
					},
				},
			},
			want: []*Messages{
				{
					CreatedAt:     date1,
					UpdatedAt:     date1,
					DateTimestamp: 1397566815000,
					Date:          "15 Апреля 2014 г.",
					DateTime:      "19:00",
				},
			},
		},
		"Заполненные данные 2": {
			args: args{
				messages: []*Messages{
					{
						CreatedAt: date1,
						UpdatedAt: date1,
					},
					{
						CreatedAt: date1,
						UpdatedAt: date2,
					},
				},
			},
			want: []*Messages{
				{
					CreatedAt:     date1,
					UpdatedAt:     date1,
					DateTimestamp: 1397566815000,
					Date:          "15 Апреля 2014 г.",
					DateTime:      "19:00",
				},
				{
					CreatedAt:     date1,
					UpdatedAt:     date2,
					DateTimestamp: 1397584815000,
					Date:          "16 Апреля 2014 г.",
					DateTime:      "00:00",
				},
			},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			GetMessagesDate(tt.args.messages)
			assert.Equalf(t, tt.want, tt.args.messages, "GetMessagesDate(%v)", tt.args.messages)
		})
	}
}

func TestGetMessagesForAdd(t *testing.T) {
	loc, _ := time.LoadLocation("Asia/Almaty")

	now := time.Now().In(loc).Format(time.RFC3339)
	date, err := time.Parse(time.RFC3339, now)
	if err != nil {
		log.Println("[ERROR] Time parse GetMessagesForAdd", err.Error())
	}

	type args struct {
		userFrom any
		message  string
	}
	tests := map[string]struct {
		args args
		want *Messages
	}{
		"Пустые данные": {
			args: args{},
			want: nil,
		},
		"Заполненные данные": {
			args: args{
				userFrom: "11111",
				message:  "aaa",
			},
			want: &Messages{
				UserFrom:  "11111",
				Message:   "aaa",
				CreatedAt: date,
				UpdatedAt: date,
			},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			assert.Equalf(t, tt.want, GetMessagesForAdd(tt.args.userFrom, tt.args.message), "GetMessagesForAdd(%v, %v)", tt.args.userFrom, tt.args.message)
		})
	}
}

func TestAddMessagesForArray(t *testing.T) {
	type args struct {
		chat *Messages
	}
	tests := map[string]struct {
		args args
		want []*Messages
	}{
		"Пустые данные": {
			args: args{},
			want: []*Messages{nil},
		},
		"Заполненные данные": {
			args: args{
				chat: &Messages{
					UserFrom: "11111",
					UserTo:   "22222",
				},
			},
			want: []*Messages{
				{
					UserFrom: "11111",
					UserTo:   "22222",
				},
			},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			assert.Equalf(t, tt.want, AddMessagesForArray(tt.args.chat), "AddMessagesForArray(%v)", tt.args.chat)
		})
	}
}
