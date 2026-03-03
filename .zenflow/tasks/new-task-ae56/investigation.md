# Investigation: Menu Options Not Displayed Correctly

## Bug Summary

The Customer screen (and potentially other screens) does not display all menu options (1-4) as stated in the BMS specification. The user reports:
- Customer screen only shows options 1 and 2
- BMS reference has "4. Cust Update"
- Option 3 appears to be dynamic/conditional

## Investigation Findings

### Current Implementation Analysis

#### Menu System (`internal/ui/components/menu.go`)

The menu system uses the `updateDisplay()` function (lines 73-85) to render options:

```go
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
```

**Key observation**: Disabled options ARE rendered, but with `[gray]..[-]` color formatting. If the label is empty, the line appears as `"3. "` with nothing visible.

#### Customer Screen (`internal/ui/views/customer.go`)

Current menu definition (lines 31-36):
```go
v.menu = components.NewMenu()
v.menu.AddOption("1", "Cust Inquiry", true)
v.menu.AddOption("2", "Cust Add", true)
v.menu.AddOption("3", "", false)           // Reserved/blank in original
v.menu.AddOption("4", "Cust Update", true)
```

**Issue identified**: All 4 options ARE defined, but:
1. Option 3 has an empty label and is disabled - renders as `"3. "` (essentially invisible)
2. Option 4 SHOULD be visible, but the user reports it's not showing

#### All Screen Menu Configurations

| Screen | Option 1 | Option 2 | Option 3 | Option 4 |
|--------|----------|----------|----------|----------|
| **Customer** | Cust Inquiry | Cust Add | *(empty/disabled)* | Cust Update |
| **Motor** | Policy Inquiry | Policy Add | Policy Delete | Policy Update |
| **Endowment** | Policy Inquiry | Policy Add | Policy Delete | Policy Update |
| **House** | Policy Inquiry | Policy Add | Policy Delete | Policy Update |
| **Commercial** | Policy Inquiry | Policy Add | Policy Delete | *(not defined)* |
| **Claim** | Claim Inquiry | Claim Add | *(not defined)* | *(not defined)* |

### Root Cause Analysis

There are potentially **two issues**:

#### Issue 1: Option 4 Not Rendering
The user reports "Customer only shows option 1 and 2" even though option 4 ("Cust Update") is defined and enabled. This suggests a rendering problem where option 4 is either:
- Being clipped/cut off in the display area
- The newline formatting may cause display issues
- A rendering bug in tview/tcell

#### Issue 2: Option 3 Displays as Empty Line
When option 3 is added with an empty label (`v.menu.AddOption("3", "", false)`), it renders as:
```
1. Cust Inquiry
2. Cust Add
3.
4. Cust Update
```

The "3. " line appears as a blank/gray line which is confusing. The user sees this as "missing" options.

#### Issue 3: Dynamic Option 3 Not Implemented
The bug report mentions "option 3 that is dynamic available". This suggests the original BMS had conditional logic:
- Option 3 might be "Cust Delete" that appears only when certain conditions are met
- This dynamic behavior is not implemented in the current code

### Affected Components

1. `internal/ui/components/menu.go` - Menu rendering logic
2. `internal/ui/views/customer.go` - Customer screen menu setup
3. Potentially other screens with similar patterns

### Proposed Solution

#### Fix 1: Ensure Option 4 Renders Properly
Verify the menu display area has sufficient height to show all 4 options. Check `internal/ui/components/screen.go` for the menu area configuration.

#### Fix 2: Handle Empty/Disabled Options Better
Modify `updateDisplay()` to either:
- Skip disabled options with empty labels entirely (don't render)
- Show a placeholder like `"3. (Reserved)"` for disabled options
- Only show enabled options

**Recommended approach**: Skip options that are disabled AND have empty labels:

```go
func (m *Menu) updateDisplay() {
    text := ""
    for _, opt := range m.options {
        if opt.Label == "" && !opt.Enabled {
            // Skip disabled options with empty labels
            continue
        }
        if opt.Enabled {
            text += fmt.Sprintf("%s. %s\n", opt.Key, opt.Label)
        } else {
            // Show disabled options in a different color
            text += fmt.Sprintf("[gray]%s. %s[-]\n", opt.Key, opt.Label)
        }
    }
    m.optionsDisplay.SetText(text)
}
```

#### Fix 3: Consider Dynamic Option 3 (Future Enhancement)
If the original BMS had dynamic behavior for option 3, this would require:
- Adding a method to enable/disable options at runtime
- Implementing business logic for when "Cust Delete" should appear
- This is likely a separate enhancement, not part of the current bug fix

### Implementation Plan

1. **Step 1**: Modify `updateDisplay()` in `menu.go` to skip disabled options with empty labels
2. **Step 2**: Verify all screens render all their enabled options correctly
3. **Step 3**: Test the Customer screen specifically to ensure options 1, 2, and 4 all appear
4. **Step 4**: Review other screens (Commercial, Claim) for similar issues

### Test Verification

After implementing the fix:
1. Customer screen should display:
   ```
   1. Cust Inquiry
   2. Cust Add
   4. Cust Update
   ```
2. Policy screens should continue displaying all 4 options
3. Commercial should display options 1-3
4. Claim should display options 1-2

### Additional Notes

The codebase has several compilation errors related to `components.FormatCustomerNum` and `components.FormatPolicyNum` being undefined. These are separate issues that should be addressed but are not directly related to the menu display bug.
