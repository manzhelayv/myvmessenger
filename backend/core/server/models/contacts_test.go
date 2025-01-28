package models

import (
	"github.com/stretchr/testify/assert"
	"gitlab.com/myvmessenger/client/f3-client/protobuf"
	"testing"
	"time"
)

func TestAvatarProto(t *testing.T) {
	type args struct {
		contactsReturn []*Contacts
	}
	tests := map[string]struct {
		args args
		want *protobuf.Files
	}{
		"Пустые данные": {
			args: args{},
			want: nil,
		},
		"Заполненные данные без картинки": {
			args: args{
				contactsReturn: []*Contacts{
					{
						Image:  "",
						UserTo: "11111",
					},
				},
			},
			want: &protobuf.Files{
				Files: []*protobuf.File{
					{
						Tdid: "11111",
						Name: EMPTY_AVATAR,
					},
				},
			},
		},
		"Заполненные данные 1": {
			args: args{
				contactsReturn: []*Contacts{
					{
						Image:  "111.png",
						UserTo: "11111",
					},
				},
			},
			want: &protobuf.Files{
				Files: []*protobuf.File{
					{
						Tdid: "11111",
						Name: "111.png",
					},
				},
			},
		},
		"Заполненные данные 2": {
			args: args{
				contactsReturn: []*Contacts{
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
			want: &protobuf.Files{
				Files: []*protobuf.File{
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
			assert.Equalf(t, tt.want, AvatarProto(tt.args.contactsReturn), "AvatarProto(%v)", tt.args.contactsReturn)
		})
	}
}

func TestContactsAvatar(t *testing.T) {
	type args struct {
		contactsReturn []*Contacts
		images         *protobuf.Files
	}
	tests := map[string]struct {
		args args
		want []*Contacts
	}{
		"Пустые данные": {
			args: args{},
			want: nil,
		},
		"Заполненные данные без картинки": {
			args: args{
				contactsReturn: []*Contacts{
					{
						Avatar: []byte{},
					},
				},
				images: &protobuf.Files{
					Files: []*protobuf.File{
						{
							Content: []byte{},
						},
					},
				},
			},
			want: []*Contacts{
				{
					Avatar: []byte{},
				},
			},
		},
		"Заполненные данные 1": {
			args: args{
				contactsReturn: []*Contacts{
					{
						UserTo: "11111",
						Avatar: []byte{},
					},
				},
				images: &protobuf.Files{
					Files: []*protobuf.File{
						{
							Tdid:    "11111",
							Content: []byte("111"),
						},
					},
				},
			},
			want: []*Contacts{
				{
					UserTo: "11111",
					Avatar: []byte("111"),
				},
			},
		},
		"Заполненные данные 2": {
			args: args{
				contactsReturn: []*Contacts{
					{
						UserTo: "11111",
						Avatar: []byte("111"),
					},
					{
						UserTo: "22222",
						Avatar: []byte("222"),
					},
				},
				images: &protobuf.Files{
					Files: []*protobuf.File{
						{
							Tdid:    "22222",
							Content: []byte("333"),
						},
						{
							Tdid:    "11111",
							Content: []byte("444"),
						},
					},
				},
			},
			want: []*Contacts{
				{
					UserTo: "11111",
					Avatar: []byte("444"),
				},
				{
					UserTo: "22222",
					Avatar: []byte("333"),
				},
			},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			ContactsAvatar(tt.args.contactsReturn, tt.args.images)
			assert.Equalf(t, tt.want, tt.args.contactsReturn, "ContactsAvatar(%v)", tt.args.contactsReturn)

		})
	}
}

func TestGetPhones(t *testing.T) {
	type args struct {
		contacts RequestContactsPhones
	}
	tests := map[string]struct {
		args args
		want []string
	}{
		"Пустые данные": {
			args: args{},
			want: nil,
		},
		"Заполненные данные 1": {
			args: args{
				contacts: RequestContactsPhones{
					Contacts: []ContactByPhones{
						{
							Phone: "7 777 777 77 77",
							Name:  "111",
						},
					},
				},
			},
			want: []string{
				"7 777 777 77 77",
			},
		},
		"Заполненные данные 2": {
			args: args{
				contacts: RequestContactsPhones{
					Contacts: []ContactByPhones{
						{
							Phone: "7 777 777 77 77",
							Name:  "111",
						},
						{
							Phone: "8 888 888 88 88",
							Name:  "111",
						},
					},
				},
			},
			want: []string{
				"7 777 777 77 77",
				"8 888 888 88 88",
			},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			assert.Equalf(t, tt.want, GetPhones(tt.args.contacts), "GetPhones(%v)", tt.args.contacts)
		})
	}
}

func TestGetTdids(t *testing.T) {
	type args struct {
		users []User
	}
	tests := map[string]struct {
		args args
		want []string
	}{
		"Пустые данные": {
			args: args{},
			want: nil,
		},
		"Заполненные данные 1": {
			args: args{
				users: []User{
					{
						Tdid: "1111",
					},
				},
			},
			want: []string{
				"1111",
			},
		},
		"Заполненные данные 2": {
			args: args{
				users: []User{
					{
						Tdid: "2222",
					},
					{
						Tdid: "1111",
					},
				},
			},
			want: []string{
				"2222",
				"1111",
			},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			assert.Equalf(t, tt.want, GetTdids(tt.args.users), "GetTdids(%v)", tt.args.users)
		})
	}
}

func TestGetContacts(t *testing.T) {
	now := time.Now().Format(time.RFC3339)
	date, _ := time.Parse(time.RFC3339, now)

	type args struct {
		contactsIsset []Contact
		users         []User
		tdid          string
		contacts      RequestContactsPhones
	}
	tests := map[string]struct {
		args args
		want []Contact
	}{
		"Пустые данные": {
			args: args{},
			want: nil,
		},
		"Заполненные данные 1": {
			args: args{
				contactsIsset: []Contact{
					{
						UserTo:   "2222",
						UserFrom: "3333",
					},
				},
				users: []User{
					{
						Tdid:  "4444",
						Phone: "777 777 77 77",
					},
				},
				contacts: RequestContactsPhones{
					Contacts: []ContactByPhones{
						{
							Phone: "777 777 77 77",
							Name:  "Name 1",
						},
					},
				},
				tdid: "1111",
			},
			want: []Contact{
				{
					UserFrom:   "1111",
					UserTo:     "4444",
					UserToName: "Name 1",
					CreatedAt:  date,
					UpdatedAt:  date,
				},
			},
		},
		"Заполненные данные 2": {
			args: args{
				contactsIsset: []Contact{
					{
						UserTo:   "1111",
						UserFrom: "3333",
					},
					{
						UserTo:   "1111",
						UserFrom: "6666",
					},
					{
						UserTo:   "1111",
						UserFrom: "7777",
					},
				},
				users: []User{
					{
						Tdid:  "4444",
						Phone: "7 777 777 77 77",
					},
					{
						Tdid:  "5555",
						Phone: "8 888 888 88 88",
					},
					{
						Tdid:  "7777",
						Phone: "9 999 999 999 999",
					},
				},
				contacts: RequestContactsPhones{
					Contacts: []ContactByPhones{
						{
							Phone: "7 777 777 77 77",
							Name:  "Name 1",
						},
						{
							Phone: "8 888 888 88 88",
							Name:  "Name 2",
						},
						{
							Phone: "9 999 999 999 999",
							Name:  "Name 3",
						},
					},
				},
				tdid: "1111",
			},
			want: []Contact{
				{
					UserFrom:   "1111",
					UserTo:     "4444",
					UserToName: "Name 1",
					CreatedAt:  date,
					UpdatedAt:  date,
				},
				{
					UserFrom:   "1111",
					UserTo:     "5555",
					UserToName: "Name 2",
					CreatedAt:  date,
					UpdatedAt:  date,
				},
			},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			GetContacts(tt.args.contactsIsset, tt.args.users, tt.args.tdid, tt.args.contacts)
			cont := GetContacts(tt.args.contactsIsset, tt.args.users, tt.args.tdid, tt.args.contacts)
			assert.Equalf(t, tt.want, cont, "GetContacts(%v, %v, %v, %v)", tt.args.contactsIsset, tt.args.users, tt.args.tdid, tt.args.contacts)
		})
	}
}
