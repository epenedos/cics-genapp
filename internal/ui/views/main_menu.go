package views

import (
	"github.com/cicsdev/genapp/internal/ui/components"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// MainMenuView is the main menu screen for navigating to different modules.
type MainMenuView struct {
	*BaseView
	menu          *components.Menu
	onNavigate    func(screen string)
}

// NewMainMenuView creates a new main menu view.
func NewMainMenuView() *MainMenuView {
	v := &MainMenuView{
		BaseView: NewBaseView("main-menu", "MENU", "General Insurance Application - Main Menu"),
	}

	// Create the main menu
	v.menu = components.NewMenu()
	v.menu.AddOption("1", "Customer Menu", true)
	v.menu.AddOption("2", "Motor Policy", true)
	v.menu.AddOption("3", "Endowment Policy", true)
	v.menu.AddOption("4", "House Policy", true)
	v.menu.AddOption("5", "Commercial Policy", true)
	v.menu.AddOption("6", "Claims", true)
	v.menu.AddOption("9", "Exit", true)

	v.screen.SetMenu(v.menu)

	// Set up Enter key handler
	v.SetOnSubmit(v.handleSubmit)

	return v
}

// SetOnNavigate sets the callback for screen navigation.
func (v *MainMenuView) SetOnNavigate(fn func(screen string)) {
	v.onNavigate = fn
}

// handleSubmit processes the menu selection.
func (v *MainMenuView) handleSubmit() {
	option, valid := v.menu.ProcessSelection()
	if !valid {
		v.ShowError("Invalid option selected")
		return
	}

	v.ClearError()

	if v.onNavigate != nil {
		switch option.Key {
		case "1":
			v.onNavigate("customer")
		case "2":
			v.onNavigate("motor")
		case "3":
			v.onNavigate("endowment")
		case "4":
			v.onNavigate("house")
		case "5":
			v.onNavigate("commercial")
		case "6":
			v.onNavigate("claim")
		case "9":
			v.onNavigate("exit")
		}
	}
}

// HandleKey handles key events specific to the main menu.
func (v *MainMenuView) HandleKey(event *tcell.EventKey) *tcell.EventKey {
	// Handle numeric key shortcuts
	if event.Key() == tcell.KeyRune {
		switch event.Rune() {
		case '1', '2', '3', '4', '5', '6', '9':
			v.menu.SetSelectedOption(string(event.Rune()))
			v.handleSubmit()
			return nil
		}
	}

	return v.BaseView.HandleKey(event)
}

// SetFocus sets focus to the menu option input.
func (v *MainMenuView) SetFocus(app *tview.Application) {
	v.app = app
	app.SetFocus(v.menu.OptionInput())
}
