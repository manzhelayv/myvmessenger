package models

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

var hashPassword = "N0OrfyFimTw9mnTNGceDNVSg$#$9$#$782f5a34e4bce4af1ee2a2a01490d9554c058f11370fee4b9e2778d5$#$cf44c7752a9f6de9ca526234efbc653f908f789982f45ddc9c4a30fc446b1ca6"
var userPassword = "1111"

func TestUsers_GeneratedTdid(t *testing.T) {
	date, _ := time.Parse(time.RFC3339, "2014-04-15T18:00:15-06:00")
	type args struct {
		max int
		min int
	}
	tests := map[string]struct {
		fields *Users
		args   args
	}{
		"Пустые данные": {
			fields: nil,
			args:   args{},
		},
		"Заполненные данные": {
			fields: &Users{
				Tdid:      "11111",
				Email:     "test@mail.ru",
				Name:      "test",
				Login:     "test",
				Password:  "test",
				CreatedAt: date,
				UpdatedAt: date,
			},
			args: args{
				max: 1000,
				min: 1,
			},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			tt.fields.GeneratedTdid(tt.args.max, tt.args.min)
		})
	}
}

func TestUsers_GetToken(t *testing.T) {
	tests := map[string]struct {
		fields  *Users
		wantErr assert.ErrorAssertionFunc
	}{
		"Пустые данные": {
			fields:  nil,
			wantErr: assert.Error,
		},
		"Заполненные данные": {
			fields: &Users{
				Email: "test@mail.ru",
				Name:  "test",
				Tdid:  "11111",
				Login: "test",
			},
			wantErr: assert.NoError,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			_, err := tt.fields.GetToken()
			if !tt.wantErr(t, err, fmt.Sprintf("GetToken()")) {
				return
			}
		})
	}
}

func TestUsers_ValidEmail(t *testing.T) {
	type fields struct {
		Tdid      string
		Email     string
		Name      string
		Login     string
		Password  string
		CreatedAt time.Time
		UpdatedAt time.Time
	}
	tests := map[string]struct {
		fields  fields
		wantErr assert.ErrorAssertionFunc
	}{
		"Пустые данные": {
			fields:  fields{},
			wantErr: assert.Error,
		},
		"Заполненные данные": {
			fields: fields{
				Email:    "test@mail.ru",
				Login:    "test",
				Name:     "test",
				Password: "test",
			},
			wantErr: assert.NoError,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			u := Users{
				Tdid:      tt.fields.Tdid,
				Email:     tt.fields.Email,
				Name:      tt.fields.Name,
				Login:     tt.fields.Login,
				Password:  tt.fields.Password,
				CreatedAt: tt.fields.CreatedAt,
				UpdatedAt: tt.fields.UpdatedAt,
			}
			tt.wantErr(t, u.ValidEmail(), fmt.Sprintf("ValidEmail()"))
		})
	}
}

func TestUsers_GetDataUserForInsert(t *testing.T) {
	date, _ := time.Parse(time.RFC3339, "2014-04-15T18:00:15-06:00")

	data := make(map[string]interface{}, 7)
	data["Email"] = "test@mail.ru"
	data["Name"] = "test"
	data["Login"] = "test"
	data["Password"] = "test"
	data["Tdid"] = "11111"
	data["CreatedAt"] = date
	data["UpdatedAt"] = date

	tests := map[string]struct {
		fields  *Users
		want    map[string]interface{}
		wantErr assert.ErrorAssertionFunc
	}{
		"Пустые данные": {
			fields:  nil,
			want:    nil,
			wantErr: assert.Error,
		},
		"Заполненные данные": {
			fields: &Users{
				Email:     "test@mail.ru",
				Name:      "test",
				Login:     "test",
				Password:  "test",
				Tdid:      "11111",
				CreatedAt: date,
				UpdatedAt: date,
			},
			want:    data,
			wantErr: assert.NoError,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := tt.fields.GetDataUserForInsert()
			if !tt.wantErr(t, err, fmt.Sprintf("GetDataUserForInsert()")) {
				return
			}
			assert.Equalf(t, tt.want, got, "GetDataUserForInsert()")
		})
	}
}

func TestFindUser_Validate(t *testing.T) {
	type fields struct {
		EmailOrLogin string
		Password     string
	}
	tests := map[string]struct {
		fields  *fields
		wantErr assert.ErrorAssertionFunc
	}{
		"Пустые данные": {
			fields:  &fields{},
			wantErr: assert.Error,
		},
		"Заполненные данные": {
			fields: &fields{
				EmailOrLogin: "test",
				Password:     "test",
			},
			wantErr: assert.NoError,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			f := FindUser{
				EmailOrLogin: tt.fields.EmailOrLogin,
				Password:     tt.fields.Password,
			}
			tt.wantErr(t, f.Validate(), fmt.Sprintf("Validate()"))
		})
	}
}

func TestHashPassword(t *testing.T) {
	type args struct {
		psw string
	}
	tests := map[string]struct {
		args    args
		want    string
		wantErr assert.ErrorAssertionFunc
	}{
		"Пустые данные": {
			args:    args{},
			want:    "",
			wantErr: assert.Error,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := HashPassword(tt.args.psw)
			if !tt.wantErr(t, err, fmt.Sprintf("HashPassword(%v)", tt.args.psw)) {
				return
			}
			assert.Equalf(t, tt.want, got, "HashPassword(%v)", tt.args.psw)
		})
	}
}

func TestVerifyHashPassword(t *testing.T) {
	type args struct {
		hash string
		psw  string
	}
	tests := map[string]struct {
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		"Пустые данные": {
			args:    args{},
			wantErr: assert.Error,
		},
		"Заполненные данные": {
			args: args{
				hash: hashPassword,
				psw:  userPassword,
			},
			wantErr: assert.NoError,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			tt.wantErr(t, VerifyHashPassword(tt.args.hash, tt.args.psw), fmt.Sprintf("VerifyHashPassword(%v, %v)", tt.args.hash, tt.args.psw))
		})
	}
}
