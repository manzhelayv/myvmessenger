package models

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRegisterRequest_Validate(t *testing.T) {
	type fields struct {
		Email    string
		Login    string
		Name     string
		Password string
	}
	tests := map[string]struct {
		fields  *fields
		wantErr assert.ErrorAssertionFunc
	}{
		"Пустые данные": {
			fields:  nil,
			wantErr: assert.Error,
		},
		"Заполненные данные": {
			fields: &fields{
				Email:    "test@mail.ru",
				Login:    "test",
				Name:     "test",
				Password: "test",
			},
			wantErr: assert.NoError,
		},
		"Не заполнен email": {
			fields: &fields{
				Email:    "",
				Login:    "test",
				Name:     "test",
				Password: "test",
			},
			wantErr: assert.Error,
		},
		"Не заполнен login": {
			fields: &fields{
				Email:    "test@mail.ru",
				Login:    "",
				Name:     "test",
				Password: "test",
			},
			wantErr: assert.Error,
		},
		"Не заполнен password": {
			fields: &fields{
				Email:    "test@mail.ru",
				Login:    "test",
				Name:     "test",
				Password: "",
			},
			wantErr: assert.Error,
		},
		"email сиволы не латиница": {
			fields: &fields{
				Email:    "вввв@mail.ru",
				Login:    "test",
				Name:     "test",
				Password: "test",
			},
			wantErr: assert.Error,
		},
		"login сиволы не латиница": {
			fields: &fields{
				Email:    "test@mail.ru",
				Login:    "вввв",
				Name:     "test",
				Password: "test",
			},
			wantErr: assert.Error,
		},
		"password сиволы не латиница": {
			fields: &fields{
				Email:    "test@mail.ru",
				Login:    "test",
				Name:     "test",
				Password: "вввв",
			},
			wantErr: assert.Error,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			if tt.fields == nil {
				var u *RegisterRequest
				tt.wantErr(t, u.Validate(), fmt.Sprintf("Validate()"))
				return
			}
			u := &RegisterRequest{
				Email:    tt.fields.Email,
				Login:    tt.fields.Login,
				Name:     tt.fields.Name,
				Password: tt.fields.Password,
			}
			tt.wantErr(t, u.Validate(), fmt.Sprintf("Validate()"))
		})
	}
}

func TestRegisterRequest_ValidEmail(t *testing.T) {
	type fields struct {
		Email    string
		Login    string
		Name     string
		Password string
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
			u := RegisterRequest{
				Email:    tt.fields.Email,
				Login:    tt.fields.Login,
				Name:     tt.fields.Name,
				Password: tt.fields.Password,
			}
			tt.wantErr(t, u.ValidEmail(), fmt.Sprintf("ValidEmail()"))
		})
	}
}

func TestAllowedHeaders(t *testing.T) {
	tests := map[string]struct {
		want []string
	}{
		"Добиваем -cover до 100% :)": {
			want: []string{
				ClientType,
				ClientVersion,
				DeviceID,
				Channel,
				Affiliation,
				Uuid,
			},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			assert.Equalf(t, tt.want, AllowedHeaders(), "AllowedHeaders()")
		})
	}
}
