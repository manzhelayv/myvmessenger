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
	"time"
)

type Contacts struct {
	dbPostgres drivers.DbInterfase
	redisCli   *redisCli.Redis
	f3Client   protobuf.FileClientClient
	infoLog    *log.Logger
	errorLog   *log.Logger
}

func NewContacts(dbPostgres drivers.DbInterfase, redisCli *redisCli.Redis, f3GrpcClient protobuf.FileClientClient, infoLog *log.Logger, errorLog *log.Logger) *Contacts {
	return &Contacts{
		dbPostgres: dbPostgres,
		redisCli:   redisCli,
		f3Client:   f3GrpcClient,
		infoLog:    infoLog,
		errorLog:   errorLog,
	}
}

func (c Contacts) Routes() http.Handler {
	r := chi.NewRouter()

	r.Get("/", c.GetContacts)
	r.Post("/", c.AddContact)
	r.Put("/", c.AddContacts)

	return r
}

// GetContacts Получение контактов пользователя
// @Summary Получение контактов
// @Description Получение контактов пользователя по tdid
// @Tags contacts
// @Accept json
// @Produce json
// @Success 200 {array}  models.Contacts
// @Failure 400 {object} models.Response
// @Router /contacts [get]
func (c *Contacts) GetContacts(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	tdid := ctx.Value("tdid").(string)

	contacts, err := c.dbPostgres.GetContacts(ctx, tdid)
	if err != nil {
		c.errorLog.Printf("GetContacts, method GetContacts: %s\n", err)

		_ = render.Render(w, r, models.BadRequest(err))
		return
	}

	files := models.AvatarProto(contacts)

	images, err := c.f3Client.LoadFiles(ctx, files)
	if err != nil {
		c.errorLog.Printf("GetContacts, method LoadFiles: %s\n", err)

		_ = render.Render(w, r, models.BadRequest(err))
		return
	}

	models.ContactsAvatar(contacts, images)

	render.JSON(w, r, contacts)
}

// AddContact Добавление контакта пользователю
// @Summary Добавление контакта
// @Description Добавление контакта пользователю по email или login
// @Tags contacts
// @Accept json
// @Produce json
// @Param email_or_login query string true "email или login пользователя"
// @Success 200 {array}  models.Contacts
// @Failure 400 {object} models.Response
// @Failure 404 {object} models.Response
// @Router /contacts [post]
func (c *Contacts) AddContact(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	reqContact := models.RequestContact{}
	if err := json.NewDecoder(r.Body).Decode(&reqContact); err != nil {
		c.errorLog.Printf("AddContact, method Decode: %s\n", err)

		_ = render.Render(w, r, models.BadRequest(models.ErrJsonRequestContact))
		return
	}

	userTo, err := c.dbPostgres.GetUserTdidForEmailOrLogin(ctx, reqContact.EmailOrLogin)
	if err != nil {
		if err.Error() == models.ErrNoResultLogin {
			err = models.ErrFindContact
		}
		c.errorLog.Printf("AddContact, method GetUserTdidForEmailOrLogin: %s\n", err)

		_ = render.Render(w, r, models.ResourceNotFound(err))
		return
	}

	tdid := ctx.Value("tdid").(string)

	err = c.dbPostgres.FindContact(ctx, tdid, userTo.Tdid)
	if err != nil {
		c.errorLog.Printf("AddContact, method FindContact: %s\n", err)

		_ = render.Render(w, r, models.BadRequest(err))
		return
	}

	now := time.Now().Format(time.RFC3339)
	date, _ := time.Parse(time.RFC3339, now)

	contact := &models.Contact{
		UserFrom:   tdid,
		UserTo:     userTo.Tdid,
		UserToName: userTo.Name,
		CreatedAt:  date,
		UpdatedAt:  date,
	}

	err = c.dbPostgres.AddContact(ctx, contact)
	if err != nil {
		c.errorLog.Printf("AddContact, method AddContact: %s\n", err)

		_ = render.Render(w, r, models.BadRequest(err))
		return
	}

	contacts, err := c.dbPostgres.GetContacts(ctx, tdid)
	if err != nil {
		c.errorLog.Printf("AddContact, method GetContacts: %s\n", err)

		_ = render.Render(w, r, models.BadRequest(err))
		return
	}

	files := models.AvatarProto(contacts)

	images, err := c.f3Client.LoadFiles(ctx, files)
	if err != nil {
		c.errorLog.Printf("AddContact, method LoadFiles: %s\n", err)

		_ = render.Render(w, r, models.BadRequest(err))
		return
	}

	models.ContactsAvatar(contacts, images)

	render.JSON(w, r, contacts)
}

// AddContacts Добавление контактов пользователю
// @Summary Добавление контактов
// @Description Добавление контактов пользователю по номеру телефона
// @Tags contacts
// @Accept json
// @Produce json
// @Param email_or_login query string true "номеру телефона пользователя"
// @Success 200 {array}  models.Contacts
// @Failure 400 {object} models.Response
// @Failure 404 {object} models.Response
// @Router /contacts [put]
func (c *Contacts) AddContacts(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	contacts := models.RequestContactsPhones{}
	if err := json.NewDecoder(r.Body).Decode(&contacts); err != nil {
		c.errorLog.Printf("AddContacts, method Decode: %s\n", err)

		_ = render.Render(w, r, models.BadRequest(models.ErrJsonRequestContact))
		return
	}

	tdid := ctx.Value("tdid").(string)

	phones := models.GetPhones(contacts)

	users, err := c.dbPostgres.GetUsersByPhones(ctx, phones)
	if err != nil {
		if err.Error() == models.ErrNoResultLogin {
			err = models.ErrFindContact
		}
		c.errorLog.Printf("AddContacts, method GetUsersByPhones: %s\n", err)

		_ = render.Render(w, r, models.ResourceNotFound(err))
		return
	}

	tdids := models.GetTdids(users)

	contactsIsset, err := c.dbPostgres.FindContacts(ctx, tdid, tdids)
	if err != nil {
		c.errorLog.Printf("AddContacts, method FindContacts: %s\n", err)

		_ = render.Render(w, r, models.BadRequest(err))
		return
	}

	contactsIn := models.GetContacts(contactsIsset, users, tdid, contacts)
	if len(contactsIn) == 0 {
		c.errorLog.Printf("AddContacts, method GetContacts: %s\n", err)

		_ = render.Render(w, r, models.BadRequest(models.ErrorNewContactsPhone))
		return
	}

	err = c.dbPostgres.AddContacts(ctx, contactsIn)
	if err != nil {
		c.errorLog.Printf("AddContacts, method AddContacts: %s\n", err)

		_ = render.Render(w, r, models.BadRequest(err))
		return
	}

	contactsReturn, err := c.dbPostgres.GetContacts(ctx, tdid)
	if err != nil {
		c.errorLog.Printf("AddContacts, method GetContacts: %s\n", err)

		_ = render.Render(w, r, models.BadRequest(err))
		return
	}

	files := models.AvatarProto(contactsReturn)

	images, err := c.f3Client.LoadFiles(ctx, files)
	if err != nil {
		c.errorLog.Printf("AddContacts, method LoadFiles: %s\n", err)

		_ = render.Render(w, r, models.BadRequest(err))
		return
	}

	models.ContactsAvatar(contactsReturn, images)

	render.JSON(w, r, contactsReturn)
}
