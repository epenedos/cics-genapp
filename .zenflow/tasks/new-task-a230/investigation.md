# Bug Investigation: Screen Layout Issue

## Bug Summary
The screen layout displays correctly with the right columns (menu items, labels, input fields), but the visual positioning does not match the expected BMS 24x80 terminal layout. The user reports "the code now shows the right columns but the screen layout looks the same."

## Evidence from Screenshot
The SSC1 "General Insurance Customer Menu" screen shows:
- Menu options on the left: "1. Cust Inquiry", "2. Cust Add", "4. Cust Update"
- Form fields on the right: Customer Number, First Name, Last Name, DOB, etc.
- Labels and values appear to be rendering, but the layout spacing may not match the original BMS design

## Root Cause Analysis

### Expected BMS Layout (from `base/src/ssmap.bms`)
The original SSMAPC1 BMS definition specifies:
- Screen ID at position (1,1)
- Title at position (1,12)
- Menu options at column 8 (rows 4-7)
- Labels at column 30 (starting row 4)
- Input fields at column 50 (starting row 4)
- Select Option prompt at (22,8)
- Error field at (24,8)

### Current Implementation
The `screen.go` uses a tview Grid with fixed column widths:
```go
SetColumns(7, 4, 18, 16, 20, -1)  // 6 columns
```

Grid item placement:
- Row 0, columns 0-1: Screen ID (spans 2 columns = 11 chars)
- Row 0, columns 2-5: Title (spans 4 columns)
- Rows 3-6, columns 0-1: Menu area (spans 2 columns)
- Rows 3-17, columns 2-5: Form area (spans 4 columns)
- Row 21, columns 0-5: Option selection area
- Row 23, columns 0-5: Error area

### Identified Issues

1. **Grid Column Widths vs BMS Columns**: The grid columns (7, 4, 18, 16, 20, flex) total approximately 65+ fixed characters, but the mapping to BMS 80-column positions is approximate, not exact.

2. **Form Area Position Calculations**: In `setFormWithPositions()`, the code assumes:
   ```go
   const formAreaStartCol = 30
   ```
   This creates relative positioning within the form area Flex container, but the Flex container itself doesn't start at BMS column 30 - it starts at grid column 2.

3. **Label Width Calculation**:
   ```go
   labelWidth := inputCol - labelCol  // 50 - 30 = 20
   ```
   This gives a 20-character label width, but the actual label text varies (e.g., "Cust Number" = 11 chars, "Email Addr" = 10 chars).

4. **No Absolute Positioning**: The tview Grid and Flex layouts use relative positioning and proportional sizing, not absolute character coordinates like BMS terminals.

## Affected Components
- `internal/ui/components/screen.go` - Grid layout and `setFormWithPositions()` function
- `internal/ui/views/customer.go` - Field position definitions
- All other view files that use BMS positioning (motor.go, house.go, endowment.go, policy_placeholders.go)

## Proposed Solution

The fix needs to ensure the visual layout matches BMS specifications. Options:

### Option 1: Adjust Grid Column Widths
Recalculate grid column widths to match BMS positions:
- Columns 1-7: Screen ID area (7 chars)
- Columns 8-29: Menu area (22 chars)
- Columns 30-79: Form area (50 chars)

New grid columns: `SetColumns(7, 22, 50, -1)` or similar.

### Option 2: Use Single-Row Layout per Field
Instead of relying on column calculations, render each field row with explicit character widths:
- Label: fixed width (e.g., 20 chars, left-aligned)
- Input: actual MaxLength + padding

### Option 3: Use tview.TextView with Fixed-Width Font
Create a single TextView for the entire screen and render content at exact character positions, mimicking a true 80x24 terminal.

## Recommended Approach
Option 1 (Adjust Grid Column Widths) is the least invasive fix. The grid columns should be recalculated to provide proper spacing between menu area and form area, matching the BMS column 30 start position for labels.

Current columns sum: 7 + 4 + 18 + 16 + 20 = 65 fixed + flex
- First 2 columns (0-1): 7 + 4 = 11 chars for menu/screen ID
- Need: 7 chars for screen ID area, ~22 chars for menu (cols 8-29), ~50 chars for form (cols 30-79)

New proposed columns: `SetColumns(7, 1, 22, 20, 30, -1)`
- Column 0 (7): Screen ID + left padding
- Column 1 (1): Spacer
- Column 2 (22): Menu area (columns 8-29 in BMS)
- Column 3 (20): Label area (columns 30-49 in BMS)
- Column 4 (30): Input area (columns 50-79 in BMS)
- Column 5 (flex): Right margin

## Implementation Notes
1. Update `NewScreen()` in `screen.go` to use corrected column widths
2. Simplify `setFormWithPositions()` to use the grid structure directly
3. Verify menu items display fully without truncation
4. Test with all screen types (customer, motor, house, endowment, commercial, claims)

## Test Plan
1. Run the application and verify SSC1 screen matches BMS layout
2. Verify all menu options are fully visible
3. Verify labels align at column 30 equivalent
4. Verify input fields align at column 50 equivalent
5. Verify tab navigation works correctly between fields
6. Test all policy screens (SSP1, SSP2, SSP3, SSP4, SSP5)

---

## Implementation Notes (Completed)

### Changes Made to `internal/ui/components/screen.go`

**1. Updated Grid Column Widths**

Changed from:
```go
SetColumns(7, 4, 18, 16, 20, -1)  // 6 columns
```

To:
```go
SetColumns(7, 22, 50, -1)  // 4 columns matching BMS positions
```

This correctly maps to BMS 80-column layout:
- Column 0 (7 chars): Screen ID area (BMS cols 1-7)
- Column 1 (22 chars): Menu area (BMS cols 8-29)
- Column 2 (50 chars): Form area (BMS cols 30-79)
- Column 3 (flex): Right margin

**2. Updated Grid Item Placements**

Changed item placements to use new column structure:
- Screen ID: row 0, col 0, spans 1 column
- Title: row 0, cols 1-3
- Menu area: rows 3-6, col 1 only (not overlapping with form)
- Form area: rows 3-17, cols 2-3
- Option/Error areas: span all 4 columns

**3. Fixed Form Area Position Calculations**

Updated `setFormWithPositions()` to correctly offset BMS positions:
- Form area container starts at BMS column 30
- Field positions within form area are now calculated relative to position 0
- BMS column positions are converted by subtracting the offset (30)

### Test Results
- `go build ./...` - Passed
- `go test ./...` - All tests pass:
  - `internal/repository` - OK
  - `internal/service` - OK
  - `internal/ui/components` - OK
