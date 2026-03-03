# Bug Investigation: TAB Key Not Working

## Bug Summary

The TAB key does not navigate between form fields in the terminal UI application. Users cannot use TAB to move to the next input field or Shift+TAB (Backtab) to move to the previous field.

## Root Cause Analysis

### The Problem

The TAB key handling in `internal/ui/components/screen.go` (lines 213-224) requires `s.app != nil` to function:

```go
case tcell.KeyTab:
    // Move to next field
    if s.form != nil && s.app != nil {
        s.form.NextField(s.app)
        return nil
    }
case tcell.KeyBacktab:
    // Move to previous field
    if s.form != nil && s.app != nil {
        s.form.PrevField(s.app)
        return nil
    }
```

The `s.app` field is only set in `Screen.SetFocus()`:

```go
func (s *Screen) SetFocus(app *tview.Application) {
    s.app = app  // <-- This is where s.app gets set
    if s.form != nil {
        s.form.SetFocus(app)
    } else if s.menu != nil {
        app.SetFocus(s.menu.OptionInput())
    }
}
```

### Why `s.app` Is Never Set

Multiple views override the `SetFocus()` method and **bypass** `v.screen.SetFocus(app)`:

| View | File | Line | Issue |
|------|------|------|-------|
| `CustomerView` | `views/customer.go` | 439-442 | Calls `v.form.SetFocus(app)` directly |
| `MotorPolicyView` | `views/motor.go` | 596-599 | Calls `v.form.SetFocus(app)` directly |
| `EndowmentPolicyView` | `views/endowment.go` | 543-546 | Calls `v.form.SetFocus(app)` directly |
| `HousePolicyView` | `views/house.go` | 519-522 | Calls `v.form.SetFocus(app)` directly |
| `CommercialPolicyView` | `views/policy_placeholders.go` | 133-136 | Calls `v.form.SetFocus(app)` directly |
| `ClaimView` | `views/policy_placeholders.go` | 265-268 | Calls `v.form.SetFocus(app)` directly |
| `MainMenuView` | `views/main_menu.go` | 91-94 | Sets focus directly to menu input |

Example of the problematic pattern in `CustomerView`:

```go
func (v *CustomerView) SetFocus(app *tview.Application) {
    v.app = app
    v.form.SetFocus(app)  // Does NOT call v.screen.SetFocus(app)
}
```

Compare to the correct implementation in `BaseView`:

```go
func (v *BaseView) SetFocus(app *tview.Application) {
    v.app = app
    v.screen.SetFocus(app)  // Correctly initializes Screen.app
}
```

### Event Flow

1. User presses TAB key
2. `App.handleGlobalKeys()` receives the event (`internal/ui/app.go:118-157`)
3. Event is passed to `view.HandleKey(event)` (line 153)
4. `CustomerView.HandleKey()` delegates to `BaseView.HandleKey()` (line 435)
5. `BaseView.HandleKey()` calls `v.screen.HandleKey(event)` (line 51)
6. `Screen.HandleKey()` checks `s.app != nil` - **fails because `s.app` was never set**
7. TAB event is returned unhandled instead of calling `s.form.NextField(s.app)`

## Affected Components

- `internal/ui/components/screen.go` - TAB key handler has correct logic but depends on `s.app`
- `internal/ui/views/customer.go` - SetFocus override bypasses screen initialization
- `internal/ui/views/motor.go` - SetFocus override bypasses screen initialization
- `internal/ui/views/endowment.go` - SetFocus override bypasses screen initialization
- `internal/ui/views/house.go` - SetFocus override bypasses screen initialization
- `internal/ui/views/policy_placeholders.go` - Two views with SetFocus overrides
- `internal/ui/views/main_menu.go` - SetFocus override (menu-only, no form)

## Proposed Solution

### Option 1: Fix Each View's SetFocus Method (Recommended)

Update each view's `SetFocus()` method to call `v.screen.SetFocus(app)` so that `Screen.app` is properly initialized:

```go
func (v *CustomerView) SetFocus(app *tview.Application) {
    v.app = app
    v.screen.SetFocus(app)  // Add this line to initialize Screen.app
}
```

This approach:
- Maintains the existing architecture
- Initializes `Screen.app` so TAB key handling works
- `Screen.SetFocus()` already handles setting focus to the form or menu appropriately

### Option 2: Remove `s.app != nil` Check

Remove the `s.app != nil` check from `Screen.HandleKey()` and pass `app` as a parameter:

```go
func (s *Screen) HandleKey(event *tcell.EventKey, app *tview.Application) *tcell.EventKey {
    switch event.Key() {
    case tcell.KeyTab:
        if s.form != nil {
            s.form.NextField(app)
            return nil
        }
    ...
```

This would require updating the call chain but avoids the need to maintain state.

### Recommendation

**Option 1** is recommended because:
1. It's a minimal change (modify one line in each affected view)
2. It aligns with how `BaseView.SetFocus()` was designed to work
3. It doesn't require changing interfaces or the event handling chain
4. `Screen.SetFocus()` already contains the correct logic for focus management

## Edge Cases

- **MainMenuView**: This view has no form, only a menu. The TAB key wouldn't navigate anything meaningful, but the fix should still set `v.screen.SetFocus(app)` for consistency. The `Screen.HandleKey()` condition `s.form != nil` will prevent any action when there's no form.

## Test Plan

1. Start the application
2. Navigate to Customer screen
3. Press TAB - should move to next form field
4. Press Shift+TAB - should move to previous form field
5. Repeat for Motor, Endowment, House, Commercial, and Claim screens
6. Verify Main Menu screen doesn't break (no form navigation expected)

---

## Implementation Notes

### Fix Applied

Applied **Option 1** as recommended. Updated the `SetFocus()` method in all affected views to call `v.screen.SetFocus(app)` instead of bypassing it.

### Files Modified

| File | Change |
|------|--------|
| `internal/ui/views/customer.go:441` | Changed `v.form.SetFocus(app)` to `v.screen.SetFocus(app)` |
| `internal/ui/views/motor.go:598` | Changed `v.form.SetFocus(app)` to `v.screen.SetFocus(app)` |
| `internal/ui/views/endowment.go:545` | Changed `v.form.SetFocus(app)` to `v.screen.SetFocus(app)` |
| `internal/ui/views/house.go:521` | Changed `v.form.SetFocus(app)` to `v.screen.SetFocus(app)` |
| `internal/ui/views/policy_placeholders.go:135` | Changed `v.form.SetFocus(app)` to `v.screen.SetFocus(app)` (CommercialPolicyView) |
| `internal/ui/views/policy_placeholders.go:267` | Changed `v.form.SetFocus(app)` to `v.screen.SetFocus(app)` (ClaimView) |
| `internal/ui/views/main_menu.go:93` | Changed `app.SetFocus(v.menu.OptionInput())` to `v.screen.SetFocus(app)` |

### Test Results

- Build: PASS (`go build ./...`)
- Tests: PASS (`go test ./...`)
  - `internal/repository`: ok
  - `internal/service`: ok
  - No UI tests exist in the codebase

### Why This Fix Works

The fix ensures that when any view's `SetFocus()` is called, it delegates to `Screen.SetFocus(app)` which:
1. Sets `s.app = app` on the Screen component
2. Properly initializes focus on the form or menu

With `s.app` properly initialized, the TAB key handler in `Screen.HandleKey()` can now execute `s.form.NextField(s.app)` and `s.form.PrevField(s.app)` successfully.
