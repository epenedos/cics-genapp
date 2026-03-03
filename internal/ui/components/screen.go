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

	// Focus tracking: true when menu option input has focus
	focusOnMenu bool
}

// NewScreen creates a new screen with the standard BMS layout.
func NewScreen(screenID, title string) *Screen {
	s := &Screen{}

	// Create the main grid with 24 rows and columns matching BMS 80-column positions
	// BMS layout: col 1-7 = screen ID, col 8-29 = menu area, col 30-79 = form area
	s.grid = tview.NewGrid().
		SetRows(1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1). // 24 rows
		SetColumns(7, 22, 50, -1). // col 0: screen ID (7), col 1: menu (22), col 2: form (50), col 3: flex
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
	// Row 0 (line 1): Screen ID at col 0, Title at col 1-3 (starts at BMS col 8/12)
	s.grid.AddItem(s.screenID, 0, 0, 1, 1, 0, 0, false)
	s.grid.AddItem(s.title, 0, 1, 1, 3, 0, 0, false)

	// Rows 3-6 (lines 4-7): Menu area at col 1 (BMS col 8-29)
	s.grid.AddItem(s.menuArea, 3, 1, 4, 1, 0, 0, false)

	// Rows 3-17 (lines 4-18): Form area at col 2 (BMS col 30-79)
	s.grid.AddItem(s.formArea, 3, 2, 15, 2, 0, 0, true)

	// Row 21 (line 22): Option selection spans all columns
	s.grid.AddItem(s.optionArea, 21, 0, 1, 4, 0, 0, false)

	// Row 23 (line 24): Error message spans all columns
	s.grid.AddItem(s.errorArea, 23, 0, 1, 4, 0, 0, false)

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

	// Check if any field has position data
	hasPositionData := false
	for _, field := range form.Fields() {
		if field.Row > 0 && field.Column > 0 {
			hasPositionData = true
			break
		}
	}

	if hasPositionData {
		// Use position-based layout matching BMS coordinates
		s.setFormWithPositions(form)
	} else {
		// Fall back to legacy sequential layout
		s.setFormSequential(form)
	}

	return s
}

// setFormSequential renders fields in sequential order (legacy behavior).
func (s *Screen) setFormSequential(form *Form) {
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
}

// setFormWithPositions renders fields at their BMS-specified positions.
// BMS uses 1-indexed row/column positions on a 24x80 terminal.
// The form area starts at screen row 4 (BMS row 4), so we offset accordingly.
func (s *Screen) setFormWithPositions(form *Form) {
	// Group fields by their row position
	rowFields := make(map[int][]*FormField)
	minRow := 999
	maxRow := 0

	for _, field := range form.Fields() {
		if field.Row > 0 {
			rowFields[field.Row] = append(rowFields[field.Row], field)
			if field.Row < minRow {
				minRow = field.Row
			}
			if field.Row > maxRow {
				maxRow = field.Row
			}
		}
	}

	// Create rows from minRow to maxRow
	for row := minRow; row <= maxRow; row++ {
		fields := rowFields[row]
		if len(fields) == 0 {
			// Empty row - add a spacer
			spacer := tview.NewBox()
			spacer.SetBackgroundColor(tcell.ColorDefault)
			s.formArea.AddItem(spacer, 1, 0, false)
			continue
		}

		// Create a flex row for this line
		rowFlex := tview.NewFlex().SetDirection(tview.FlexColumn)

		// Sort fields by their column position
		sortedFields := sortFieldsByColumn(fields)

		// Form area starts at BMS column 30, so we offset BMS positions accordingly
		// A field at BMS column 30 should be at form area position 0
		const formAreaOffset = 30

		currentCol := 0 // Position within the form area (starts at 0)

		for _, field := range sortedFields {
			labelCol := field.LabelColumn
			if labelCol == 0 {
				labelCol = formAreaOffset // Default to column 30
			}
			inputCol := field.Column
			if inputCol == 0 {
				inputCol = labelCol + 20 // Default offset
			}

			// Convert BMS columns to form area positions (subtract offset)
			labelPosInForm := labelCol - formAreaOffset
			if labelPosInForm < 0 {
				labelPosInForm = 0
			}

			// Add spacer before label if needed
			if labelPosInForm > currentCol {
				spacerWidth := labelPosInForm - currentCol
				spacer := tview.NewBox()
				spacer.SetBackgroundColor(tcell.ColorDefault)
				rowFlex.AddItem(spacer, spacerWidth, 0, false)
				currentCol = labelPosInForm
			}

			// Calculate label width (from labelCol to inputCol)
			labelWidth := inputCol - labelCol
			if labelWidth < 1 {
				labelWidth = 16 // Minimum label width
			}

			// Add label
			rowFlex.AddItem(field.LabelView(), labelWidth, 0, false)
			currentCol += labelWidth

			// Add input field
			inputWidth := field.MaxLength + 2 // +2 for input field padding
			rowFlex.AddItem(field.InputField(), inputWidth, 0, true)
			currentCol += inputWidth
		}

		// Add trailing spacer to fill the rest of the row
		spacer := tview.NewBox()
		spacer.SetBackgroundColor(tcell.ColorDefault)
		rowFlex.AddItem(spacer, 0, 1, false)

		s.formArea.AddItem(rowFlex, 1, 0, true)
	}
}

// sortFieldsByColumn sorts fields by their LabelColumn or Column position.
func sortFieldsByColumn(fields []*FormField) []*FormField {
	// Simple insertion sort since we typically have few fields per row
	result := make([]*FormField, len(fields))
	copy(result, fields)

	for i := 1; i < len(result); i++ {
		j := i
		for j > 0 && getFieldStartCol(result[j-1]) > getFieldStartCol(result[j]) {
			result[j-1], result[j] = result[j], result[j-1]
			j--
		}
	}
	return result
}

// getFieldStartCol returns the starting column for a field (label column or input column).
func getFieldStartCol(f *FormField) int {
	if f.LabelColumn > 0 {
		return f.LabelColumn
	}
	if f.Column > 0 {
		return f.Column
	}
	return 30 // Default BMS form area start
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
		// Move to next field, including menu option input in the cycle
		if s.app != nil {
			if s.focusOnMenu {
				// Currently on menu option input, go to first form field
				if s.form != nil {
					s.form.FocusFirstField(s.app)
					s.focusOnMenu = false
				}
				return nil
			}
			if s.form != nil {
				// Check if we're at the last form field
				if s.form.IsAtLastEditableField() && s.menu != nil {
					// Move to menu option input
					s.app.SetFocus(s.menu.OptionInput())
					s.focusOnMenu = true
					return nil
				}
				// Otherwise, move to next form field
				s.form.NextField(s.app)
				return nil
			}
		}
	case tcell.KeyBacktab:
		// Move to previous field, including menu option input in the cycle
		if s.app != nil {
			if s.focusOnMenu {
				// Currently on menu option input, go to last form field
				if s.form != nil {
					s.form.FocusLastField(s.app)
					s.focusOnMenu = false
				}
				return nil
			}
			if s.form != nil {
				// Check if we're at the first form field
				if s.form.IsAtFirstEditableField() && s.menu != nil {
					// Move to menu option input
					s.app.SetFocus(s.menu.OptionInput())
					s.focusOnMenu = true
					return nil
				}
				// Otherwise, move to previous form field
				s.form.PrevField(s.app)
				return nil
			}
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
