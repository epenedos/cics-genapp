package components

import (
	"strings"
	"testing"
)

// TestMenuDisplaySkipsDisabledEmptyLabels verifies that disabled options with
// empty labels are not displayed. This is a regression test for the bug where
// the Customer screen only showed options 1 and 2, when option 4 should also
// be visible (option 3 being a reserved/placeholder slot with empty label).
func TestMenuDisplaySkipsDisabledEmptyLabels(t *testing.T) {
	menu := NewMenu()
	menu.AddOption("1", "Cust Inquiry", true)
	menu.AddOption("2", "Cust Add", true)
	menu.AddOption("3", "", false) // Reserved/placeholder - should NOT be displayed
	menu.AddOption("4", "Cust Update", true)

	displayText := menu.OptionsDisplay().GetText(false)

	// Verify option 1 is displayed
	if !strings.Contains(displayText, "1. Cust Inquiry") {
		t.Error("Option 1 'Cust Inquiry' should be displayed")
	}

	// Verify option 2 is displayed
	if !strings.Contains(displayText, "2. Cust Add") {
		t.Error("Option 2 'Cust Add' should be displayed")
	}

	// Verify option 3 (disabled with empty label) is NOT displayed
	if strings.Contains(displayText, "3. ") {
		t.Error("Option 3 (disabled with empty label) should NOT be displayed")
	}

	// Verify option 4 is displayed
	if !strings.Contains(displayText, "4. Cust Update") {
		t.Error("Option 4 'Cust Update' should be displayed")
	}
}

// TestMenuDisplayShowsDisabledOptionsWithLabels verifies that disabled options
// WITH labels are still displayed (in gray).
func TestMenuDisplayShowsDisabledOptionsWithLabels(t *testing.T) {
	menu := NewMenu()
	menu.AddOption("1", "Active Option", true)
	menu.AddOption("2", "Disabled Option", false) // Disabled but has a label

	displayText := menu.OptionsDisplay().GetText(false)

	// Verify option 1 is displayed
	if !strings.Contains(displayText, "1. Active Option") {
		t.Error("Option 1 'Active Option' should be displayed")
	}

	// Verify option 2 (disabled with label) IS displayed
	if !strings.Contains(displayText, "2. Disabled Option") {
		t.Error("Option 2 'Disabled Option' should be displayed even though disabled")
	}
}

// TestCustomerMenuDisplaysCorrectOptions verifies that the CustomerMenu factory
// function creates a menu that displays options 1, 2, and 4 (not 3).
func TestCustomerMenuDisplaysCorrectOptions(t *testing.T) {
	menu := CustomerMenu()
	displayText := menu.OptionsDisplay().GetText(false)

	// Should show: 1. Cust Inquiry, 2. Cust Add, 4. Cust Update
	// Should NOT show: 3.
	expectedOptions := []string{
		"1. Cust Inquiry",
		"2. Cust Add",
		"4. Cust Update",
	}

	for _, expected := range expectedOptions {
		if !strings.Contains(displayText, expected) {
			t.Errorf("CustomerMenu should display '%s'", expected)
		}
	}

	// Option 3 should not be visible
	if strings.Contains(displayText, "3. ") {
		t.Error("CustomerMenu should NOT display option 3 (reserved placeholder)")
	}
}

// TestPolicyMenuDisplaysAllOptions verifies that the PolicyMenu shows all 4 options.
func TestPolicyMenuDisplaysAllOptions(t *testing.T) {
	menu := PolicyMenu()
	displayText := menu.OptionsDisplay().GetText(false)

	expectedOptions := []string{
		"1. Policy Inquiry",
		"2. Policy Add",
		"3. Policy Delete",
		"4. Policy Update",
	}

	for _, expected := range expectedOptions {
		if !strings.Contains(displayText, expected) {
			t.Errorf("PolicyMenu should display '%s'", expected)
		}
	}
}

// TestCommercialPolicyMenuDisplaysThreeOptions verifies Commercial menu shows options 1-3.
func TestCommercialPolicyMenuDisplaysThreeOptions(t *testing.T) {
	menu := CommercialPolicyMenu()
	displayText := menu.OptionsDisplay().GetText(false)

	expectedOptions := []string{
		"1. Policy Inquiry",
		"2. Policy Add",
		"3. Policy Delete",
	}

	for _, expected := range expectedOptions {
		if !strings.Contains(displayText, expected) {
			t.Errorf("CommercialPolicyMenu should display '%s'", expected)
		}
	}

	// Option 4 should not exist
	if strings.Contains(displayText, "4. ") {
		t.Error("CommercialPolicyMenu should NOT display option 4")
	}
}

// TestClaimMenuDisplaysTwoOptions verifies Claim menu shows options 1-2 only.
func TestClaimMenuDisplaysTwoOptions(t *testing.T) {
	menu := ClaimMenu()
	displayText := menu.OptionsDisplay().GetText(false)

	expectedOptions := []string{
		"1. Claim Inquiry",
		"2. Claim Add",
	}

	for _, expected := range expectedOptions {
		if !strings.Contains(displayText, expected) {
			t.Errorf("ClaimMenu should display '%s'", expected)
		}
	}

	// Options 3 and 4 should not exist
	if strings.Contains(displayText, "3. ") {
		t.Error("ClaimMenu should NOT display option 3")
	}
	if strings.Contains(displayText, "4. ") {
		t.Error("ClaimMenu should NOT display option 4")
	}
}
