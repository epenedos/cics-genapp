// Package main is the entry point for the GENAPP application.
// This application is a Go + OpenTUI migration of the original CICS/COBOL
// General Insurance Policy Management System.
package main

import (
	"os"

	"github.com/cicsdev/genapp/internal/config"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	// Set up zerolog
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to load configuration")
	}

	// Set log level from config
	level, err := zerolog.ParseLevel(cfg.LogLevel)
	if err != nil {
		level = zerolog.InfoLevel
	}
	zerolog.SetGlobalLevel(level)

	log.Info().
		Str("app", "genapp").
		Str("version", "1.0.0").
		Msg("Starting General Insurance Application")

	// TODO: Initialize database connection
	// TODO: Create service instances
	// TODO: Start OpenTUI application

	log.Info().Msg("Application initialized successfully")
}
