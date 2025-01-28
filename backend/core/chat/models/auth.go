package models

var JwtSecretKey = []byte("secret-key")

const (
	ClientType    = "Client-Type"
	ClientVersion = "Client-Version"
	DeviceID      = "Device-id"
	Channel       = "Channel"
	Affiliation   = "Affiliation"
	Uuid          = "Uuid"
)

func AllowedHeaders() []string {
	return []string{
		ClientType,
		ClientVersion,
		DeviceID,
		Channel,
		Affiliation,
		Uuid,
	}
}

// LoginResponse Структура пользователя для фронтенда
type LoginResponse struct {
	AccessToken string `json:"access_token"` // Токен доступа
	Login       string `json:"login"`        // Логин пользователя
	Email       string `json:"email"`        // Email пользователя
	Name        string `json:"name"`         // Имя пользователя
	Tdid        string `json:"tdid"`         // Уникальный идентификатор пользователя
}
