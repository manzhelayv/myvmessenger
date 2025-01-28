package servers

import (
	"context"
	"f3/config"
	redisCli "f3/servers/redis"
	"log"
	"net/http"
	"sync"
)

type HttpServer struct {
	handler  http.Handler
	server   *http.Server
	ctx      context.Context
	wg       *sync.WaitGroup
	redisCli *redisCli.Redis
	opt      *config.Configuration
	InfoLog  *log.Logger
	ErrorLog *log.Logger
}

func NewHTTPServer(a *application) *HttpServer {
	srv := &HttpServer{
		ctx:      a.ctx,
		wg:       a.wg,
		redisCli: &a.redisCli,
		opt:      a.opt,
		InfoLog:  a.InfoLog,
		ErrorLog: a.ErrorLog,
	}

	return srv
}

func (h *HttpServer) Start() error {
	go h.gracefulShutdown()

	h.handler = h.setupRouter()

	h.server = &http.Server{
		Addr:    h.opt.ListenAddr,
		Handler: h.handler,
	}

	log.Println("[INFO] Server HTTP started port", h.opt.ListenAddr)

	return h.server.ListenAndServe()
}

func (h *HttpServer) Stop() error {
	return h.server.Shutdown(h.ctx)
}

func (h *HttpServer) gracefulShutdown() {
	defer h.wg.Done()
	<-h.ctx.Done()
	log.Println("Shutting down HTTP server")

	grCtx, cancel := context.WithTimeout(context.Background(), 10)
	defer cancel()

	if err := h.server.Shutdown(grCtx); err != nil {
		log.Println(err, "[ERROR] HTTP server Shutdown")
	}
}
