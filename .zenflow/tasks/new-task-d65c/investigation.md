# Bug Investigation: Input Field Displays Only 9 Digits Instead of 10

## Bug Summary

When a user types a 10-digit customer number (e.g., `1000000011`), the input field only displays 9 digits (`000000011`). The first digit disappears from the visible display, though the full value may still be stored internally.

**Expected**: Customer number field displays all 10 digits (e.g., `1000000011`)
**Actual**: Only 9 digits are visible (e.g., `000000011`)

## Root Cause Analysis

### Investigation Steps

1. **Checked field configuration**: The `MaxLength` property is correctly set to 10 in `internal/ui/views/customer.go:43`

2. **Checked acceptance function**: The validation logic in `internal/ui/components/form.go:96-116` is correct:
   ```go
   if len(text) > maxLen {
       return false
   }
   ```
   This correctly allows up to 10 characters.

3. **Checked tview InputField behavior**: The issue is in how tview's `SetFieldWidth()` handles display width.

### Root Cause

The bug is in `internal/ui/components/form.go:79`:

```go
field.inputView = tview.NewInputField().
    SetFieldWidth(field.MaxLength).  // <-- Sets display width to 10
    SetAcceptanceFunc(f.getAcceptanceFunc(field.FieldType, field.MaxLength))
```

**The problem**: tview's `InputField.SetFieldWidth(n)` reserves space for the cursor, effectively displaying only `n-1` characters when the field is full. When you set `SetFieldWidth(10)`, the field can only display 9 characters visually, even though 10 characters are stored.

From tview's source code (`inputfield.go`), the Draw function calls:
```go
i.textArea.setMinCursorPadding(fieldWidth-1, 1)
```

This `fieldWidth-1` causes the off-by-one display issue.

## Affected Components

| File | Line | Description |
|------|------|-------------|
| `internal/ui/components/form.go` | 79 | `SetFieldWidth(field.MaxLength)` - needs +1 adjustment |

All form fields using `FieldTypeNumeric` with exact-length requirements are affected:
- Customer Number (10 digits) - `internal/ui/views/customer.go:40-52`
- Policy Number (10 digits) - multiple policy views
- Other numeric fields that need to display their full MaxLength

## Proposed Solution

**Fix**: Add 1 to the field width when creating the InputField to account for cursor space:

```go
// In internal/ui/components/form.go, line 79
field.inputView = tview.NewInputField().
    SetFieldWidth(field.MaxLength + 1).  // +1 to account for cursor space
    SetAcceptanceFunc(f.getAcceptanceFunc(field.FieldType, field.MaxLength))
```

This change ensures that:
1. The field displays all `MaxLength` characters
2. The acceptance function still correctly limits input to `MaxLength` characters
3. No changes needed to validation or data storage logic

### Alternative Considered

Another option would be to increase `MaxLength` in each field definition, but this would:
- Require changes to multiple files
- Allow users to type more characters than intended
- Break validation at the service layer

The proposed solution is cleaner as it only requires one change in one location.

## Test Plan

1. Start the application and navigate to customer screen
2. Type a 10-digit customer number: `1000000011`
3. Verify all 10 digits are visible in the input field
4. Verify the customer inquiry returns the correct customer
5. Test other numeric fields (policy numbers) to ensure they also display correctly
