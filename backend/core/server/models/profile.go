package models

import "gitlab.com/myvmessenger/client/f3-client/protobuf"

const PROFILE_IMAGE = "profile"

type Profile struct {
	Image  string `json:"image" bson:"image"`     // Картинка пользователя в f3 сервисе
	UserId int    `json:"user_id" bson:"user_id"` // Связанное поле с таблицей user
}

type ProfileResponse struct {
	Image  string `json:"image"`   // Картинка пользователя в f3 сервисе
	UserId int    `json:"user_id"` // Связанное поле с таблицей user
	Avatar []byte `json:"avatar"`  // Аватарка пользователя в f3 сервисе
}

// FilesProto Функция для получения f3 структуры Files
func FilesProto(tdid string, profile *Profile) *protobuf.Files {
	if profile == nil {
		return nil
	}

	img := profile.Image
	if img == "" {
		img = EMPTY_AVATAR
	}

	file := &protobuf.Files{
		Files: []*protobuf.File{
			{
				Tdid: tdid,
				Name: img,
			},
		},
	}

	return file
}

// FileProto Функция для получения f3 струтуры File
func FileProto(f *File) *protobuf.File {
	if f == nil {
		return nil
	}

	file := &protobuf.File{
		Name:    PROFILE_IMAGE,
		Content: []byte(f.Image),
	}

	return file
}
