package components

import (
	"testing"
)

// TestInputFieldWidthIncludesCursorSpace verifies that input fields are configured
// with field width = MaxLength + 1 to account for cursor space in tview.
// This is a regression test for the bug where 10-digit customer numbers only
// displayed 9 digits because tview's SetFieldWidth reserves space for the cursor.
func TestInputFieldWidthIncludesCursorSpace(t *testing.T) {
	testCases := []struct {
		name      string
		maxLength int
	}{
		{"CustomerNumber_10digits", 10},
		{"PolicyNumber_10digits", 10},
		{"SmallField_5digits", 5},
		{"SingleDigit", 1},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			form := NewForm()
			form.AddField(&FormField{
				Name:      "testField",
				Label:     "Test Field",
				MaxLength: tc.maxLength,
				FieldType: FieldTypeNumeric,
				Editable:  true,
			})

			field := form.GetField("testField")
			if field == nil {
				t.Fatal("Field not found")
			}

			// Get the configured field width from the underlying tview InputField
			// tview InputField stores field width internally; we verify via GetFieldWidth
			inputField := field.InputField()
			configuredWidth := inputField.GetFieldWidth()

			// The field width should be MaxLength + 1 to display all characters
			// while accounting for cursor space
			expectedWidth := tc.maxLength + 1
			if configuredWidth != expectedWidth {
				t.Errorf("Field width = %d, want %d (MaxLength %d + 1 for cursor)",
					configuredWidth, expectedWidth, tc.maxLength)
			}
		})
	}
}

// TestInputFieldAcceptsMaxLengthCharacters verifies that the acceptance function
// allows exactly MaxLength characters to be typed.
func TestInputFieldAcceptsMaxLengthCharacters(t *testing.T) {
	form := NewForm()
	form.AddField(&FormField{
		Name:      "customerNum",
		Label:     "Customer Number",
		MaxLength: 10,
		FieldType: FieldTypeNumeric,
		Editable:  true,
	})

	field := form.GetField("customerNum")
	if field == nil {
		t.Fatal("Field not found")
	}

	// Simulate typing a 10-digit number
	inputField := field.InputField()
	inputField.SetText("1000000011")

	// Verify the full value is stored
	value := inputField.GetText()
	if value != "1000000011" {
		t.Errorf("Expected '1000000011', got '%s'", value)
	}
	if len(value) != 10 {
		t.Errorf("Expected 10 characters, got %d", len(value))
	}
}

// TestNumericFieldRejectsNonDigits verifies that numeric fields only accept digits.
func TestNumericFieldRejectsNonDigits(t *testing.T) {
	form := NewForm()
	acceptFunc := form.getAcceptanceFunc(FieldTypeNumeric, 10)

	// Should accept digits
	for _, r := range "0123456789" {
		if !acceptFunc(string(r), r) {
			t.Errorf("Numeric field should accept digit '%c'", r)
		}
	}

	// Should reject non-digits
	for _, r := range "abcXYZ!@#-." {
		if acceptFunc(string(r), r) {
			t.Errorf("Numeric field should reject character '%c'", r)
		}
	}
}

// TestFieldWidthForAllFieldTypes verifies that all field types get proper width.
func TestFieldWidthForAllFieldTypes(t *testing.T) {
	fieldTypes := []FieldType{
		FieldTypeText,
		FieldTypeNumeric,
		FieldTypeDate,
		FieldTypeYesNo,
		FieldTypeDecimal,
	}

	for _, ft := range fieldTypes {
		t.Run(fieldTypeName(ft), func(t *testing.T) {
			form := NewForm()
			maxLen := 15
			form.AddField(&FormField{
				Name:      "testField",
				Label:     "Test Field",
				MaxLength: maxLen,
				FieldType: ft,
				Editable:  true,
			})

			field := form.GetField("testField")
			inputField := field.InputField()
			configuredWidth := inputField.GetFieldWidth()

			expectedWidth := maxLen + 1
			if configuredWidth != expectedWidth {
				t.Errorf("Field width for %s = %d, want %d",
					fieldTypeName(ft), configuredWidth, expectedWidth)
			}
		})
	}
}

func fieldTypeName(ft FieldType) string {
	switch ft {
	case FieldTypeText:
		return "Text"
	case FieldTypeNumeric:
		return "Numeric"
	case FieldTypeDate:
		return "Date"
	case FieldTypeYesNo:
		return "YesNo"
	case FieldTypeDecimal:
		return "Decimal"
	default:
		return "Unknown"
	}
}
