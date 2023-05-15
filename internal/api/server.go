package api

import (
	"log"
	"time"

	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
)

type api struct {
	server *fasthttp.Server
	router fasthttp.RequestHandler
	config *Config

	image   Service
	captcha Generator
}

type Config struct {
	Host         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

func New(config *Config, image Service, captcha Generator) *api {
	api := &api{
		config:  config,
		image:   image,
		captcha: captcha,
	}

	r := router.New()

	r.GET("/captcha", api.generate)
	r.POST("/api/v1/captcha/solve", api.solve)

	r.ServeFiles("/{filepath:*}", "templates/ui/dist")

	api.router = r.Handler

	return api
}

func (api *api) ListenAndServe() error {
	api.server = &fasthttp.Server{
		Handler:         api.router,
		ReadTimeout:     api.config.ReadTimeout,
		WriteTimeout:    api.config.WriteTimeout,
		IdleTimeout:     api.config.IdleTimeout,
		CloseOnShutdown: true,
	}

	log.Println(api.config.Host)

	return api.server.ListenAndServe(api.config.Host)
}

func (api *api) Shutdown() error {
	time.Sleep(1 * time.Second)
	return api.server.Shutdown()
}
