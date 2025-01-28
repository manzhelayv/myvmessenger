package http

import (
	"chat/database/drivers"
	"chat/models"
	redisCli "chat/servers/redis"
	"encoding/json"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	protobufF3 "gitlab.com/myvmessenger/client/f3-client/protobuf"
	"gitlab.com/myvmessenger/client/server-client/user/protobuf"
	"log"
	"net/http"
)

type Chat struct {
	db       drivers.DbInterfase
	redisCli *redisCli.Redis
	userGrpc protobuf.UsersClient
	f3Client protobufF3.FileClientClient
	infoLog  *log.Logger
	errorLog *log.Logger
}

func NewChat(db drivers.DbInterfase, redisCli *redisCli.Redis, userGrpc protobuf.UsersClient, f3GrpcClient protobufF3.FileClientClient, infoLog *log.Logger, errorLog *log.Logger) Chat {
	return Chat{
		db:       db,
		redisCli: redisCli,
		userGrpc: userGrpc,
		f3Client: f3GrpcClient,
		infoLog:  infoLog,
		errorLog: errorLog,
	}
}

func (m Chat) Routes() http.Handler {
	r := chi.NewRouter()

	r.Post("/", m.AddMessage)
	r.Get("/", m.GetMessages)
	r.Get("/chats", m.Chats)

	return r
}

// Chats Получение чатов пользователя
// @Summary Получение чатов
// @Description Получение чатов пользователя по tdid пользователя
// @Tags chat
// @Accept json
// @Produce json
// @Param Authorization header string true "Authorization"
// @Success 200 {array}  []models.Chats
// @Failure 400 {object} models.Response
// @Failure 403 {object} models.Response
// @Router /chat/chats [get]
func (c *Chat) Chats(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	tdid := ctx.Value("tdid").(string)

	messages, err := c.db.GetLastMessages(ctx, tdid)
	if err != nil {
		c.errorLog.Printf("Chats, method GetLastMessages: %s\n", err)

		_ = render.Render(w, r, models.BadRequest(err))
		return
	}

	if messages == nil {
		c.errorLog.Printf("Chats, method ChatNotFound: %s\n", models.ChatNotFound)

		_ = render.Render(w, r, models.BadRequest(models.ChatNotFound))
		return
	}

	usersTo, messMessage, messDate, filesMess := models.GetChatsData(messages)

	users, err := c.userGrpc.GetUsersFromTdid(ctx, usersTo)
	if err != nil {
		c.errorLog.Printf("Chats, method GetUsersFromTdid: %s\n", err)

		_ = render.Render(w, r, models.BadRequest(err))
		return
	}

	chats := models.GetChats(users, messMessage, messDate, filesMess)
	if chats == nil {
		c.errorLog.Printf("Chats, method GetChats: %s\n", models.ChatNotFound)

		_ = render.Render(w, r, models.BadRequest(models.ChatNotFound))
		return
	}

	files := models.AvatarProto(chats)
	if files == nil {
		c.errorLog.Printf("Chats, method AvatarProto is nil")

		_ = render.Render(w, r, models.BadRequest(models.ErrUserNilAvatar))
		return
	}

	images, err := c.f3Client.LoadFiles(ctx, files)
	switch err {
	case nil:
		models.GetAvatar(images, chats)

		render.JSON(w, r, chats)
		return
	case models.ErrEmptyImage:
		c.errorLog.Printf("Chats, method ErrEmptyImage: %s\n", err)

		_ = render.Render(w, r, models.BadRequest(err))
		return
	case models.ErrUserDoNotHaveAccess:
		c.errorLog.Printf("Chats, method AccessDenied: %s\n", err)

		_ = render.Render(w, r, models.AccessDenied(err))
		return
	default:
		c.errorLog.Printf("Chats, method Internal: %s\n", err)

		_ = render.Render(w, r, models.Internal(err))
		return
	}

	render.JSON(w, r, chats)
}

// GetMessages Получение сообщений чата
// @Summary Получение сообщений
// @Description Получение сообщений чата с пользователем по user_to
// @Tags chat
// @Accept json
// @Produce json
// @Param Authorization header string true "Authorization"
// @Success 200 {object}  models.MessagesUser
// @Failure 400 {object} models.Response
// @Router /chat/{user_to} [get]
func (c *Chat) GetMessages(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userTo := r.URL.Query().Get("user_to")
	chat := &models.Chat{
		UserTo: userTo,
	}

	userFrom := ctx.Value("tdid")
	if userFrom == nil {
		c.errorLog.Printf("GetMessages, method ErrTokenNotFound: %s\n", models.ErrTokenNotFound)

		_ = render.Render(w, r, models.BadRequest(models.ErrTokenNotFound))
		return
	}

	chat.UserFrom = userFrom.(string)

	messages, err := c.db.GetMessages(ctx, chat.UserFrom, chat.UserTo)
	if err != nil {
		c.errorLog.Printf("GetMessages, method GetMessages: %s\n", err)

		_ = render.Render(w, r, models.BadRequest(err))
		return
	}

	models.GetMessagesDate(messages)

	render.JSON(w, r, models.MessagesUser{
		Messages: messages,
	})
}

// AddMessage Добавление сообщения
// @Summary Добавление сообщения
// @Description Добавление сообщения в чат пользователю по user_to
// @Tags chat
// @Accept json
// @Produce json
// @Param Authorization header string true "Authorization"
// @Param user_to query string true "пользователь в чате которому отправляется сообщение"
// @Param message query string true "сообщение"
// @Success 200 {array}  models.Chats
// @Failure 400 {object} models.Response
// @Router /chat [post]
func (c *Chat) AddMessage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	message := &models.Message{}
	if err := json.NewDecoder(r.Body).Decode(message); err != nil {
		c.errorLog.Printf("AddMessage, method Decode: %s\n", err)

		_ = render.Render(w, r, models.BadRequest(err))
		return
	}

	if message.Message == "" {
		_ = render.Render(w, r, models.BadRequest(models.EmptyMessage))
		return
	}

	userFrom := ctx.Value("tdid")
	if userFrom == nil {
		c.errorLog.Printf("AddMessage, method ErrTokenNotFound: %s\n", models.ErrTokenNotFound)

		_ = render.Render(w, r, models.BadRequest(models.ErrTokenNotFound))
		return
	}

	chat := models.GetMessagesForAdd(userFrom, message.Message)

	userTo := c.redisCli.GetUser(message.UserTo)
	if userTo != "" {
		chat.UserTo = userTo
	} else {
		userFromGrpc, err := c.userGrpc.GetUserTdid(ctx, &protobuf.User{Email: message.UserTo})
		if err != nil {
			c.errorLog.Printf("AddMessage, method GetUserTdid: %s\n", err)

			_ = render.Render(w, r, models.BadRequest(err))
			return
		}

		userTo = userFromGrpc.Tdid
		chat.UserTo = userTo
	}

	err := c.db.SetMessage(ctx, chat)
	if err != nil {
		c.errorLog.Printf("AddMessage, method SetMessage: %s\n", err)

		_ = render.Render(w, r, models.BadRequest(err))
		return
	}

	messages := models.AddMessagesForArray(chat)

	usersTo, messMessage, messDate, filesMess := models.GetChatsData(messages)

	users, err := c.userGrpc.GetUsersFromTdid(ctx, usersTo)
	if err != nil {
		c.errorLog.Printf("AddMessage, method GetUsersFromTdid: %s\n", err)

		_ = render.Render(w, r, models.BadRequest(err))
		return
	}

	chats := models.GetChats(users, messMessage, messDate, filesMess)
	if chats == nil {
		c.errorLog.Printf("AddMessage, method GetChats: %s\n", models.ChatNotFound)

		_ = render.Render(w, r, models.BadRequest(models.ChatNotFound))
		return
	}

	render.JSON(w, r, chats[0])
}
