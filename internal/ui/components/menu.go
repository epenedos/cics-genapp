package components

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// MenuOption represents a single option in a menu.
type MenuOption struct {
	Key         string // e.g., "1", "2", etc.
	Label       string // Display text
	Description string // Optional description
	Enabled     bool   // Whether the option is available
}

// Menu is a component for displaying and selecting menu options.
// It mimics the BMS menu style with numbered options.
type Menu struct {
	options         []MenuOption
	selectedIndex   int
	optionInput     *tview.InputField
	optionsDisplay  *tview.TextView
	onSelect        func(option MenuOption)
	selectedOption  string
}

// NewMenu creates a new menu component.
func NewMenu() *Menu {
	m := &Menu{
		options:       make([]MenuOption, 0),
		selectedIndex: -1,
	}

	// Create the options display
	m.optionsDisplay = tview.NewTextView().
		SetDynamicColors(true).
		SetTextAlign(tview.AlignLeft)
	m.optionsDisplay.SetBackgroundColor(tcell.ColorDefault)

	// Create the option input field (single character)
	m.optionInput = tview.NewInputField().
		SetFieldWidth(1).
		SetAcceptanceFunc(func(text string, lastChar rune) bool {
			// Only accept single digit
			return len(text) <= 1 && lastChar >= '1' && lastChar <= '9'
		})
	m.optionInput.SetFieldBackgroundColor(tcell.ColorBlack)
	m.optionInput.SetFieldTextColor(tcell.ColorGreen)

	return m
}

// AddOption adds an option to the menu.
func (m *Menu) AddOption(key, label string, enabled bool) *Menu {
	m.options = append(m.options, MenuOption{
		Key:     key,
		Label:   label,
		Enabled: enabled,
	})
	m.updateDisplay()
	return m
}

// SetOptions sets all menu options at once.
func (m *Menu) SetOptions(options []MenuOption) *Menu {
	m.options = options
	m.updateDisplay()
	return m
}

// updateDisplay refreshes the menu display text.
func (m *Menu) updateDisplay() {
	text := ""
	for _, opt := range m.options {
		if opt.Enabled {
			text += fmt.Sprintf("%s. %s\n", opt.Key, opt.Label)
		} else {
			// Show disabled options in a different color
			text += fmt.Sprintf("[gray]%s. %s[-]\n", opt.Key, opt.Label)
		}
	}
	m.optionsDisplay.SetText(text)
}

// SetOnSelect sets the callback when an option is selected.
func (m *Menu) SetOnSelect(handler func(option MenuOption)) *Menu {
	m.onSelect = handler
	return m
}

// GetSelectedOption returns the currently selected option key.
func (m *Menu) GetSelectedOption() string {
	return m.optionInput.GetText()
}

// SetSelectedOption sets the selected option.
func (m *Menu) SetSelectedOption(key string) {
	m.optionInput.SetText(key)
	m.selectedOption = key
}

// Clear clears the selection.
func (m *Menu) Clear() {
	m.optionInput.SetText("")
	m.selectedOption = ""
}

// ProcessSelection validates and processes the current selection.
func (m *Menu) ProcessSelection() (MenuOption, bool) {
	selected := m.optionInput.GetText()
	for _, opt := range m.options {
		if opt.Key == selected && opt.Enabled {
			m.selectedOption = selected
			if m.onSelect != nil {
				m.onSelect(opt)
			}
			return opt, true
		}
	}
	return MenuOption{}, false
}

// OptionsDisplay returns the text view showing the menu options.
func (m *Menu) OptionsDisplay() *tview.TextView {
	return m.optionsDisplay
}

// OptionInput returns the input field for option selection.
func (m *Menu) OptionInput() *tview.InputField {
	return m.optionInput
}

// IsValidSelection checks if the current input is a valid option.
func (m *Menu) IsValidSelection() bool {
	selected := m.optionInput.GetText()
	for _, opt := range m.options {
		if opt.Key == selected && opt.Enabled {
			return true
		}
	}
	return false
}

// OperationType represents the type of CRUD operation.
type OperationType int

const (
	// OpInquiry is a read/query operation
	OpInquiry OperationType = iota
	// OpAdd is a create operation
	OpAdd
	// OpDelete is a delete operation
	OpDelete
	// OpUpdate is an update operation
	OpUpdate
)

// String returns the string representation of the operation type.
func (o OperationType) String() string {
	switch o {
	case OpInquiry:
		return "Inquiry"
	case OpAdd:
		return "Add"
	case OpDelete:
		return "Delete"
	case OpUpdate:
		return "Update"
	default:
		return "Unknown"
	}
}

// GetOperationType converts an option key to an OperationType.
func GetOperationType(optionKey string) OperationType {
	switch optionKey {
	case "1":
		return OpInquiry
	case "2":
		return OpAdd
	case "3":
		return OpDelete
	case "4":
		return OpUpdate
	default:
		return OpInquiry
	}
}

// StandardCRUDMenu creates a standard CRUD menu with Inquiry, Add, Delete, Update options.
func StandardCRUDMenu(entityName string) *Menu {
	menu := NewMenu()
	menu.AddOption("1", entityName+" Inquiry", true)
	menu.AddOption("2", entityName+" Add", true)
	menu.AddOption("3", entityName+" Delete", true)
	menu.AddOption("4", entityName+" Update", true)
	return menu
}

// StandardInquiryAddMenu creates a menu with only Inquiry and Add options.
func StandardInquiryAddMenu(entityName string) *Menu {
	menu := NewMenu()
	menu.AddOption("1", entityName+" Inquiry", true)
	menu.AddOption("2", entityName+" Add", true)
	return menu
}

// CustomerMenu creates the standard customer menu matching SSMAPC1.
func CustomerMenu() *Menu {
	menu := NewMenu()
	menu.AddOption("1", "Cust Inquiry", true)
	menu.AddOption("2", "Cust Add", true)
	menu.AddOption("3", "", false) // Reserved
	menu.AddOption("4", "Cust Update", true)
	return menu
}

// PolicyMenu creates the standard policy menu matching SSMAPP1/P2/P3.
func PolicyMenu() *Menu {
	menu := NewMenu()
	menu.AddOption("1", "Policy Inquiry", true)
	menu.AddOption("2", "Policy Add", true)
	menu.AddOption("3", "Policy Delete", true)
	menu.AddOption("4", "Policy Update", true)
	return menu
}

// CommercialPolicyMenu creates the commercial policy menu (no update option).
func CommercialPolicyMenu() *Menu {
	menu := NewMenu()
	menu.AddOption("1", "Policy Inquiry", true)
	menu.AddOption("2", "Policy Add", true)
	menu.AddOption("3", "Policy Delete", true)
	return menu
}

// ClaimMenu creates the claim menu matching SSMAPP5.
func ClaimMenu() *Menu {
	menu := NewMenu()
	menu.AddOption("1", "Claim Inquiry", true)
	menu.AddOption("2", "Claim Add", true)
	return menu
}
