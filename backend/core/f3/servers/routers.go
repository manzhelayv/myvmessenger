package servers

import (
	"f3/models"
	http2 "f3/servers/http"
	middleware2 "f3/servers/middleware"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	//"models"
	"net/http"
	"os"
	"time"
)

func (srv *HttpServer) setupRouter() chi.Router {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	allowedHeaders := []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token", "Content-Language"}
	allowedHeaders = append(allowedHeaders, models.AllowedHeaders()...)

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   allowedHeaders,
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	r.Use(middleware.Timeout(60 * time.Second))
	_, err := os.Getwd()
	if err != nil {
		os.Exit(1)
	}

	srv.Routes(r)

	return r
}

func (srv HttpServer) Routes(r *chi.Mux) http.Handler {
	r.Mount("/", http2.NewUsers(srv.redisCli, srv.InfoLog, srv.ErrorLog).Routes())

	r.Group(func(r chi.Router) {
		r.Use(middleware2.NewUserAccessCtx(models.JwtSecretKey).ChiMiddleware)
	})

	return r
}
