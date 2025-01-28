package http

import (
	redisCli "f3/servers/redis"
	"github.com/go-chi/chi"
	"log"
	"net/http"
)

type Users struct {
	redisCli *redisCli.Redis
	infoLog  *log.Logger
	errorLog *log.Logger
}

func NewUsers(redisCli *redisCli.Redis, infoLog *log.Logger, errorLog *log.Logger) *Users {
	return &Users{
		redisCli: redisCli,
		infoLog:  infoLog,
		errorLog: errorLog,
	}
}

func (m Users) Routes() http.Handler {
	r := chi.NewRouter()

	//r.Post("/user", m.Register)

	return r
}
