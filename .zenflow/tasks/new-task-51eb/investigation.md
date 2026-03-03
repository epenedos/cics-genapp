# Bug Investigation: Unable to Choose Menu Options

## Bug Summary

**Reported Issue:** "On the screen i can never choose the option. Normally there is an option on each screen ie: 1 to 4 but the input is impossible"

**Severity:** HIGH - Blocks core application functionality on screens with forms

**Affected Screens:**
- Customer View (SSC1)
- Motor Policy View (SSP1)
- Endowment Policy View (SSP2)
- House Policy View (SSP3)
- Commercial Policy View (SSP4)
- Claims View (SSP5)

**Non-affected Screens:**
- Main Menu View - works correctly (has no form, only menu)

## Root Cause Analysis

### The Problem

Screens that have **both a form and a menu** cannot focus on the menu's option input field. The TAB key only cycles through form fields, never reaching the menu option input. Since users cannot focus on the option input field, they cannot type a menu option (1, 2, 3, or 4) to perform operations.

### Technical Details

1. **Focus Initialization** (`internal/ui/components/screen.go:195-203`):
   ```go
   func (s *Screen) SetFocus(app *tview.Application) {
       s.app = app
       if s.form != nil {
           s.form.SetFocus(app)  // Form takes priority
       } else if s.menu != nil {
           app.SetFocus(s.menu.OptionInput())  // Menu only gets focus if no form
       }
   }
   ```
   - When a screen has a form, focus always goes to the form
   - The menu option input field is never focused initially

2. **TAB Key Navigation** (`internal/ui/components/screen.go:213-224`):
   ```go
   case tcell.KeyTab:
       if s.form != nil && s.app != nil {
           s.form.NextField(s.app)  // Only cycles through form fields
           return nil
       }
   ```
   - TAB only navigates within form fields
   - The menu option input is not part of the navigation cycle

3. **Original BMS Behavior** (from `base/src/ssmap.bms`):
   - In the original CICS 3270 application, ALL unprotected fields were navigable via TAB
   - The "Select Option" field at position (22,24) was one of the TAB-able fields
   - Users could TAB from form fields to the option field and back

### Why Main Menu Works

The Main Menu (`internal/ui/views/main_menu.go`) works because:
1. It has NO form, only a menu
2. When `SetFocus` is called, focus goes directly to `menu.OptionInput()`
3. It also has a `HandleKey` override that intercepts numeric keys directly:
   ```go
   func (v *MainMenuView) HandleKey(event *tcell.EventKey) *tcell.EventKey {
       if event.Key() == tcell.KeyRune {
           switch event.Rune() {
           case '1', '2', '3', '4', '5', '6', '9':
               v.menu.SetSelectedOption(string(event.Rune()))
               v.handleSubmit()
               return nil
           }
       }
       return v.BaseView.HandleKey(event)
   }
   ```

Other views (CustomerView, MotorPolicyView, etc.) do NOT have this numeric shortcut handling because typing numbers while focused on form fields would interfere with data entry.

## Affected Components

| File | Line(s) | Issue |
|------|---------|-------|
| `internal/ui/components/screen.go` | 195-203 | `SetFocus` doesn't include menu input in focus flow |
| `internal/ui/components/screen.go` | 213-224 | TAB navigation excludes menu option input |
| `internal/ui/components/form.go` | 214-230 | `NextField`/`PrevField` only cycle through form fields |

## Proposed Solution

Modify the TAB navigation to include the menu's option input field in the navigation cycle. When the user is on the last form field and presses TAB, focus should move to the menu option input. When on the option input and pressing TAB, focus should return to the first form field.

### Implementation Approach

1. **Modify `Screen.HandleKey`** to include menu option input in TAB navigation:
   - After the last form field, TAB should focus the menu option input
   - After the menu option input, TAB should cycle back to the first form field
   - Similarly for Shift+TAB (BackTab) in reverse

2. **Track focus state** to know if focus is on form or menu

### Alternative Approach (simpler but different UX)

Add a dedicated key (e.g., F10 or similar) to toggle focus between form and menu option input. This would be less faithful to the original 3270 behavior but easier to implement.

## Test Plan

After implementing the fix:
1. Launch the application
2. Navigate to Customer View (or any policy view)
3. Verify TAB cycles through all form fields AND the menu option input
4. Verify Shift+TAB (BackTab) works in reverse
5. Verify numeric input works in both:
   - Form fields: should enter digits into the field
   - Menu option input: should enter option selection
6. Verify pressing Enter with a valid option (1, 2, or 4 for customer) triggers the appropriate action

## Related Previous Fix

**Commit 1deb777** fixed a related TAB key navigation issue where `screen.app` was not being initialized in view's `SetFocus` methods, preventing TAB from working at all. This fix corrected the initialization but did not address the menu option input being excluded from navigation.

---

## Implementation Notes

### Changes Made

**1. `internal/ui/components/screen.go`**
- Added `focusOnMenu` field to `Screen` struct to track whether focus is currently on the menu option input
- Modified `HandleKey` for `tcell.KeyTab`:
  - When focus is on menu → moves to first form field
  - When focus is on last form field → moves to menu option input
  - Otherwise → moves to next form field
- Modified `HandleKey` for `tcell.KeyBacktab`:
  - When focus is on menu → moves to last form field
  - When focus is on first form field → moves to menu option input
  - Otherwise → moves to previous form field

**2. `internal/ui/components/form.go`**
- Added `IsAtLastEditableField()` method - returns true if current focus is on the last editable field
- Added `IsAtFirstEditableField()` method - returns true if current focus is on the first editable field
- Added `FocusFirstField(app)` method - sets focus to the first editable field
- Added `FocusLastField(app)` method - sets focus to the last editable field
- Updated `NextField()` and `PrevField()` to return a bool indicating whether navigation wrapped

**3. `internal/ui/components/screen_test.go` (new file)**
- Added regression tests for TAB navigation including menu option input
- Added tests for BackTab navigation
- Added tests for screens with only form (no menu)
- Added tests for screens with only menu (no form)
- Added unit tests for `IsAtLastEditableField` and `IsAtFirstEditableField`

### Test Results

All tests pass:
```
=== RUN   TestScreenTabNavigationIncludesMenuOption
--- PASS: TestScreenTabNavigationIncludesMenuOption (0.00s)
=== RUN   TestScreenBackTabNavigationIncludesMenuOption
--- PASS: TestScreenBackTabNavigationIncludesMenuOption (0.00s)
=== RUN   TestScreenTabNavigationNoMenu
--- PASS: TestScreenTabNavigationNoMenu (0.00s)
=== RUN   TestScreenTabNavigationNoForm
--- PASS: TestScreenTabNavigationNoForm (0.00s)
=== RUN   TestFormIsAtLastEditableField
--- PASS: TestFormIsAtLastEditableField (0.00s)
=== RUN   TestFormIsAtFirstEditableField
--- PASS: TestFormIsAtFirstEditableField (0.00s)
PASS
ok      github.com/cicsdev/genapp/internal/ui/components        0.367s
```

Build succeeds with no errors.

### Behavior After Fix

1. **Customer View / Policy Views**: Users can now TAB through all form fields, and when they reach the last field, one more TAB moves focus to the "Select Option" input field where they can enter 1, 2, 3, or 4.
2. **Shift+TAB (BackTab)**: Works in reverse - from the first form field, BackTab moves to the menu option input.
3. **Main Menu**: Unchanged - continues to work as before since it has no form.
