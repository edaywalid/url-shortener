package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/edaywalid/url-shortner/internal/app"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	log.Info().Msg("Starting the application")
	app, err := app.NewApp()
	if err != nil {
		log.Error().Err(err).Msg("Failed to create the app")
		return
	}

	go func() {
		stopChan := make(chan os.Signal, 1)
		signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM)
		<-stopChan

		log.Info().Msg("Shutting down the server...")
		app.Close()
		log.Info().Msg("Server gracefully stopped")

		os.Exit(0)
	}()

	if err := app.Init(); err != nil {
		log.Error().Err(err).Msg("Failed to initialize the app")
		return
	}

}
