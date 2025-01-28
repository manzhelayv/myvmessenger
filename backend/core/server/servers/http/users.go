package http

import (
	"encoding/json"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"gitlab.com/myvmessenger/client/f3-client/protobuf"
	"log"
	"net/http"
	"server/database/drivers"
	"server/models"
	redisCli "server/servers/redis"
)

type Users struct {
	db       drivers.DbInterfase
	redisCli *redisCli.Redis
	f3Client protobuf.FileClientClient
	infoLog  *log.Logger
	errorLog *log.Logger
}

func NewUsers(db drivers.DbInterfase, redisCli *redisCli.Redis, f3Client protobuf.FileClientClient, infoLog *log.Logger, errorLog *log.Logger) *Users {
	return &Users{
		db:       db,
		redisCli: redisCli,
		f3Client: f3Client,
		infoLog:  infoLog,
		errorLog: errorLog,
	}
}

func (u Users) Routes() http.Handler {
	r := chi.NewRouter()

	r.Post("/user", u.Register)
	r.Post("/login", u.Login)
	r.Get("/user/image", u.GetImage)
	r.Put("/user", u.UpdateUser)
	r.Put("/updatepassword", u.UpdatePassword)

	return r
}

// Register Добавление Регистрация пользователя
// @Summary Регистрация пользователя
// @Description Регистрация пользователя по email, name, login, password
// @Tags users
// @Accept json
// @Produce json
// @Param email query string true "email пользователя"
// @Param name query string true "Имя пользователя"
// @Param login query string true "Логин пользователя"
// @Param password query string true "Пароль пользователя"
// @Param phone query string true "Телефон пользователя"
// @Success 200 {array}  models.LoginResponse
// @Failure 400 {object} models.Response
// @Router /user [post]
func (u *Users) Register(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	regReq := models.RegisterRequest{}
	if err := json.NewDecoder(r.Body).Decode(&regReq); err != nil {
		u.errorLog.Printf("Register, method Decode: %s\n", err)

		_ = render.Render(w, r, models.BadRequest(err))
		return
	}

	err := regReq.Validate()
	if err != nil {
		if err.Error() == models.ErrMailNoAddres {
			err = models.ErrMailNoAddresResponse
		}

		u.errorLog.Printf("Register, method Validate: %s\n", err)

		_ = render.Render(w, r, models.BadRequest(err))
		return
	}

	user := models.ConvertStructUsersCreate(regReq)

	err = user.GeneratedTdid(10000000000, 1)
	if err != nil {
		u.errorLog.Printf("Register, method GeneratedTdid: %s\n", err)

		_ = render.Render(w, r, models.BadRequest(err))
		return
	}

	password, err := models.HashPassword(user.Password)
	if err != nil {
		u.errorLog.Printf("Register, method HashPassword: %s\n", err)

		_ = render.Render(w, r, models.BadRequest(err))
		return
	}

	user.Password = password

	err = u.db.InserUser(ctx, user)
	if err != nil {
		u.errorLog.Printf("Register, method InserUser: %s\n", err)

		_ = render.Render(w, r, models.BadRequest(err))
		return
	}

	data, err := user.GetDataUserForInsert()
	if err != nil {
		u.errorLog.Printf("Register, method GetDataUserForInsert: %s\n", err)

		_ = render.Render(w, r, models.BadRequest(err))
		return
	}

	_, err = u.redisCli.CreateUser(data)
	if err != nil {
		u.errorLog.Printf("Register, method CreateUser, Ошибка редиса: %s\n", err)

		log.Println("Ошибка редиса: ", err)
	}

	t, err := user.GetToken()
	if err != nil {
		u.errorLog.Printf("Register, method GetToken: %s\n", err)

		_ = render.Render(w, r, models.BadRequest(err))
		return
	}

	render.JSON(w, r, models.LoginResponse{
		AccessToken: t,
		Login:       user.Login,
		Email:       user.Email,
		Name:        user.Name,
		Tdid:        user.Tdid,
		Phone:       user.Phone,
	})
}

// Login Авторизация пользователя
// @Summary Авторизация пользователя
// @Description Авторизация пользователя по email или login и password
// @Tags users
// @Accept json
// @Produce json
// @Param email_or_login query string true "Логин пользователя или email"
// @Param password query string true "Пароль пользователя"
// @Success 200 {array}  models.LoginResponse
// @Failure 400 {object} models.Response
// @Failure 401 {object} models.Response
// @Router /login [post]
func (u *Users) Login(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userFind := &models.FindUser{}
	if err := json.NewDecoder(r.Body).Decode(userFind); err != nil {
		u.errorLog.Printf("Login, method Decode: %s\n", err)

		_ = render.Render(w, r, models.BadRequest(err))
		return
	}

	err := userFind.Validate()
	if err != nil {
		u.errorLog.Printf("Login, method Validate: %s\n", err)

		_ = render.Render(w, r, models.BadRequest(err))
		return
	}

	emailOrLogin := userFind.EmailOrLogin

	user := u.redisCli.GetUser(emailOrLogin)
	if user != nil {
		err = models.VerifyHashPassword(user.Password, userFind.Password)
		if err == nil {
			t, err := user.GetToken()
			if err != nil {
				u.errorLog.Printf("Login, method GetUser: %s\n", err)

				_ = render.Render(w, r, models.BadRequest(err))
				return
			}

			render.JSON(w, r, models.LoginResponse{
				AccessToken: t,
				Login:       user.Login,
				Email:       user.Email,
				Name:        user.Name,
				Tdid:        user.Tdid,
				Phone:       user.Phone,
			})
			return
		}
	}

	userPg, err := u.db.GetUserEmailOrLogin(ctx, userFind.EmailOrLogin)
	if err != nil {
		if err.Error() == models.ErrNoResultLogin {
			err = models.ErrNoResultResponse
		}

		u.errorLog.Printf("Login, method GetUserEmailOrLogin: %s\n", err)

		_ = render.Render(w, r, models.Unauthorized(err))
		return
	}

	err = models.VerifyHashPassword(userPg.Password, userFind.Password)
	if err != nil {
		u.errorLog.Printf("Login, method VerifyHashPassword: %s\n", err)

		_ = render.Render(w, r, models.Unauthorized(err))
		return
	}

	t, err := userPg.GetToken()
	if err != nil {
		u.errorLog.Printf("Login, method GetToken: %s\n", err)

		_ = render.Render(w, r, models.BadRequest(err))
		return
	}

	render.JSON(w, r, models.LoginResponse{
		AccessToken: t,
		Login:       userPg.Login,
		Email:       userPg.Email,
		Name:        userPg.Name,
		Tdid:        userPg.Tdid,
		Phone:       userPg.Phone,
	})
}

// GetImage Получение картинки пользователя
// @Summary Получение картинки пользователя
// @Description Получение картинки пользователя
// @Tags users
// @Accept json
// @Produce json
// @Success 200 {array}  models.ProfileResponse
// @Failure 400 {object} models.Response
// @Router /user/image/{user_to} [get]
func (u *Users) GetImage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	tdid := r.URL.Query().Get("user_to")

	profile, err := u.db.GetProfile(ctx, tdid)
	if err != nil {
		u.errorLog.Printf("GetImage, method GetProfile: %s\n", err)

		_ = render.Render(w, r, models.BadRequest(err))
		return
	}

	if profile.Image == "" {
		profile.Image = models.EMPTY_AVATAR
	}

	file := models.FilesProto(tdid, profile)

	image, err := u.f3Client.LoadFiles(ctx, file)
	if err != nil {
		u.errorLog.Printf("GetImage, method LoadFiles: %s\n", err)

		_ = render.Render(w, r, models.BadRequest(err))
		return
	}

	render.JSON(w, r, &models.ProfileResponse{
		Avatar: image.Files[0].Content,
	})
}

// UpdateUser Обновление данных пользователя
// @Summary Обновление данных пользователя
// @Description Обновление данных пользователя
// @Tags users
// @Accept json
// @Produce json
// @Param email query string true "email пользователя"
// @Param name query string true "Имя пользователя"
// @Param login query string true "Логин пользователя"
// @Param phone query string true "Телефон пользователя"
// @Success 200 {array}  models.LoginResponse
// @Failure 400 {object} models.Response
// @Router /user [put]
func (u *Users) UpdateUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	regReq := models.UpdateRequest{}
	if err := json.NewDecoder(r.Body).Decode(&regReq); err != nil {
		u.errorLog.Printf("UpdateUser, method Decode: %s\n", err)

		_ = render.Render(w, r, models.BadRequest(err))
		return
	}

	err := regReq.UpdateUserValidate()
	if err != nil {
		if err.Error() == models.ErrMailNoAddres {
			err = models.ErrMailNoAddresResponse
		}

		u.errorLog.Printf("UpdateUser, method UpdateUserValidate: %s\n", err)

		_ = render.Render(w, r, models.BadRequest(err))
		return
	}

	user := models.ConvertStructUsersUpdate(regReq)

	err = u.db.UpdateUser(ctx, user)
	if err != nil {
		u.errorLog.Printf("UpdateUser, method UpdateUser: %s\n", err)

		_ = render.Render(w, r, models.BadRequest(err))
		return
	}

	data, err := user.GetDataUserForInsert()
	if err != nil {
		u.errorLog.Printf("UpdateUser, method GetDataUserForInsert: %s\n", err)

		_ = render.Render(w, r, models.BadRequest(err))
		return
	}

	_, err = u.redisCli.CreateUser(data)
	if err != nil {
		u.errorLog.Printf("UpdateUser, method CreateUser, Ошибка редиса: %s\n", err)
		log.Println("Ошибка редиса: ", err)
	}

	t, err := user.GetToken()
	if err != nil {
		u.errorLog.Printf("UpdateUser, method GetToken: %s\n", err)

		_ = render.Render(w, r, models.BadRequest(err))
		return
	}

	render.JSON(w, r, models.LoginResponse{
		AccessToken: t,
		Login:       user.Login,
		Email:       user.Email,
		Name:        user.Name,
		Tdid:        user.Tdid,
		Phone:       user.Phone,
	})
}

// UpdatePassword Изменение пароля пользователя
// @Summary Изменение пароля пользователя
// @Description Изменение пароля пользователя
// @Tags users
// @Accept json
// @Produce json
// @Param newPassword query string true "Новый пароль пользователя"
// @Param oldPassword query string true "Старый пароль пользователя"
// @Success 200
// @Failure 400 {object} models.Response
// @Failure 401 {object} models.Response
// @Router /updatepassword [put]
func (u *Users) UpdatePassword(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userFind := &models.UpdatePassword{}
	if err := json.NewDecoder(r.Body).Decode(userFind); err != nil {
		u.errorLog.Printf("UpdatePassword, method Decode: %s\n", err)

		_ = render.Render(w, r, models.BadRequest(err))
		return
	}

	err := userFind.Validate()
	if err != nil {
		u.errorLog.Printf("UpdatePassword, method Validate: %s\n", err)

		_ = render.Render(w, r, models.BadRequest(err))
		return
	}

	emailOrLogin := userFind.EmailOrLogin

	user := u.redisCli.GetUser(emailOrLogin)
	if user != nil {
		err = models.VerifyHashPassword(user.Password, userFind.Password)
		if err == nil {
			user.Password = userFind.Password

			err = u.db.UpdatePassword(ctx, user)
			if err != nil {
				u.errorLog.Printf("UpdatePassword, method UpdatePassword: %s\n", err)

				_ = render.Render(w, r, models.BadRequest(err))
				return
			}

			return
		}
	}

	userPg, err := u.db.GetUserEmailOrLogin(ctx, userFind.EmailOrLogin)
	if err != nil {
		if err.Error() == models.ErrNoResultLogin {
			err = models.ErrNoResultResponse
		}

		u.errorLog.Printf("UpdatePassword, method GetUserEmailOrLogin: %s\n", err)

		_ = render.Render(w, r, models.Unauthorized(err))
		return
	}

	err = models.VerifyHashPassword(userPg.Password, userFind.Password)
	if err != nil {
		u.errorLog.Printf("UpdatePassword, method VerifyHashPassword: %s\n", err)

		_ = render.Render(w, r, models.Unauthorized(err))
		return
	}

	password, err := models.HashPassword(userFind.NewPassword)
	if err != nil {
		u.errorLog.Printf("UpdatePassword, method HashPassword: %s\n", err)

		_ = render.Render(w, r, models.BadRequest(err))
		return
	}

	user = models.ConvertUsersPasswordUpdate(password, userPg)

	err = u.db.UpdatePassword(ctx, user)
	if err != nil {
		u.errorLog.Printf("UpdatePassword, method UpdatePassword: %s\n", err)

		_ = render.Render(w, r, models.BadRequest(err))
		return
	}

	return
}
