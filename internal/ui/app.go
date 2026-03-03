// Package ui provides the terminal user interface for the GENAPP application.
// It uses tview (OpenTUI) to create 3270-style screens matching the original BMS maps.
package ui

import (
	"github.com/cicsdev/genapp/internal/service"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// ScreenType identifies the different screens in the application.
type ScreenType int

const (
	// ScreenMain is the main menu / customer screen
	ScreenMain ScreenType = iota
	// ScreenCustomer is the customer operations screen (SSMAPC1)
	ScreenCustomer
	// ScreenMotor is the motor policy screen (SSMAPP1)
	ScreenMotor
	// ScreenEndowment is the endowment policy screen (SSMAPP2)
	ScreenEndowment
	// ScreenHouse is the house policy screen (SSMAPP3)
	ScreenHouse
	// ScreenCommercial is the commercial policy screen (SSMAPP4)
	ScreenCommercial
	// ScreenClaim is the claim screen (SSMAPP5)
	ScreenClaim
)

// Terminal dimensions matching original 3270 terminal
const (
	TerminalWidth  = 80
	TerminalHeight = 24
)

// View defines the interface for all application screens.
type View interface {
	// Name returns the screen identifier
	Name() string
	// Layout returns the tview primitive for this screen
	Layout() tview.Primitive
	// SetFocus sets focus to the primary input field
	SetFocus(app *tview.Application)
	// HandleKey handles screen-specific key events
	HandleKey(event *tcell.EventKey) *tcell.EventKey
	// Clear resets all form fields
	Clear()
	// ShowError displays an error message
	ShowError(msg string)
	// ClearError clears the error message
	ClearError()
}

// Services holds all service instances needed by the UI.
type Services struct {
	Customer *service.CustomerService
	Policy   *service.PolicyService
	Counter  *service.CounterService
}

// App is the main terminal UI application.
type App struct {
	app      *tview.Application
	pages    *tview.Pages
	services *Services
	views    map[ScreenType]View
	current  ScreenType

	// Global key bindings
	onQuit func()
}

// NewApp creates a new terminal UI application.
func NewApp(services *Services) *App {
	a := &App{
		app:      tview.NewApplication(),
		pages:    tview.NewPages(),
		services: services,
		views:    make(map[ScreenType]View),
	}

	// Set up fixed terminal size (24x80)
	a.app.SetRoot(a.pages, true)

	// Set up global input capture for common key bindings
	a.app.SetInputCapture(a.handleGlobalKeys)

	return a
}

// RegisterView registers a view with the application.
func (a *App) RegisterView(screenType ScreenType, view View) {
	a.views[screenType] = view
	a.pages.AddPage(view.Name(), view.Layout(), true, false)
}

// SwitchTo switches to the specified screen.
func (a *App) SwitchTo(screenType ScreenType) {
	if view, ok := a.views[screenType]; ok {
		a.current = screenType
		a.pages.SwitchToPage(view.Name())
		view.SetFocus(a.app)
	}
}

// CurrentView returns the currently active view.
func (a *App) CurrentView() View {
	return a.views[a.current]
}

// SetOnQuit sets the callback to be invoked when the application quits.
func (a *App) SetOnQuit(fn func()) {
	a.onQuit = fn
}

// handleGlobalKeys handles application-wide key bindings.
func (a *App) handleGlobalKeys(event *tcell.EventKey) *tcell.EventKey {
	// Handle global keys
	switch event.Key() {
	case tcell.KeyEscape:
		// PF3 equivalent - typically go back or exit
		// If at main screen, quit; otherwise go back to customer screen
		if a.current == ScreenCustomer || a.current == ScreenMain {
			a.Stop()
			return nil
		}
		a.SwitchTo(ScreenCustomer)
		return nil

	case tcell.KeyCtrlC:
		// Emergency exit
		a.Stop()
		return nil

	case tcell.KeyF3:
		// PF3 - same as Escape
		if a.current == ScreenCustomer || a.current == ScreenMain {
			a.Stop()
			return nil
		}
		a.SwitchTo(ScreenCustomer)
		return nil

	case tcell.KeyF12:
		// Clear screen (3270 Master Clear)
		if view := a.CurrentView(); view != nil {
			view.Clear()
		}
		return nil
	}

	// Let the current view handle the key
	if view := a.CurrentView(); view != nil {
		return view.HandleKey(event)
	}

	return event
}

// Run starts the terminal UI application.
func (a *App) Run() error {
	// Switch to the first registered screen
	if len(a.views) > 0 {
		a.SwitchTo(ScreenCustomer)
	}

	return a.app.Run()
}

// Stop stops the terminal UI application.
func (a *App) Stop() {
	if a.onQuit != nil {
		a.onQuit()
	}
	a.app.Stop()
}

// Application returns the underlying tview application.
func (a *App) Application() *tview.Application {
	return a.app
}

// Services returns the services container.
func (a *App) Services() *Services {
	return a.services
}

// Pages returns the pages container for navigation.
func (a *App) Pages() *tview.Pages {
	return a.pages
}
