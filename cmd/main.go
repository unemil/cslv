package main

import (
	"cslv/internal/api"
	"cslv/internal/generator"
	"cslv/internal/service/image"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	api := api.New(&api.Config{
		Host:         ":80",
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		IdleTimeout:  10 * time.Second,
	}, image.New(), generator.New())

	go func() {
		if err := api.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}()
	defer api.Shutdown()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop
}
