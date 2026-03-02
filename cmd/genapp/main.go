// Package main is the entry point for the GENAPP application.
// This application is a Go + OpenTUI migration of the original CICS/COBOL
// General Insurance Policy Management System.
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/cicsdev/genapp/internal/config"
	"github.com/cicsdev/genapp/internal/ui"
	"github.com/cicsdev/genapp/internal/ui/views"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Version information
const (
	AppName    = "genapp"
	AppVersion = "1.0.0"
)

func main() {
	// Parse command line flags
	showVersion := flag.Bool("version", false, "Show version information")
	flag.Parse()

	if *showVersion {
		fmt.Printf("%s version %s\n", AppName, AppVersion)
		os.Exit(0)
	}

	// Set up zerolog - disable for TUI mode as it interferes with the terminal
	// In production, logs should go to a file
	zerolog.SetGlobalLevel(zerolog.Disabled)

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load configuration: %v\n", err)
		os.Exit(1)
	}

	// Set log level from config (if logging to file is enabled later)
	level, err := zerolog.ParseLevel(cfg.LogLevel)
	if err != nil {
		level = zerolog.InfoLevel
	}
	_ = level // Will be used when file logging is enabled

	log.Info().
		Str("app", AppName).
		Str("version", AppVersion).
		Msg("Starting General Insurance Application")

	// Create UI application
	// Note: Services are nil for now - they will be wired up in the Integration step
	app := ui.NewApp(nil)

	// Create and register views
	customerView := views.NewCustomerView()
	motorView := views.NewMotorPolicyView()
	endowmentView := views.NewEndowmentPolicyView()
	houseView := views.NewHousePolicyView()
	commercialView := views.NewCommercialPolicyView()
	claimView := views.NewClaimView()

	// Set up navigation callbacks
	navigateTo := func(screen string) {
		switch screen {
		case "customer":
			app.SwitchTo(ui.ScreenCustomer)
		case "motor":
			app.SwitchTo(ui.ScreenMotor)
		case "endowment":
			app.SwitchTo(ui.ScreenEndowment)
		case "house":
			app.SwitchTo(ui.ScreenHouse)
		case "commercial":
			app.SwitchTo(ui.ScreenCommercial)
		case "claim":
			app.SwitchTo(ui.ScreenClaim)
		case "exit":
			app.Stop()
		}
	}

	customerView.SetOnNavigate(navigateTo)
	motorView.SetOnNavigate(navigateTo)
	endowmentView.SetOnNavigate(navigateTo)
	houseView.SetOnNavigate(navigateTo)
	commercialView.SetOnNavigate(navigateTo)
	claimView.SetOnNavigate(navigateTo)

	// Register views with the application
	app.RegisterView(ui.ScreenCustomer, customerView)
	app.RegisterView(ui.ScreenMotor, motorView)
	app.RegisterView(ui.ScreenEndowment, endowmentView)
	app.RegisterView(ui.ScreenHouse, houseView)
	app.RegisterView(ui.ScreenCommercial, commercialView)
	app.RegisterView(ui.ScreenClaim, claimView)

	// Set cleanup handler
	app.SetOnQuit(func() {
		log.Info().Msg("Application shutting down")
	})

	// Run the application
	if err := app.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Application error: %v\n", err)
		os.Exit(1)
	}
}
