package servers

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	"net/http"
	"os"
	"server/models"
	httpResource "server/servers/http"
	middlewareServer "server/servers/middleware"
	docs "server/swagger"
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
	docs.SwaggerInfo.Version = "1"
	docs.SwaggerInfo.BasePath = "/"

	r.Mount("/swagger", httpResource.NewSwaggerResource(httpResource.BasePath, "files").Routes())

	r.Mount("/", httpResource.NewUsers(srv.dbPostgres, srv.redisCli, srv.f3GrpcClient, srv.InfoLog, srv.ErrorLog).Routes())

	r.Group(func(r chi.Router) {
		r.Use(middlewareServer.NewUserAccessCtx(models.JwtSecretKey).ChiMiddleware)

		r.Mount("/contacts", httpResource.NewContacts(srv.dbPostgres, srv.redisCli, srv.f3GrpcClient, srv.InfoLog, srv.ErrorLog).Routes())

		r.Mount("/profile", httpResource.NewProfile(srv.dbPostgres, srv.redisCli, srv.f3GrpcClient, srv.InfoLog, srv.ErrorLog).Routes())
	})

	return r
}
