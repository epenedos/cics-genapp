package components

import (
	"testing"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// TestScreenTabNavigationIncludesMenuOption verifies that TAB navigation cycles
// through form fields AND the menu option input field.
// This is a regression test for the bug where users could not select menu options
// on screens with forms because TAB only cycled through form fields.
func TestScreenTabNavigationIncludesMenuOption(t *testing.T) {
	app := tview.NewApplication()

	// Create a screen with both a form and a menu
	screen := NewScreen("TEST", "Test Screen")

	// Create a form with two editable fields
	form := NewForm()
	form.AddField(&FormField{
		Name:      "field1",
		Label:     "Field 1",
		MaxLength: 10,
		Editable:  true,
	})
	form.AddField(&FormField{
		Name:      "field2",
		Label:     "Field 2",
		MaxLength: 10,
		Editable:  true,
	})

	// Create a menu
	menu := NewMenu()
	menu.AddOption("1", "Option 1", true)
	menu.AddOption("2", "Option 2", true)

	screen.SetForm(form)
	screen.SetMenu(menu)
	screen.SetFocus(app)

	// Initially, focus should be on the first form field
	if form.focusIndex != 0 {
		t.Errorf("Expected initial focus on field 0, got %d", form.focusIndex)
	}

	// Simulate TAB key press - should move to field 2
	tabEvent := tcell.NewEventKey(tcell.KeyTab, 0, tcell.ModNone)
	screen.HandleKey(tabEvent)

	if form.focusIndex != 1 {
		t.Errorf("After first TAB, expected focus on field 1, got %d", form.focusIndex)
	}

	// Simulate TAB key press - should move to menu option input
	screen.HandleKey(tabEvent)

	if !screen.focusOnMenu {
		t.Error("After second TAB (from last form field), expected focus on menu option input")
	}

	// Simulate TAB key press - should wrap back to first form field
	screen.HandleKey(tabEvent)

	if screen.focusOnMenu {
		t.Error("After TAB from menu, expected focus back on form")
	}
	if form.focusIndex != 0 {
		t.Errorf("After TAB from menu, expected focus on field 0, got %d", form.focusIndex)
	}
}

// TestScreenBackTabNavigationIncludesMenuOption verifies that Shift+TAB (BackTab)
// navigation cycles through form fields AND the menu option input field in reverse.
func TestScreenBackTabNavigationIncludesMenuOption(t *testing.T) {
	app := tview.NewApplication()

	// Create a screen with both a form and a menu
	screen := NewScreen("TEST", "Test Screen")

	// Create a form with two editable fields
	form := NewForm()
	form.AddField(&FormField{
		Name:      "field1",
		Label:     "Field 1",
		MaxLength: 10,
		Editable:  true,
	})
	form.AddField(&FormField{
		Name:      "field2",
		Label:     "Field 2",
		MaxLength: 10,
		Editable:  true,
	})

	// Create a menu
	menu := NewMenu()
	menu.AddOption("1", "Option 1", true)

	screen.SetForm(form)
	screen.SetMenu(menu)
	screen.SetFocus(app)

	// Initially, focus should be on the first form field
	if form.focusIndex != 0 {
		t.Errorf("Expected initial focus on field 0, got %d", form.focusIndex)
	}

	// Simulate BackTab key press - should move to menu option input (from first field)
	backTabEvent := tcell.NewEventKey(tcell.KeyBacktab, 0, tcell.ModNone)
	screen.HandleKey(backTabEvent)

	if !screen.focusOnMenu {
		t.Error("After BackTab from first form field, expected focus on menu option input")
	}

	// Simulate BackTab key press - should move to last form field
	screen.HandleKey(backTabEvent)

	if screen.focusOnMenu {
		t.Error("After BackTab from menu, expected focus back on form")
	}
	if form.focusIndex != 1 {
		t.Errorf("After BackTab from menu, expected focus on last field (1), got %d", form.focusIndex)
	}
}

// TestScreenTabNavigationNoMenu verifies TAB works normally when screen has no menu.
func TestScreenTabNavigationNoMenu(t *testing.T) {
	app := tview.NewApplication()

	// Create a screen with only a form (no menu)
	screen := NewScreen("TEST", "Test Screen")

	// Create a form with two editable fields
	form := NewForm()
	form.AddField(&FormField{
		Name:      "field1",
		Label:     "Field 1",
		MaxLength: 10,
		Editable:  true,
	})
	form.AddField(&FormField{
		Name:      "field2",
		Label:     "Field 2",
		MaxLength: 10,
		Editable:  true,
	})

	screen.SetForm(form)
	// Note: no menu set
	screen.SetFocus(app)

	// Simulate TAB key presses - should cycle through form fields only
	tabEvent := tcell.NewEventKey(tcell.KeyTab, 0, tcell.ModNone)

	screen.HandleKey(tabEvent)
	if form.focusIndex != 1 {
		t.Errorf("Expected focus on field 1 after TAB, got %d", form.focusIndex)
	}

	// TAB again should wrap to field 0 (no menu to go to)
	screen.HandleKey(tabEvent)
	if form.focusIndex != 0 {
		t.Errorf("Expected focus to wrap to field 0, got %d", form.focusIndex)
	}
}

// TestScreenTabNavigationNoForm verifies TAB works normally when screen has only menu.
func TestScreenTabNavigationNoForm(t *testing.T) {
	app := tview.NewApplication()

	// Create a screen with only a menu (no form)
	screen := NewScreen("TEST", "Test Screen")

	// Create a menu
	menu := NewMenu()
	menu.AddOption("1", "Option 1", true)

	screen.SetMenu(menu)
	// Note: no form set
	screen.SetFocus(app)

	// TAB should be a no-op (focus stays on menu)
	tabEvent := tcell.NewEventKey(tcell.KeyTab, 0, tcell.ModNone)
	result := screen.HandleKey(tabEvent)

	// The event should be returned (not consumed) since there's nothing to tab to
	if result == nil {
		t.Error("Expected TAB event to be returned (not consumed) when only menu exists")
	}
}

// TestFormIsAtLastEditableField verifies the IsAtLastEditableField method.
func TestFormIsAtLastEditableField(t *testing.T) {
	form := NewForm()
	form.AddField(&FormField{
		Name:      "field1",
		Label:     "Field 1",
		MaxLength: 10,
		Editable:  true,
	})
	form.AddField(&FormField{
		Name:      "field2",
		Label:     "Field 2",
		MaxLength: 10,
		Editable:  false, // non-editable
	})
	form.AddField(&FormField{
		Name:      "field3",
		Label:     "Field 3",
		MaxLength: 10,
		Editable:  true,
	})

	// Focus on first field
	form.focusIndex = 0
	if form.IsAtLastEditableField() {
		t.Error("Field 0 should not be the last editable field")
	}

	// Focus on last editable field
	form.focusIndex = 2
	if !form.IsAtLastEditableField() {
		t.Error("Field 2 should be the last editable field")
	}
}

// TestFormIsAtFirstEditableField verifies the IsAtFirstEditableField method.
func TestFormIsAtFirstEditableField(t *testing.T) {
	form := NewForm()
	form.AddField(&FormField{
		Name:      "field1",
		Label:     "Field 1",
		MaxLength: 10,
		Editable:  false, // non-editable
	})
	form.AddField(&FormField{
		Name:      "field2",
		Label:     "Field 2",
		MaxLength: 10,
		Editable:  true,
	})
	form.AddField(&FormField{
		Name:      "field3",
		Label:     "Field 3",
		MaxLength: 10,
		Editable:  true,
	})

	// Focus on first editable field (index 1)
	form.focusIndex = 1
	if !form.IsAtFirstEditableField() {
		t.Error("Field 1 should be the first editable field")
	}

	// Focus on last field
	form.focusIndex = 2
	if form.IsAtFirstEditableField() {
		t.Error("Field 2 should not be the first editable field")
	}
}
