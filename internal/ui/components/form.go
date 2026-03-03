// Package components provides reusable UI components for the GENAPP application.
package components

import (
	"regexp"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// FieldType defines the type of input validation for a form field.
type FieldType int

const (
	// FieldTypeText accepts any text input
	FieldTypeText FieldType = iota
	// FieldTypeNumeric accepts only numeric input (digits)
	FieldTypeNumeric
	// FieldTypeDate accepts date input in yyyy-mm-dd format
	FieldTypeDate
	// FieldTypeYesNo accepts only Y or N
	FieldTypeYesNo
	// FieldTypeDecimal accepts decimal numbers
	FieldTypeDecimal
)

// FormField represents a single input field in a form.
type FormField struct {
	Label       string
	Name        string
	MaxLength   int
	FieldType   FieldType
	Required    bool
	Editable    bool
	InitialFocus bool
	RightJustify bool
	ZeroPad      bool

	// BMS position fields (1-indexed to match BMS coordinates)
	Row         int // Row position (1-24), 0 means auto-layout
	LabelColumn int // Column position for label (1-80), 0 means auto-layout
	Column      int // Column position for input field (1-80), 0 means auto-layout

	// Internal tview components
	labelView *tview.TextView
	inputView *tview.InputField
}

// Form is a reusable form component that mimics BMS form behavior.
type Form struct {
	flex       *tview.Flex
	fields     []*FormField
	focusIndex int
	errorView  *tview.TextView
	onSubmit   func(data map[string]string)
	onCancel   func()
}

// NewForm creates a new form component.
func NewForm() *Form {
	f := &Form{
		flex:   tview.NewFlex().SetDirection(tview.FlexRow),
		fields: make([]*FormField, 0),
	}
	return f
}

// AddField adds a new field to the form.
func (f *Form) AddField(field *FormField) *Form {
	// Create the label
	field.labelView = tview.NewTextView().
		SetText(field.Label).
		SetTextAlign(tview.AlignLeft)
	field.labelView.SetBackgroundColor(tcell.ColorDefault)

	// Create the input field
	field.inputView = tview.NewInputField().
		SetFieldWidth(field.MaxLength).
		SetAcceptanceFunc(f.getAcceptanceFunc(field.FieldType, field.MaxLength))

	// Style the input field
	field.inputView.SetFieldBackgroundColor(tcell.ColorBlack)
	field.inputView.SetFieldTextColor(tcell.ColorGreen)

	if !field.Editable {
		field.inputView.SetFieldBackgroundColor(tcell.ColorDarkBlue)
		field.inputView.SetDisabled(true)
	}

	f.fields = append(f.fields, field)
	return f
}

// getAcceptanceFunc returns the appropriate acceptance function for the field type.
func (f *Form) getAcceptanceFunc(fieldType FieldType, maxLen int) func(text string, lastChar rune) bool {
	return func(text string, lastChar rune) bool {
		if len(text) > maxLen {
			return false
		}

		switch fieldType {
		case FieldTypeNumeric:
			return lastChar >= '0' && lastChar <= '9'
		case FieldTypeDate:
			// Allow digits and hyphens for date format
			return (lastChar >= '0' && lastChar <= '9') || lastChar == '-'
		case FieldTypeYesNo:
			return lastChar == 'Y' || lastChar == 'y' || lastChar == 'N' || lastChar == 'n'
		case FieldTypeDecimal:
			return (lastChar >= '0' && lastChar <= '9') || lastChar == '.'
		default:
			return true
		}
	}
}

// SetErrorView sets the text view used to display errors.
func (f *Form) SetErrorView(errorView *tview.TextView) *Form {
	f.errorView = errorView
	return f
}

// SetOnSubmit sets the callback for form submission (Enter key).
func (f *Form) SetOnSubmit(handler func(data map[string]string)) *Form {
	f.onSubmit = handler
	return f
}

// SetOnCancel sets the callback for form cancellation (Escape key).
func (f *Form) SetOnCancel(handler func()) *Form {
	f.onCancel = handler
	return f
}

// GetField returns a field by name.
func (f *Form) GetField(name string) *FormField {
	for _, field := range f.fields {
		if field.Name == name {
			return field
		}
	}
	return nil
}

// GetValue returns the value of a field by name.
func (f *Form) GetValue(name string) string {
	if field := f.GetField(name); field != nil {
		return field.inputView.GetText()
	}
	return ""
}

// SetValue sets the value of a field by name.
func (f *Form) SetValue(name, value string) {
	if field := f.GetField(name); field != nil {
		field.inputView.SetText(value)
	}
}

// GetAllValues returns all field values as a map.
func (f *Form) GetAllValues() map[string]string {
	values := make(map[string]string)
	for _, field := range f.fields {
		values[field.Name] = field.inputView.GetText()
	}
	return values
}

// Clear resets all field values.
func (f *Form) Clear() {
	for _, field := range f.fields {
		field.inputView.SetText("")
	}
}

// ShowError displays an error message.
func (f *Form) ShowError(msg string) {
	if f.errorView != nil {
		f.errorView.SetText(msg)
		f.errorView.SetTextColor(tcell.ColorRed)
	}
}

// ClearError clears the error message.
func (f *Form) ClearError() {
	if f.errorView != nil {
		f.errorView.SetText("")
	}
}

// Validate validates all required fields and field formats.
func (f *Form) Validate() (bool, string) {
	for _, field := range f.fields {
		value := field.inputView.GetText()

		// Check required fields
		if field.Required && strings.TrimSpace(value) == "" {
			return false, field.Label + " is required"
		}

		// Validate field format
		if value != "" {
			switch field.FieldType {
			case FieldTypeDate:
				if !isValidDate(value) {
					return false, field.Label + " must be in yyyy-mm-dd format"
				}
			case FieldTypeNumeric:
				if !isNumeric(value) {
					return false, field.Label + " must be numeric"
				}
			}
		}
	}
	return true, ""
}

// NextField moves focus to the next field.
// Returns true if focus wrapped to the beginning (useful for external navigation).
func (f *Form) NextField(app *tview.Application) bool {
	if len(f.fields) == 0 {
		return true
	}

	// Find next editable field
	startIndex := f.focusIndex
	for {
		f.focusIndex = (f.focusIndex + 1) % len(f.fields)
		if f.fields[f.focusIndex].Editable || f.focusIndex == startIndex {
			break
		}
	}

	// Check if we wrapped around (went past the last editable field)
	wrapped := f.focusIndex <= startIndex && f.focusIndex != startIndex

	app.SetFocus(f.fields[f.focusIndex].inputView)
	return wrapped
}

// IsAtLastEditableField returns true if the current focus is on the last editable field.
func (f *Form) IsAtLastEditableField() bool {
	if len(f.fields) == 0 {
		return true
	}
	// Check if there's any editable field after the current focus
	for i := f.focusIndex + 1; i < len(f.fields); i++ {
		if f.fields[i].Editable {
			return false
		}
	}
	return true
}

// IsAtFirstEditableField returns true if the current focus is on the first editable field.
func (f *Form) IsAtFirstEditableField() bool {
	if len(f.fields) == 0 {
		return true
	}
	// Check if there's any editable field before the current focus
	for i := 0; i < f.focusIndex; i++ {
		if f.fields[i].Editable {
			return false
		}
	}
	return true
}

// FocusFirstField sets focus to the first editable field.
func (f *Form) FocusFirstField(app *tview.Application) {
	for i, field := range f.fields {
		if field.Editable {
			f.focusIndex = i
			app.SetFocus(field.inputView)
			return
		}
	}
}

// FocusLastField sets focus to the last editable field.
func (f *Form) FocusLastField(app *tview.Application) {
	for i := len(f.fields) - 1; i >= 0; i-- {
		if f.fields[i].Editable {
			f.focusIndex = i
			app.SetFocus(f.fields[i].inputView)
			return
		}
	}
}

// PrevField moves focus to the previous field.
// Returns true if focus wrapped to the end (useful for external navigation).
func (f *Form) PrevField(app *tview.Application) bool {
	if len(f.fields) == 0 {
		return true
	}

	// Find previous editable field
	startIndex := f.focusIndex
	for {
		f.focusIndex = f.focusIndex - 1
		if f.focusIndex < 0 {
			f.focusIndex = len(f.fields) - 1
		}
		if f.fields[f.focusIndex].Editable || f.focusIndex == startIndex {
			break
		}
	}

	// Check if we wrapped around (went before the first editable field)
	wrapped := f.focusIndex >= startIndex && f.focusIndex != startIndex

	app.SetFocus(f.fields[f.focusIndex].inputView)
	return wrapped
}

// SetFocus sets focus to the first editable field or the field marked as InitialFocus.
func (f *Form) SetFocus(app *tview.Application) {
	if len(f.fields) == 0 {
		return
	}

	// Look for field marked as InitialFocus
	for i, field := range f.fields {
		if field.InitialFocus && field.Editable {
			f.focusIndex = i
			app.SetFocus(field.inputView)
			return
		}
	}

	// Otherwise, find first editable field
	for i, field := range f.fields {
		if field.Editable {
			f.focusIndex = i
			app.SetFocus(field.inputView)
			return
		}
	}
}

// Fields returns the slice of form fields.
func (f *Form) Fields() []*FormField {
	return f.fields
}

// InputField returns the underlying tview.InputField for a form field.
func (ff *FormField) InputField() *tview.InputField {
	return ff.inputView
}

// LabelView returns the underlying tview.TextView for a form field's label.
func (ff *FormField) LabelView() *tview.TextView {
	return ff.labelView
}

// Validation helper functions

var dateRegex = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`)

func isValidDate(value string) bool {
	return dateRegex.MatchString(value)
}

func isNumeric(value string) bool {
	for _, ch := range value {
		if ch < '0' || ch > '9' {
			return false
		}
	}
	return true
}

// FormatWithLeadingZeros pads a numeric string with leading zeros.
func FormatWithLeadingZeros(value string, length int) string {
	if len(value) >= length {
		return value
	}
	return strings.Repeat("0", length-len(value)) + value
}

// FormatCustomerNum formats a customer number with leading zeros (10 digits).
func FormatCustomerNum(value string) string {
	// Remove any non-digits
	digits := ""
	for _, ch := range value {
		if ch >= '0' && ch <= '9' {
			digits += string(ch)
		}
	}
	return FormatWithLeadingZeros(digits, 10)
}

// FormatPolicyNum formats a policy number with leading zeros (10 digits).
func FormatPolicyNum(value string) string {
	return FormatCustomerNum(value) // Same format
}
