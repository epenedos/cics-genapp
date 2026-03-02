package components

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// Screen provides a base layout matching the BMS 24x80 terminal format.
// Layout structure:
// - Row 1: Screen ID (col 1-4) + Title (col 12+)
// - Rows 4-7: Menu options (col 8-24)
// - Rows 4-18: Form fields (col 30+)
// - Row 22: Option selection prompt
// - Row 24: Error/status message
type Screen struct {
	// Layout components
	grid       *tview.Grid
	screenID   *tview.TextView
	title      *tview.TextView
	menuArea   *tview.Flex
	formArea   *tview.Flex
	optionArea *tview.Flex
	errorArea  *tview.TextView

	// Components
	menu *Menu
	form *Form

	// Key handling
	onEnter  func()
	onEscape func()

	// Application reference for focus management
	app *tview.Application
}

// NewScreen creates a new screen with the standard BMS layout.
func NewScreen(screenID, title string) *Screen {
	s := &Screen{}

	// Create the main grid with 24 rows and flexible columns
	s.grid = tview.NewGrid().
		SetRows(1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1). // 24 rows
		SetColumns(7, 4, 18, 16, 20, -1). // Column layout to match BMS positions
		SetBorders(false)

	// Row 1: Screen ID (bold, position 1)
	s.screenID = tview.NewTextView().
		SetText(screenID).
		SetTextAlign(tview.AlignLeft).
		SetTextColor(tcell.ColorWhite)
	s.screenID.SetBackgroundColor(tcell.ColorDefault)
	s.screenID.SetTextStyle(tcell.StyleDefault.Bold(true))

	// Row 1: Title (bold, position 12)
	s.title = tview.NewTextView().
		SetText(title).
		SetTextAlign(tview.AlignLeft).
		SetTextColor(tcell.ColorWhite)
	s.title.SetBackgroundColor(tcell.ColorDefault)
	s.title.SetTextStyle(tcell.StyleDefault.Bold(true))

	// Menu area (rows 4-7, col 8-24)
	s.menuArea = tview.NewFlex().SetDirection(tview.FlexRow)
	s.menuArea.SetBackgroundColor(tcell.ColorDefault)

	// Form area (rows 4-18, col 30+)
	s.formArea = tview.NewFlex().SetDirection(tview.FlexRow)
	s.formArea.SetBackgroundColor(tcell.ColorDefault)

	// Option selection area (row 22)
	s.optionArea = tview.NewFlex().SetDirection(tview.FlexColumn)
	s.optionArea.SetBackgroundColor(tcell.ColorDefault)

	// Error area (row 24)
	s.errorArea = tview.NewTextView().
		SetTextAlign(tview.AlignLeft).
		SetTextColor(tcell.ColorRed)
	s.errorArea.SetBackgroundColor(tcell.ColorDefault)

	// Assemble the grid
	// Row 0 (line 1): Screen ID and Title
	s.grid.AddItem(s.screenID, 0, 0, 1, 2, 0, 0, false)
	s.grid.AddItem(s.title, 0, 2, 1, 4, 0, 0, false)

	// Rows 3-6 (lines 4-7): Menu area (left side)
	s.grid.AddItem(s.menuArea, 3, 0, 4, 2, 0, 0, false)

	// Rows 3-17 (lines 4-18): Form area (right side)
	s.grid.AddItem(s.formArea, 3, 2, 15, 4, 0, 0, true)

	// Row 21 (line 22): Option selection
	s.grid.AddItem(s.optionArea, 21, 0, 1, 6, 0, 0, false)

	// Row 23 (line 24): Error message
	s.grid.AddItem(s.errorArea, 23, 0, 1, 6, 0, 0, false)

	return s
}

// SetMenu sets the menu component for this screen.
func (s *Screen) SetMenu(menu *Menu) *Screen {
	s.menu = menu
	s.menuArea.Clear()
	s.menuArea.AddItem(menu.OptionsDisplay(), 0, 1, false)

	// Add option prompt and input to option area
	s.optionArea.Clear()
	promptLabel := tview.NewTextView().
		SetText("Select Option ").
		SetTextAlign(tview.AlignLeft)
	promptLabel.SetBackgroundColor(tcell.ColorDefault)
	s.optionArea.AddItem(promptLabel, 14, 0, false)
	s.optionArea.AddItem(menu.OptionInput(), 2, 0, true)

	return s
}

// SetForm sets the form component for this screen.
func (s *Screen) SetForm(form *Form) *Screen {
	s.form = form
	form.SetErrorView(s.errorArea)
	s.formArea.Clear()

	// Add each field as a row with label and input
	for _, field := range form.Fields() {
		row := tview.NewFlex().SetDirection(tview.FlexColumn)
		row.AddItem(field.LabelView(), 16, 0, false)
		row.AddItem(field.InputField(), field.MaxLength+2, 0, true)
		// Add spacer
		spacer := tview.NewBox()
		spacer.SetBackgroundColor(tcell.ColorDefault)
		row.AddItem(spacer, 0, 1, false)
		s.formArea.AddItem(row, 1, 0, true)
	}

	return s
}

// SetOnEnter sets the callback for Enter key press.
func (s *Screen) SetOnEnter(handler func()) *Screen {
	s.onEnter = handler
	return s
}

// SetOnEscape sets the callback for Escape key press.
func (s *Screen) SetOnEscape(handler func()) *Screen {
	s.onEscape = handler
	return s
}

// Layout returns the screen's grid layout as a tview Primitive.
func (s *Screen) Layout() tview.Primitive {
	return s.grid
}

// Menu returns the screen's menu component.
func (s *Screen) Menu() *Menu {
	return s.menu
}

// Form returns the screen's form component.
func (s *Screen) Form() *Form {
	return s.form
}

// ShowError displays an error message in the error area.
func (s *Screen) ShowError(msg string) {
	s.errorArea.SetText(msg)
	s.errorArea.SetTextColor(tcell.ColorRed)
}

// ShowSuccess displays a success message in the error area.
func (s *Screen) ShowSuccess(msg string) {
	s.errorArea.SetText(msg)
	s.errorArea.SetTextColor(tcell.ColorGreen)
}

// ClearError clears the error/message area.
func (s *Screen) ClearError() {
	s.errorArea.SetText("")
}

// Clear resets both the form and menu.
func (s *Screen) Clear() {
	if s.form != nil {
		s.form.Clear()
	}
	if s.menu != nil {
		s.menu.Clear()
	}
	s.ClearError()
}

// SetFocus sets focus to the first input field.
func (s *Screen) SetFocus(app *tview.Application) {
	s.app = app
	if s.form != nil {
		s.form.SetFocus(app)
	} else if s.menu != nil {
		app.SetFocus(s.menu.OptionInput())
	}
}

// HandleKey handles screen-specific key events.
func (s *Screen) HandleKey(event *tcell.EventKey) *tcell.EventKey {
	switch event.Key() {
	case tcell.KeyEnter:
		if s.onEnter != nil {
			s.onEnter()
			return nil
		}
	case tcell.KeyTab:
		// Move to next field
		if s.form != nil && s.app != nil {
			s.form.NextField(s.app)
			return nil
		}
	case tcell.KeyBacktab:
		// Move to previous field
		if s.form != nil && s.app != nil {
			s.form.PrevField(s.app)
			return nil
		}
	}
	return event
}

// ErrorArea returns the error text view for direct manipulation.
func (s *Screen) ErrorArea() *tview.TextView {
	return s.errorArea
}

// Grid returns the underlying grid for advanced customization.
func (s *Screen) Grid() *tview.Grid {
	return s.grid
}

// FormArea returns the form area flex for adding custom content.
func (s *Screen) FormArea() *tview.Flex {
	return s.formArea
}

// MenuArea returns the menu area flex for adding custom content.
func (s *Screen) MenuArea() *tview.Flex {
	return s.menuArea
}

// OptionArea returns the option selection area for customization.
func (s *Screen) OptionArea() *tview.Flex {
	return s.optionArea
}
