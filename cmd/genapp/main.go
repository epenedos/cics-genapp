// Package main is the entry point for the GENAPP application.
// This application is a Go + OpenTUI migration of the original CICS/COBOL
// General Insurance Policy Management System.
package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/cicsdev/genapp/internal/config"
	"github.com/cicsdev/genapp/internal/repository"
	"github.com/cicsdev/genapp/internal/service"
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

// exitCode is used to track the exit status
var exitCode = 0

func main() {
	// Parse command line flags
	showVersion := flag.Bool("version", false, "Show version information")
	noDb := flag.Bool("no-db", false, "Run without database connection (demo mode)")
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

	// Initialize services (nil if running without database)
	var services *ui.Services
	var db *repository.DB

	if !*noDb {
		// Initialize database connection
		db, err = initDatabase(cfg)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Database connection failed: %v\n", err)
			fmt.Fprintf(os.Stderr, "Hint: Use --no-db flag to run in demo mode without database\n")
			os.Exit(1)
		}

		// Create services with repositories
		services = createServices(db)

		// Initialize counters on startup
		if err := initializeCounters(services.Counter); err != nil {
			log.Warn().Err(err).Msg("Failed to initialize counters (non-fatal)")
		}
	}

	// Create UI application with services
	app := ui.NewApp(services)

	// Create and register views
	customerView := views.NewCustomerView()
	motorView := views.NewMotorPolicyView()
	endowmentView := views.NewEndowmentPolicyView()
	houseView := views.NewHousePolicyView()
	commercialView := views.NewCommercialPolicyView()
	claimView := views.NewClaimView()

	// Wire up services to all views
	if services != nil {
		customerView.SetServices(services.Customer, services.Policy, services.Counter)
		motorView.SetServices(services.Customer, services.Policy, services.Counter)
		endowmentView.SetServices(services.Customer, services.Policy, services.Counter)
		houseView.SetServices(services.Customer, services.Policy, services.Counter)
		commercialView.SetServices(services.Customer, services.Policy, services.Counter)
		claimView.SetServices(services.Customer, services.Policy, services.Counter)
	}

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

	// Set up graceful shutdown handling
	shutdownChan := make(chan os.Signal, 1)
	signal.Notify(shutdownChan, os.Interrupt, syscall.SIGTERM)

	// Set cleanup handler
	app.SetOnQuit(func() {
		log.Info().Msg("Application shutting down")
		if db != nil {
			if err := db.Close(); err != nil {
				log.Error().Err(err).Msg("Error closing database connection")
			} else {
				log.Info().Msg("Database connection closed")
			}
		}
	})

	// Handle SIGTERM/SIGINT in a goroutine
	go func() {
		<-shutdownChan
		log.Info().Msg("Received shutdown signal")
		app.Stop()
	}()

	// Run the application
	if err := app.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Application error: %v\n", err)
		exitCode = 1
	}

	os.Exit(exitCode)
}

// initDatabase establishes the database connection with retry logic.
func initDatabase(cfg *config.Config) (*repository.DB, error) {
	dbConfig := repository.DBConfig{
		Host:            cfg.Database.Host,
		Port:            cfg.Database.Port,
		User:            cfg.Database.User,
		Password:        cfg.Database.Password,
		Database:        cfg.Database.DBName,
		SSLMode:         cfg.Database.SSLMode,
		MaxOpenConns:    25,
		MaxIdleConns:    5,
		ConnMaxLifetime: 5 * time.Minute,
	}

	// Attempt connection with retry
	var db *repository.DB
	var err error
	maxRetries := 3
	for i := 0; i < maxRetries; i++ {
		db, err = repository.NewDB(dbConfig)
		if err == nil {
			break
		}
		if i < maxRetries-1 {
			log.Warn().Err(err).Int("attempt", i+1).Msg("Database connection failed, retrying...")
			time.Sleep(time.Second * time.Duration(i+1))
		}
	}

	if err != nil {
		return nil, fmt.Errorf("failed to connect after %d attempts: %w", maxRetries, err)
	}

	// Verify connection is alive
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.Ping(ctx); err != nil {
		db.Close()
		return nil, fmt.Errorf("database ping failed: %w", err)
	}

	log.Info().
		Str("host", cfg.Database.Host).
		Int("port", cfg.Database.Port).
		Str("database", cfg.Database.DBName).
		Msg("Database connection established")

	return db, nil
}

// createServices initializes all service instances with their dependencies.
func createServices(db *repository.DB) *ui.Services {
	// Create repositories
	customerRepo := repository.NewCustomerRepository(db)
	policyRepo := repository.NewPolicyRepository(db)
	counterRepo := repository.NewCounterRepository(db)

	// Create services with repositories
	customerSvc := service.NewCustomerService(customerRepo, counterRepo)
	policySvc := service.NewPolicyService(policyRepo, customerRepo, counterRepo)
	counterSvc := service.NewCounterService(counterRepo)

	return &ui.Services{
		Customer: customerSvc,
		Policy:   policySvc,
		Counter:  counterSvc,
	}
}

// initializeCounters ensures all application counters exist in the database.
func initializeCounters(counterSvc *service.CounterService) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return counterSvc.InitializeCounters(ctx)
}
