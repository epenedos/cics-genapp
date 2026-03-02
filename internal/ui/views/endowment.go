package views

import (
	"context"
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/cicsdev/genapp/internal/models"
	"github.com/cicsdev/genapp/internal/service"
	"github.com/cicsdev/genapp/internal/ui/components"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// EndowmentPolicyView implements the endowment policy screen (SSMAPP2 equivalent).
// Provides endowment policy inquiry, add, update, and delete operations.
type EndowmentPolicyView struct {
	*BaseView
	menu       *components.Menu
	form       *components.Form
	onNavigate func(screen string)
}

// NewEndowmentPolicyView creates a new endowment policy view matching SSMAPP2.
func NewEndowmentPolicyView() *EndowmentPolicyView {
	v := &EndowmentPolicyView{
		BaseView: NewBaseView("endowment", "SSP2", "General Insurance Endowment Policy Menu"),
	}

	// Create the menu matching SSMAPP2
	v.menu = components.PolicyMenu()

	// Create the form with endowment policy fields
	v.form = components.NewForm()
	v.form.AddField(&components.FormField{
		Label:        "Policy Number",
		Name:         "policy_num",
		MaxLength:    10,
		FieldType:    components.FieldTypeNumeric,
		Editable:     true,
		InitialFocus: true,
		RightJustify: true,
		ZeroPad:      true,
	})
	v.form.AddField(&components.FormField{
		Label:        "Cust Number",
		Name:         "customer_num",
		MaxLength:    10,
		FieldType:    components.FieldTypeNumeric,
		Editable:     true,
		RightJustify: true,
		ZeroPad:      true,
	})
	v.form.AddField(&components.FormField{
		Label:     "Issue date",
		Name:      "issue_date",
		MaxLength: 10,
		FieldType: components.FieldTypeDate,
		Editable:  true,
	})
	v.form.AddField(&components.FormField{
		Label:     "Expiry date",
		Name:      "expiry_date",
		MaxLength: 10,
		FieldType: components.FieldTypeDate,
		Editable:  true,
	})
	v.form.AddField(&components.FormField{
		Label:     "Fund Name",
		Name:      "fund_name",
		MaxLength: 10,
		FieldType: components.FieldTypeText,
		Editable:  true,
	})
	v.form.AddField(&components.FormField{
		Label:        "Term (years)",
		Name:         "term",
		MaxLength:    2,
		FieldType:    components.FieldTypeNumeric,
		Editable:     true,
		RightJustify: true,
	})
	v.form.AddField(&components.FormField{
		Label:        "Sum Assured",
		Name:         "sum_assured",
		MaxLength:    10,
		FieldType:    components.FieldTypeNumeric,
		Editable:     true,
		RightJustify: true,
	})
	v.form.AddField(&components.FormField{
		Label:     "Life Assured",
		Name:      "life_assured",
		MaxLength: 31,
		FieldType: components.FieldTypeText,
		Editable:  true,
	})
	v.form.AddField(&components.FormField{
		Label:     "With Profits",
		Name:      "with_profits",
		MaxLength: 1,
		FieldType: components.FieldTypeYesNo,
		Editable:  true,
	})
	v.form.AddField(&components.FormField{
		Label:     "Equities",
		Name:      "equities",
		MaxLength: 1,
		FieldType: components.FieldTypeYesNo,
		Editable:  true,
	})
	v.form.AddField(&components.FormField{
		Label:     "Managed Funds",
		Name:      "managed_fund",
		MaxLength: 1,
		FieldType: components.FieldTypeYesNo,
		Editable:  true,
	})

	// Set up the screen
	v.screen.SetMenu(v.menu)
	v.screen.SetForm(v.form)

	// Set up Enter key handler
	v.SetOnSubmit(v.handleSubmit)

	return v
}

// SetOnNavigate sets the callback for screen navigation.
func (v *EndowmentPolicyView) SetOnNavigate(fn func(screen string)) {
	v.onNavigate = fn
}

// handleSubmit processes the form submission.
func (v *EndowmentPolicyView) handleSubmit() {
	option, valid := v.menu.ProcessSelection()
	if !valid {
		v.ShowError("Please select an option (1-4)")
		return
	}

	v.ClearError()

	// Process based on option
	switch option.Key {
	case "1": // Inquiry
		v.handleInquiry()
	case "2": // Add
		v.handleAdd()
	case "3": // Delete
		v.handleDelete()
	case "4": // Update
		v.handleUpdate()
	}
}

// handleInquiry retrieves and displays endowment policy data.
func (v *EndowmentPolicyView) handleInquiry() {
	policyNum := v.getFormattedPolicyNum()
	if policyNum == "" {
		v.ShowError("Policy Number required for inquiry")
		return
	}

	if v.PolicyService() == nil {
		v.ShowError("Service not available")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	policy, err := v.PolicyService().Get(ctx, policyNum)
	if err != nil {
		if errors.Is(err, service.ErrPolicyNotFound) {
			v.ShowError("Policy not found")
		} else if errors.Is(err, service.ErrInvalidPolicyNumber) {
			v.ShowError("Invalid policy number format")
		} else {
			v.ShowError("Error: " + truncateError(err.Error(), 50))
		}
		return
	}

	// Verify it's an endowment policy
	if policy.PolicyType != models.PolicyTypeEndowment {
		v.ShowError("Policy is not an endowment policy")
		return
	}

	// Populate form with policy data
	v.populateFormFromPolicy(policy)
	v.ShowSuccess("Endowment policy inquiry successful")
}

// handleAdd creates a new endowment policy from form data.
func (v *EndowmentPolicyView) handleAdd() {
	// For add operations, policy number field should be empty
	existingNum := strings.TrimSpace(v.form.GetValue("policy_num"))
	if existingNum != "" && existingNum != "0000000000" {
		v.ShowError("Clear policy number for new policy")
		return
	}

	// Customer number is required
	customerNum := v.getFormattedCustomerNum()
	if customerNum == "" || customerNum == "0000000000" {
		v.ShowError("Customer Number is required")
		return
	}

	// Validate form for add operation
	if errMsg := v.validateAddForm(); errMsg != "" {
		v.ShowError(errMsg)
		return
	}

	if v.PolicyService() == nil {
		v.ShowError("Service not available")
		return
	}

	// Build the add input from form values
	input, err := v.buildAddInput()
	if err != nil {
		v.ShowError(err.Error())
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := v.PolicyService().Add(ctx, input)
	if err != nil {
		if errors.Is(err, service.ErrValidationFailed) {
			v.ShowError("Validation: " + extractValidationMessage(err))
		} else if errors.Is(err, service.ErrCustomerNotFound) {
			v.ShowError("Customer not found")
		} else {
			v.ShowError("Error: " + truncateError(err.Error(), 50))
		}
		return
	}

	// Set the generated policy number
	v.form.SetValue("policy_num", components.FormatPolicyNum(result.PolicyNum))
	v.ShowSuccess("Endowment policy " + result.PolicyNum + " added")
}

// handleUpdate modifies an existing endowment policy.
func (v *EndowmentPolicyView) handleUpdate() {
	policyNum := v.getFormattedPolicyNum()
	if policyNum == "" {
		v.ShowError("Policy Number required for update")
		return
	}

	// Validate form for update operation
	if errMsg := v.validateUpdateForm(); errMsg != "" {
		v.ShowError(errMsg)
		return
	}

	if v.PolicyService() == nil {
		v.ShowError("Service not available")
		return
	}

	// Build the update input from form values
	input := v.buildUpdateInput()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := v.PolicyService().Update(ctx, policyNum, input)
	if err != nil {
		if errors.Is(err, service.ErrPolicyNotFound) {
			v.ShowError("Policy not found")
		} else if errors.Is(err, service.ErrInvalidPolicyNumber) {
			v.ShowError("Invalid policy number format")
		} else if errors.Is(err, service.ErrValidationFailed) {
			v.ShowError("Validation: " + extractValidationMessage(err))
		} else {
			v.ShowError("Error: " + truncateError(err.Error(), 50))
		}
		return
	}

	v.ShowSuccess("Endowment policy " + policyNum + " updated")
}

// handleDelete removes an endowment policy.
func (v *EndowmentPolicyView) handleDelete() {
	policyNum := v.getFormattedPolicyNum()
	if policyNum == "" {
		v.ShowError("Policy Number required for delete")
		return
	}

	if v.PolicyService() == nil {
		v.ShowError("Service not available")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := v.PolicyService().Delete(ctx, policyNum)
	if err != nil {
		if errors.Is(err, service.ErrPolicyNotFound) {
			v.ShowError("Policy not found")
		} else if errors.Is(err, service.ErrInvalidPolicyNumber) {
			v.ShowError("Invalid policy number format")
		} else {
			v.ShowError("Error: " + truncateError(err.Error(), 50))
		}
		return
	}

	// Clear the form after successful delete
	v.form.Clear()
	v.ShowSuccess("Endowment policy " + policyNum + " deleted")
}

// getFormattedPolicyNum returns the policy number padded to 10 digits.
func (v *EndowmentPolicyView) getFormattedPolicyNum() string {
	raw := strings.TrimSpace(v.form.GetValue("policy_num"))
	if raw == "" {
		return ""
	}
	return components.FormatPolicyNum(raw)
}

// getFormattedCustomerNum returns the customer number padded to 10 digits.
func (v *EndowmentPolicyView) getFormattedCustomerNum() string {
	raw := strings.TrimSpace(v.form.GetValue("customer_num"))
	if raw == "" {
		return ""
	}
	return components.FormatCustomerNum(raw)
}

// validateAddForm validates fields required for adding an endowment policy.
func (v *EndowmentPolicyView) validateAddForm() string {
	// Validate date formats
	if err := v.validateDates(); err != "" {
		return err
	}

	// Validate term is within range
	if term := strings.TrimSpace(v.form.GetValue("term")); term != "" {
		if val, err := strconv.Atoi(term); err == nil {
			if val < 0 || val > 99 {
				return "Term must be between 0 and 99 years"
			}
		}
	}

	return ""
}

// validateUpdateForm validates fields for updating an endowment policy.
func (v *EndowmentPolicyView) validateUpdateForm() string {
	return v.validateAddForm()
}

// validateDates validates date format fields.
func (v *EndowmentPolicyView) validateDates() string {
	dateFields := []struct {
		name  string
		label string
	}{
		{"issue_date", "Issue date"},
		{"expiry_date", "Expiry date"},
	}

	for _, df := range dateFields {
		dateVal := strings.TrimSpace(v.form.GetValue(df.name))
		if dateVal != "" {
			if _, err := time.Parse("2006-01-02", dateVal); err != nil {
				return df.label + " must be in yyyy-mm-dd format"
			}
		}
	}

	return ""
}

// buildAddInput creates AddPolicyInput from form values.
func (v *EndowmentPolicyView) buildAddInput() (*service.AddPolicyInput, error) {
	values := v.form.GetAllValues()

	input := &service.AddPolicyInput{
		CustomerNum: components.FormatCustomerNum(values["customer_num"]),
		PolicyType:  models.PolicyTypeEndowment,
		Endowment:   &service.AddEndowmentInput{},
	}

	// Parse dates
	if issueDate := strings.TrimSpace(values["issue_date"]); issueDate != "" {
		t, err := time.Parse("2006-01-02", issueDate)
		if err != nil {
			return nil, errors.New("Issue date must be in yyyy-mm-dd format")
		}
		input.IssueDate = &t
	}

	if expiryDate := strings.TrimSpace(values["expiry_date"]); expiryDate != "" {
		t, err := time.Parse("2006-01-02", expiryDate)
		if err != nil {
			return nil, errors.New("Expiry date must be in yyyy-mm-dd format")
		}
		input.ExpiryDate = &t
	}

	// Endowment-specific fields
	input.Endowment.FundName = strings.TrimSpace(values["fund_name"])
	input.Endowment.LifeAssured = strings.TrimSpace(values["life_assured"])

	if term := strings.TrimSpace(values["term"]); term != "" {
		if val, err := strconv.Atoi(term); err == nil {
			input.Endowment.Term = val
		}
	}

	if sumAssured := strings.TrimSpace(values["sum_assured"]); sumAssured != "" {
		if val, err := strconv.ParseFloat(sumAssured, 64); err == nil {
			input.Endowment.SumAssured = val
		}
	}

	// Y/N fields
	input.Endowment.WithProfits = parseYesNo(values["with_profits"])
	input.Endowment.Equities = parseYesNo(values["equities"])
	input.Endowment.ManagedFund = parseYesNo(values["managed_fund"])

	return input, nil
}

// buildUpdateInput creates UpdatePolicyInput from form values.
func (v *EndowmentPolicyView) buildUpdateInput() *service.UpdatePolicyInput {
	values := v.form.GetAllValues()

	input := &service.UpdatePolicyInput{
		Endowment: &service.UpdateEndowmentInput{},
	}

	// Parse dates
	if issueDate := strings.TrimSpace(values["issue_date"]); issueDate != "" {
		if t, err := time.Parse("2006-01-02", issueDate); err == nil {
			input.IssueDate = &t
		}
	}

	if expiryDate := strings.TrimSpace(values["expiry_date"]); expiryDate != "" {
		if t, err := time.Parse("2006-01-02", expiryDate); err == nil {
			input.ExpiryDate = &t
		}
	}

	// Endowment-specific fields
	if val := strings.TrimSpace(values["fund_name"]); val != "" {
		input.Endowment.FundName = &val
	}
	if val := strings.TrimSpace(values["life_assured"]); val != "" {
		input.Endowment.LifeAssured = &val
	}

	if term := strings.TrimSpace(values["term"]); term != "" {
		if val, err := strconv.Atoi(term); err == nil {
			input.Endowment.Term = &val
		}
	}

	if sumAssured := strings.TrimSpace(values["sum_assured"]); sumAssured != "" {
		if val, err := strconv.ParseFloat(sumAssured, 64); err == nil {
			input.Endowment.SumAssured = &val
		}
	}

	// Y/N fields - always set these as they may be explicitly toggled
	withProfits := parseYesNo(values["with_profits"])
	input.Endowment.WithProfits = &withProfits

	equities := parseYesNo(values["equities"])
	input.Endowment.Equities = &equities

	managedFund := parseYesNo(values["managed_fund"])
	input.Endowment.ManagedFund = &managedFund

	return input
}

// populateFormFromPolicy fills form fields with policy data.
func (v *EndowmentPolicyView) populateFormFromPolicy(p *models.Policy) {
	v.form.SetValue("policy_num", components.FormatPolicyNum(p.PolicyNum))
	v.form.SetValue("customer_num", components.FormatCustomerNum(p.CustomerNum))

	// Format dates
	if p.IssueDate.Valid {
		v.form.SetValue("issue_date", p.IssueDate.Time.Format("2006-01-02"))
	} else {
		v.form.SetValue("issue_date", "")
	}

	if p.ExpiryDate.Valid {
		v.form.SetValue("expiry_date", p.ExpiryDate.Time.Format("2006-01-02"))
	} else {
		v.form.SetValue("expiry_date", "")
	}

	// Endowment-specific fields
	if p.Endowment != nil {
		v.form.SetValue("fund_name", p.Endowment.GetFundName())
		v.form.SetValue("term", formatInt(p.Endowment.GetTerm()))
		v.form.SetValue("sum_assured", formatNumeric(p.Endowment.GetSumAssured()))
		v.form.SetValue("life_assured", p.Endowment.GetLifeAssured())
		v.form.SetValue("with_profits", formatYesNo(p.Endowment.GetWithProfits()))
		v.form.SetValue("equities", formatYesNo(p.Endowment.GetEquities()))
		v.form.SetValue("managed_fund", formatYesNo(p.Endowment.GetManagedFund()))
	}
}

// HandleKey handles key events specific to the endowment policy view.
func (v *EndowmentPolicyView) HandleKey(event *tcell.EventKey) *tcell.EventKey {
	// Handle F-keys for navigation
	switch event.Key() {
	case tcell.KeyF6:
		// Navigate back to customer screen
		if v.onNavigate != nil {
			v.onNavigate("customer")
		}
		return nil
	}

	return v.BaseView.HandleKey(event)
}

// SetFocus sets focus to the policy number field.
func (v *EndowmentPolicyView) SetFocus(app *tview.Application) {
	v.app = app
	v.form.SetFocus(app)
}

// Clear resets all form fields and the menu selection.
func (v *EndowmentPolicyView) Clear() {
	v.form.Clear()
	v.menu.Clear()
	v.ClearError()
}

// GetPolicyNumber returns the current policy number value.
func (v *EndowmentPolicyView) GetPolicyNumber() string {
	return v.form.GetValue("policy_num")
}

// SetPolicyNumber sets the policy number field.
func (v *EndowmentPolicyView) SetPolicyNumber(num string) {
	v.form.SetValue("policy_num", components.FormatPolicyNum(num))
}

// SetCustomerNumber sets the customer number field.
func (v *EndowmentPolicyView) SetCustomerNumber(num string) {
	v.form.SetValue("customer_num", components.FormatCustomerNum(num))
}

// parseYesNo converts Y/N string to boolean.
func parseYesNo(val string) bool {
	return strings.ToUpper(strings.TrimSpace(val)) == "Y"
}

// formatYesNo converts boolean to Y/N string.
func formatYesNo(val bool) string {
	if val {
		return "Y"
	}
	return "N"
}
