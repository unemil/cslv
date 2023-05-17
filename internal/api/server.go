package api

import (
	"time"

	"github.com/fasthttp/router"
	"github.com/rs/zerolog/log"
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

	v1 := r.Group("/api/v1")
	{
		v1.POST("/captcha/solve", api.solve)
	}

	r.ServeFiles("/{filepath:*}", "templates/ui/dist")

	r.ServeFiles("/swagger/{filepath:*}", "templates/swagger-ui/dist")
	r.ServeFiles("/docs/{filepath:*}", "templates/swagger-ui")

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

	log.Info().Msgf("Starting http server on %s ...", api.config.Host)

	return api.server.ListenAndServe(api.config.Host)
}

func (api *api) Shutdown() error {
	log.Info().Msg("Stopping http server ...")
	time.Sleep(1 * time.Second)

	return api.server.Shutdown()
}
