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

	captcha Service
}

type Config struct {
	Host         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

func New(config *Config, captcha Service) *api {
	api := &api{
		config:  config,
		captcha: captcha,
	}

	r := router.New()
	v1 := r.Group("/api/v1")
	{
		captcha := v1.Group("/captcha")
		{
			captcha.GET("/generate", api.generate)
			captcha.POST("/solve", api.solve)
			captcha.GET("/analyze", api.analyze)
		}
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
