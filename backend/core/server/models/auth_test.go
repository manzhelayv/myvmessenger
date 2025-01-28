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
		Phone    string
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
				Phone:    "7 777 777 77 77",
			},
			wantErr: assert.NoError,
		},
		"Не заполнен email": {
			fields: &fields{
				Email:    "",
				Login:    "test",
				Name:     "test",
				Password: "test",
				Phone:    "7 777 777 77 77",
			},
			wantErr: assert.Error,
		},
		"Не заполнен login": {
			fields: &fields{
				Email:    "test@mail.ru",
				Login:    "",
				Name:     "test",
				Password: "test",
				Phone:    "7 777 777 77 77",
			},
			wantErr: assert.Error,
		},
		"Не заполнен password": {
			fields: &fields{
				Email:    "test@mail.ru",
				Login:    "test",
				Name:     "test",
				Password: "",
				Phone:    "7 777 777 77 77",
			},
			wantErr: assert.Error,
		},
		"email сиволы не латиница": {
			fields: &fields{
				Email:    "вввв@mail.ru",
				Login:    "test",
				Name:     "test",
				Password: "test",
				Phone:    "7 777 777 77 77",
			},
			wantErr: assert.Error,
		},
		"login сиволы не латиница": {
			fields: &fields{
				Email:    "test@mail.ru",
				Login:    "вввв",
				Name:     "test",
				Password: "test",
				Phone:    "7 777 777 77 77",
			},
			wantErr: assert.Error,
		},
		"password сиволы не латиница": {
			fields: &fields{
				Email:    "test@mail.ru",
				Login:    "test",
				Name:     "test",
				Password: "вввв",
				Phone:    "7 777 777 77 77",
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
				Phone:    tt.fields.Phone,
			}
			tt.wantErr(t, u.Validate(), fmt.Sprintf("Validate()"))
		})
	}
}

func TestRegisterRequest_ValidEmail(t *testing.T) {
	type fields struct {
		Email string
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
				Email: "test@mail.ru",
			},
			wantErr: assert.NoError,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			tt.wantErr(t, ValidEmail(tt.fields.Email), fmt.Sprintf("ValidEmail()"))
		})
	}
}

func TestUpdateRequest_UpdateUserValidate(t *testing.T) {
	type fields struct {
		Email string
		Login string
		Name  string
		Tdid  string
		Phone string
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
				Email: "test@mail.ru",
				Login: "test",
				Name:  "test",
				Phone: "7 777 777 77 77",
			},
			wantErr: assert.NoError,
		},
		"Не заполнен email": {
			fields: fields{
				Email: "",
				Login: "test",
				Name:  "test",
				Phone: "7 777 777 77 77",
			},
			wantErr: assert.Error,
		},
		"Не заполнен login": {
			fields: fields{
				Email: "test@mail.ru",
				Login: "",
				Name:  "test",
				Phone: "7 777 777 77 77",
			},
			wantErr: assert.Error,
		},
		"email сиволы не латиница": {
			fields: fields{
				Email: "вввв@mail.ru",
				Login: "test",
				Name:  "test",
				Phone: "7 777 777 77 77",
			},
			wantErr: assert.Error,
		},
		"login сиволы не латиница": {
			fields: fields{
				Email: "test@mail.ru",
				Login: "вввв",
				Name:  "test",
				Phone: "7 777 777 77 77",
			},
			wantErr: assert.Error,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			u := &UpdateRequest{
				Email: tt.fields.Email,
				Login: tt.fields.Login,
				Name:  tt.fields.Name,
				Tdid:  tt.fields.Tdid,
				Phone: tt.fields.Phone,
			}
			tt.wantErr(t, u.UpdateUserValidate(), fmt.Sprintf("UpdateUserValidate()"))
		})
	}
}

func TestValidateUserData(t *testing.T) {
	type args struct {
		email string
		login string
		phone string
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
				email: "test@mail.ru",
				login: "test",
				phone: "7 777 777 77 77",
			},
			wantErr: assert.NoError,
		},
		"Не заполнен email": {
			args: args{
				email: "",
				login: "test",
				phone: "7 777 777 77 77",
			},
			wantErr: assert.Error,
		},
		"Не заполнен login": {
			args: args{
				email: "test@mail.ru",
				login: "",
				phone: "7 777 777 77 77",
			},
			wantErr: assert.Error,
		},
		"email сиволы не латиница": {
			args: args{
				email: "test@mail.ru",
				login: "test",
				phone: "7",
			},
			wantErr: assert.Error,
		},
		"login сиволы не латиница": {
			args: args{
				email: "test@mail.ru",
				login: "вввв",
				phone: "7 777 777 77 77",
			},
			wantErr: assert.Error,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			tt.wantErr(t, ValidateUserData(tt.args.email, tt.args.login, tt.args.phone), fmt.Sprintf("ValidateUserData(%v, %v, %v)", tt.args.email, tt.args.login, tt.args.phone))
		})
	}
}

func TestValidatePassword(t *testing.T) {
	type args struct {
		password string
	}
	tests := map[string]struct {
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		"Пустые данные": {
			args:    args{},
			wantErr: assert.Error,
		},
		"Заполненные данные, латиница": {
			args: args{
				password: "aaaa",
			},
			wantErr: assert.NoError,
		},
		"Заполненные данные, кирилица 2 символа": {
			args: args{
				password: "фф",
			},
			wantErr: assert.Error,
		},
		"Заполненные данные, кирилица 4 символа": {
			args: args{
				password: "фффф",
			},
			wantErr: assert.Error,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			tt.wantErr(t, ValidatePassword(tt.args.password), fmt.Sprintf("ValidatePassword(%v)", tt.args.password))
		})
	}
}

func TestValidateName(t *testing.T) {
	type args struct {
		name string
	}
	tests := map[string]struct {
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		"Пустые данные": {
			args:    args{},
			wantErr: assert.Error,
		},
		"Заполненные данные, латиница": {
			args: args{
				name: "aaaa",
			},
			wantErr: assert.NoError,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			tt.wantErr(t, ValidateName(tt.args.name), fmt.Sprintf("ValidateName(%v)", tt.args.name))
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
