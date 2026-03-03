# Investigation: BMS Field Start Position Bug

## Bug Summary

Fields on screens do not respect their start position defined in the BMS file. The screen layout shows fields at incorrect column positions compared to what is specified in the BMS `POS=(row,column)` attributes.

**Example from screenshot:**
- The label "Cust Number" and its input field should start at specific BMS positions
- In the BMS file (SSMAPC1): Label at `POS=(04,30)`, Input at `POS=(04,50)`
- But the current rendering ignores these positions entirely

## Root Cause Analysis

### The Problem

The BMS position information (`POS=(row,column)`) is **completely ignored** during screen rendering. The current implementation:

1. **Does NOT parse** position data from the BMS file
2. **Does NOT store** position information in the `FormField` structure
3. **Does NOT use** positions when rendering fields on screen
4. Fields are rendered **sequentially** in the order they are defined, not at their BMS-specified positions

### Affected Components

| File | Issue |
|------|-------|
| `internal/ui/components/form.go:28-43` | `FormField` struct has no Row/Column fields |
| `internal/ui/components/screen.go:128-137` | `SetForm()` renders fields sequentially, ignoring position |
| `internal/ui/views/*.go` | View files define fields without position data |

### Code Evidence

**1. FormField struct missing position fields (`form.go:28-43`):**
```go
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
    // NO Row, Column, or Position fields!
}
```

**2. SetForm renders fields sequentially (`screen.go:128-137`):**
```go
for _, field := range form.Fields() {
    row := tview.NewFlex().SetDirection(tview.FlexColumn)
    row.AddItem(field.LabelView(), 16, 0, false)           // Fixed 16-char label width
    row.AddItem(field.InputField(), field.MaxLength+2, 0, true)  // Variable input width
    spacer := tview.NewBox()
    spacer.SetBackgroundColor(tcell.ColorDefault)
    row.AddItem(spacer, 0, 1, false)
    s.formArea.AddItem(row, 1, 0, true)  // Added to vertical FlexBox, ignoring BMS row
}
```

**3. BMS defines precise positions (`ssmap.bms:27-31`):**
```
        DFHMDF POS=(04,30),LENGTH=12,ATTRB=(NORM,ASKIP),               X
               INITIAL='Cust Number '
ENT1CNO DFHMDF POS=(04,50),LENGTH=10,ATTRB=(NORM,UNPROT,IC,FSET),      *
               JUSTIFY=(RIGHT,ZERO)
```

### Current Rendering Pipeline

```
BMS File (ssmap.bms)
    ↓ [NOT PARSED - positions ignored]
FormField Definition (in Go code)
    ↓ [No position info stored]
Screen.SetForm()
    ↓ [Sequential vertical layout]
Fields stacked in definition order
    ↓
WRONG POSITIONS!
```

### Expected vs Actual Layout

**BMS Specification (Customer Screen - SSMAPC1):**
| Field | Label Position | Input Position |
|-------|---------------|----------------|
| Cust Number | Row 4, Col 30 | Row 4, Col 50 |
| First Name | Row 5, Col 30 | Row 5, Col 50 |
| Last Name | Row 6, Col 30 | Row 6, Col 50 |
| DOB | Row 7, Col 30 | Row 7, Col 50 |
| House Name | Row 8, Col 30 | Row 8, Col 50 |
| House Number | Row 9, Col 30 | Row 9, Col 50 |
| Postcode | Row 10, Col 30 | Row 10, Col 50 |
| Phone: Home | Row 11, Col 30 | Row 11, Col 50 |
| Phone: Mob | Row 12, Col 30 | Row 12, Col 50 |
| Email Addr | Row 13, Col 30 | Row 13, Col 50 |

**Current Behavior:**
- All labels use fixed 16-character width
- All inputs start immediately after their label
- Fields are stacked vertically starting from form area row 3
- Column positions (30, 50) are completely ignored

## Proposed Solution

### Option 1: Add Position Fields to FormField (Recommended)

1. **Extend FormField struct** to include Row and Column:
   ```go
   type FormField struct {
       // ... existing fields ...
       Row    int  // BMS row position (1-24)
       Column int  // BMS column position for input (1-80)
       LabelColumn int  // BMS column position for label
   }
   ```

2. **Update view files** to specify positions when creating fields:
   ```go
   v.form.AddField(&components.FormField{
       Label:       "Cust Number",
       Name:        "customer_num",
       Row:         4,
       LabelColumn: 30,
       Column:      50,
       // ...
   })
   ```

3. **Modify SetForm() in screen.go** to use absolute positioning:
   - Use tview's Grid or manual positioning based on Row/Column
   - Place labels at their LabelColumn position
   - Place inputs at their Column position

### Option 2: Grid-Based Layout

Use tview's Grid layout to position fields at exact row/column coordinates matching the BMS 24x80 terminal format.

### Implementation Notes

- The current grid setup (`screen.go:45-48`) already has 24 rows
- Column widths need adjustment to support character-precise positioning
- Consider using a single 80-column layout to match BMS exactly
- Labels and inputs should be positioned independently based on their BMS coordinates

## Edge Cases and Considerations

1. **Varying label lengths**: BMS fields have different label lengths, but the current implementation uses fixed 16-char width
2. **Input field positioning**: Inputs should start exactly at their POS column, not after the label
3. **Screen boundaries**: Ensure fields don't overflow the 80-column boundary
4. **All screens affected**: This issue affects all view files (customer, motor, house, endowment, etc.)

## Test Strategy

1. Compare rendered output against BMS position specifications
2. Verify each field appears at its correct (row, column) position
3. Test tab order respects field row ordering
4. Validate on standard 80x24 terminal dimensions

## Implementation Notes

The fix requires changes to:
1. `internal/ui/components/form.go` - Add position fields to FormField
2. `internal/ui/components/screen.go` - Update SetForm() to use position-based layout
3. `internal/ui/views/customer.go` - Add position data to field definitions
4. `internal/ui/views/motor.go` - Add position data to field definitions
5. `internal/ui/views/house.go` - Add position data to field definitions
6. `internal/ui/views/endowment.go` - Add position data to field definitions
7. `internal/ui/views/policy_placeholders.go` - Add position data to field definitions

Note: There are pre-existing compilation errors in the codebase (SetOnSubmit, CustomerService, PolicyService undefined). These should be addressed as part of the fix or in a separate task.

---

## Implementation Complete

### Changes Made

**1. `internal/ui/components/form.go`** - Added position fields to FormField struct:
```go
type FormField struct {
    // ... existing fields ...
    Row         int // Row position (1-24), 0 means auto-layout
    LabelColumn int // Column position for label (1-80), 0 means auto-layout
    Column      int // Column position for input field (1-80), 0 means auto-layout
    // ...
}
```

**2. `internal/ui/components/screen.go`** - Refactored SetForm() to support position-based layout:
- Added `setFormWithPositions()` for BMS-compliant positioning
- Added `setFormSequential()` for legacy fallback behavior
- Added `sortFieldsByColumn()` helper to order fields by column position
- Added `getFieldStartCol()` helper to determine field start column
- Fields are now grouped by row and positioned based on their BMS coordinates

**3. View files updated with BMS positions:**
- `internal/ui/views/customer.go` - SSMAPC1 positions (rows 4-13, labels at col 30, inputs at col 50)
- `internal/ui/views/motor.go` - SSMAPP1 positions (rows 4-16, labels at col 30, inputs at col 50)
- `internal/ui/views/house.go` - SSMAPP3 positions (rows 4-13, labels at col 30, inputs at col 50)
- `internal/ui/views/endowment.go` - SSMAPP2 positions (rows 4-14, labels at col 30, inputs at col 50)
- `internal/ui/views/policy_placeholders.go` - SSMAPP4/SSMAPP5 positions for Commercial and Claim views

### Test Results

```
$ go build ./...
# Successful - no errors

$ go test ./...
ok  	github.com/cicsdev/genapp/internal/repository	0.541s
ok  	github.com/cicsdev/genapp/internal/service	0.939s
ok  	github.com/cicsdev/genapp/internal/ui/components	0.672s
```

All tests pass. The implementation is backward-compatible - views without position data will continue to use the sequential layout.

### Key Implementation Details

1. **Automatic detection**: The `SetForm()` method automatically detects if any field has position data and switches between position-based and sequential layout modes.

2. **Row grouping**: Fields are grouped by their BMS row number, allowing multiple fields on the same row (though current screens only have one field per row).

3. **Column spacing**: Labels are positioned at their `LabelColumn` position, with the label width calculated as the difference between `Column` and `LabelColumn`. Input fields start exactly at their `Column` position.

4. **Backward compatibility**: Views without position data (Row=0, Column=0) fall back to the original sequential layout behavior.
