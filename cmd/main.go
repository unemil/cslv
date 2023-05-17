package main

import (
	"cslv/internal/api"
	"cslv/internal/generator/captcha"
	"cslv/internal/service/image"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func init() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
}

func main() {
	api := api.New(&api.Config{
		Host:         ":80",
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		IdleTimeout:  10 * time.Second,
	}, image.New(), captcha.New())

	go func() {
		if err := api.ListenAndServe(); err != nil {
			log.Fatal().Msg(err.Error())
		}
	}()
	defer api.Shutdown()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop
}
