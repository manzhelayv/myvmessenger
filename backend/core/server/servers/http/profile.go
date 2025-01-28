package http

import (
	"encoding/base64"
	"encoding/json"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"gitlab.com/myvmessenger/client/f3-client/protobuf"
	"log"
	"net/http"
	"server/database/drivers"
	"server/models"
	redisCli "server/servers/redis"
	"strings"
)

const ERROR_AVATAR_SIZE = "ResourceExhausted desc"

type Profile struct {
	db       drivers.DbInterfase
	redisCli *redisCli.Redis
	f3Client protobuf.FileClientClient
	infoLog  *log.Logger
	errorLog *log.Logger
}

func NewProfile(db drivers.DbInterfase, redisCli *redisCli.Redis, f3GrpcClient protobuf.FileClientClient, infoLog *log.Logger, errorLog *log.Logger) *Profile {
	return &Profile{
		db:       db,
		redisCli: redisCli,
		f3Client: f3GrpcClient,
		infoLog:  infoLog,
		errorLog: errorLog,
	}
}

func (p Profile) Routes() http.Handler {
	r := chi.NewRouter()

	r.Post("/", p.Profile)
	r.Get("/", p.GetProfile)

	return r
}

// GetProfile Получение профайла пользователя
// @Summary Получение профайла пользователя
// @Description Получение профайла
// @Tags profile
// @Accept json
// @Produce json
// @Success 200 {array}  models.ProfileResponse
// @Failure 400 {object} models.Response
// @Failure 403 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /profile [get]
func (p *Profile) GetProfile(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	tdid := ctx.Value("tdid").(string)

	profile, err := p.db.GetProfile(ctx, tdid)
	if err != nil {
		p.errorLog.Printf("GetProfile, method GetProfile: %s\n", err)

		_ = render.Render(w, r, models.BadRequest(err))
		return
	}

	file := models.FilesProto(tdid, profile)

	image, err := p.f3Client.LoadFiles(ctx, file)
	switch err {
	case nil:
		render.JSON(w, r, &models.ProfileResponse{
			Avatar: image.Files[0].Content,
		})
		return
	case models.ErrEmptyImage:
		p.errorLog.Printf("GetProfile, method LoadFiles: %s\n", err)

		_ = render.Render(w, r, models.BadRequest(err))
		return
	case models.ErrUserDoNotHaveAccess:
		p.errorLog.Printf("GetProfile, method LoadFiles: %s\n", err)

		_ = render.Render(w, r, models.AccessDenied(err))
		return
	default:
		p.errorLog.Printf("GetProfile, method LoadFiles: %s\n", err)

		_ = render.Render(w, r, models.Internal(err))
		return
	}
}

// Profile Изменение профайла пользователя
// @Summary Изменение профайла пользователя
// @Description Изменение профайла
// @Tags profile
// @Accept json
// @Produce json
// @Param image query string true "Фото пользователя"
// @Success 200 {array}  protobuf.ID
// @Failure 400 {object} models.Response
// @Router /profile [post]
func (p *Profile) Profile(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	tdid := ctx.Value("tdid").(string)

	f := &models.File{}
	if err := json.NewDecoder(r.Body).Decode(&f); err != nil {
		p.errorLog.Printf("Profile, method Decode: %s\n", err)

		_ = render.Render(w, r, models.BadRequest(err))
		return
	}

	f.Image = f.Image[strings.IndexByte(f.Image, ',')+1:]

	file := models.FileProto(f)

	content, err := base64.StdEncoding.DecodeString(f.Image)
	if err != nil {
		p.errorLog.Printf("Profile, method DecodeString: %s\n", err)

		_ = render.Render(w, r, models.BadRequest(err))
		return
	}

	file.Content = content

	imagef3, err := p.f3Client.UploadFile(ctx, file)
	if err != nil {
		contain := strings.Contains(err.Error(), ERROR_AVATAR_SIZE)
		switch contain {
		case true:
			p.errorLog.Printf("Profile, method UploadFile: %s\n", err)

			err = models.ErrorSizeAvatar

			p.errorLog.Printf("Profile, method UploadFile: %s\n", err)

			_ = render.Render(w, r, models.BadRequest(err))
			return
		default:
			p.errorLog.Printf("Profile, method UploadFile: %s\n", err)

			_ = render.Render(w, r, models.BadRequest(err))
			return
		}
	}

	user, err := p.db.GetUsersFromTdid(ctx, []string{tdid})
	if err != nil {
		p.errorLog.Printf("Profile, method GetUsersFromTdid: %s\n", err)

		_ = render.Render(w, r, models.BadRequest(err))
		return
	}

	err = p.db.UpdateProfile(ctx, user[0].Id, imagef3.Id)
	if err != nil {
		p.errorLog.Printf("Profile, method UpdateProfile: %s\n", err)

		_ = render.Render(w, r, models.BadRequest(err))
		return
	}

	render.JSON(w, r, imagef3)
}
