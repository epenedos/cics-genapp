package views

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/cicsdev/genapp/internal/models"
	"github.com/cicsdev/genapp/internal/service"
	"github.com/cicsdev/genapp/internal/ui/components"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// MotorPolicyView implements the motor policy screen (SSMAPP1 equivalent).
// Provides motor policy inquiry, add, update, and delete operations.
type MotorPolicyView struct {
	*BaseView
	menu       *components.Menu
	form       *components.Form
	onNavigate func(screen string)
}

// NewMotorPolicyView creates a new motor policy view matching SSMAPP1.
func NewMotorPolicyView() *MotorPolicyView {
	v := &MotorPolicyView{
		BaseView: NewBaseView("motor", "SSP1", "General Insurance Motor Policy Menu"),
	}

	// Create the menu matching SSMAPP1
	v.menu = components.PolicyMenu()

	// Create the form with motor policy fields (BMS positions from SSMAPP1)
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
		Row:          4,
		LabelColumn:  30,
		Column:       50,
	})
	v.form.AddField(&components.FormField{
		Label:        "Cust Number",
		Name:         "customer_num",
		MaxLength:    10,
		FieldType:    components.FieldTypeNumeric,
		Editable:     true,
		RightJustify: true,
		ZeroPad:      true,
		Row:          5,
		LabelColumn:  30,
		Column:       50,
	})
	v.form.AddField(&components.FormField{
		Label:       "Issue date",
		Name:        "issue_date",
		MaxLength:   10,
		FieldType:   components.FieldTypeDate,
		Editable:    true,
		Row:         6,
		LabelColumn: 30,
		Column:      50,
	})
	v.form.AddField(&components.FormField{
		Label:       "Expiry date",
		Name:        "expiry_date",
		MaxLength:   10,
		FieldType:   components.FieldTypeDate,
		Editable:    true,
		Row:         7,
		LabelColumn: 30,
		Column:      50,
	})
	v.form.AddField(&components.FormField{
		Label:       "Car Make",
		Name:        "car_make",
		MaxLength:   15,
		FieldType:   components.FieldTypeText,
		Editable:    true,
		Row:         8,
		LabelColumn: 30,
		Column:      50,
	})
	v.form.AddField(&components.FormField{
		Label:       "Car Model",
		Name:        "car_model",
		MaxLength:   15,
		FieldType:   components.FieldTypeText,
		Editable:    true,
		Row:         9,
		LabelColumn: 30,
		Column:      50,
	})
	v.form.AddField(&components.FormField{
		Label:        "Car Value",
		Name:         "car_value",
		MaxLength:    8,
		FieldType:    components.FieldTypeNumeric,
		Editable:     true,
		RightJustify: true,
		Row:          10,
		LabelColumn:  30,
		Column:       50,
	})
	v.form.AddField(&components.FormField{
		Label:       "Registration",
		Name:        "registration",
		MaxLength:   7,
		FieldType:   components.FieldTypeText,
		Editable:    true,
		Row:         11,
		LabelColumn: 30,
		Column:      50,
	})
	v.form.AddField(&components.FormField{
		Label:       "Car Colour",
		Name:        "car_colour",
		MaxLength:   8,
		FieldType:   components.FieldTypeText,
		Editable:    true,
		Row:         12,
		LabelColumn: 30,
		Column:      50,
	})
	v.form.AddField(&components.FormField{
		Label:        "CC",
		Name:         "cc",
		MaxLength:    8,
		FieldType:    components.FieldTypeNumeric,
		Editable:     true,
		RightJustify: true,
		Row:          13,
		LabelColumn:  30,
		Column:       50,
	})
	v.form.AddField(&components.FormField{
		Label:       "Manufacture Date",
		Name:        "manufactured",
		MaxLength:   10,
		FieldType:   components.FieldTypeDate,
		Editable:    true,
		Row:         14,
		LabelColumn: 30,
		Column:      50,
	})
	v.form.AddField(&components.FormField{
		Label:        "No. Accidents",
		Name:         "accidents",
		MaxLength:    6,
		FieldType:    components.FieldTypeNumeric,
		Editable:     true,
		RightJustify: true,
		Row:          15,
		LabelColumn:  30,
		Column:       50,
	})
	v.form.AddField(&components.FormField{
		Label:        "Policy Premium",
		Name:         "premium",
		MaxLength:    8,
		FieldType:    components.FieldTypeNumeric,
		Editable:     true,
		RightJustify: true,
		Row:          16,
		LabelColumn:  30,
		Column:       50,
	})

	// Set up the screen
	v.screen.SetMenu(v.menu)
	v.screen.SetForm(v.form)

	// Set up Enter key handler
	v.SetOnSubmit(v.handleSubmit)

	return v
}

// SetOnNavigate sets the callback for screen navigation.
func (v *MotorPolicyView) SetOnNavigate(fn func(screen string)) {
	v.onNavigate = fn
}

// handleSubmit processes the form submission.
func (v *MotorPolicyView) handleSubmit() {
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

// handleInquiry retrieves and displays motor policy data.
func (v *MotorPolicyView) handleInquiry() {
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

	// Verify it's a motor policy
	if policy.PolicyType != models.PolicyTypeMotor {
		v.ShowError("Policy is not a motor policy")
		return
	}

	// Populate form with policy data
	v.populateFormFromPolicy(policy)
	v.ShowSuccess("Motor policy inquiry successful")
}

// handleAdd creates a new motor policy from form data.
func (v *MotorPolicyView) handleAdd() {
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
	v.ShowSuccess("Motor policy " + result.PolicyNum + " added")
}

// handleUpdate modifies an existing motor policy.
func (v *MotorPolicyView) handleUpdate() {
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

	v.ShowSuccess("Motor policy " + policyNum + " updated")
}

// handleDelete removes a motor policy.
func (v *MotorPolicyView) handleDelete() {
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
	v.ShowSuccess("Motor policy " + policyNum + " deleted")
}

// getFormattedPolicyNum returns the policy number padded to 10 digits.
func (v *MotorPolicyView) getFormattedPolicyNum() string {
	raw := strings.TrimSpace(v.form.GetValue("policy_num"))
	if raw == "" {
		return ""
	}
	return components.FormatPolicyNum(raw)
}

// getFormattedCustomerNum returns the customer number padded to 10 digits.
func (v *MotorPolicyView) getFormattedCustomerNum() string {
	raw := strings.TrimSpace(v.form.GetValue("customer_num"))
	if raw == "" {
		return ""
	}
	return components.FormatCustomerNum(raw)
}

// validateAddForm validates fields required for adding a motor policy.
func (v *MotorPolicyView) validateAddForm() string {
	carMake := strings.TrimSpace(v.form.GetValue("car_make"))
	if carMake == "" {
		return "Car Make is required"
	}

	// Validate date formats
	if err := v.validateDates(); err != "" {
		return err
	}

	return ""
}

// validateUpdateForm validates fields for updating a motor policy.
func (v *MotorPolicyView) validateUpdateForm() string {
	return v.validateDates()
}

// validateDates validates date format fields.
func (v *MotorPolicyView) validateDates() string {
	dateFields := []struct {
		name  string
		label string
	}{
		{"issue_date", "Issue date"},
		{"expiry_date", "Expiry date"},
		{"manufactured", "Manufacture Date"},
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
func (v *MotorPolicyView) buildAddInput() (*service.AddPolicyInput, error) {
	values := v.form.GetAllValues()

	input := &service.AddPolicyInput{
		CustomerNum: components.FormatCustomerNum(values["customer_num"]),
		PolicyType:  models.PolicyTypeMotor,
		Motor:       &service.AddMotorInput{},
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

	// Motor-specific fields
	input.Motor.Make = strings.TrimSpace(values["car_make"])
	input.Motor.Model = strings.TrimSpace(values["car_model"])
	input.Motor.RegNumber = strings.TrimSpace(values["registration"])
	input.Motor.Colour = strings.TrimSpace(values["car_colour"])

	if carValue := strings.TrimSpace(values["car_value"]); carValue != "" {
		if val, err := strconv.ParseFloat(carValue, 64); err == nil {
			input.Motor.Value = val
		}
	}

	if cc := strings.TrimSpace(values["cc"]); cc != "" {
		if val, err := strconv.Atoi(cc); err == nil {
			input.Motor.CC = val
		}
	}

	if manufactured := strings.TrimSpace(values["manufactured"]); manufactured != "" {
		t, err := time.Parse("2006-01-02", manufactured)
		if err != nil {
			return nil, errors.New("Manufacture Date must be in yyyy-mm-dd format")
		}
		input.Motor.Manufactured = &t
	}

	if accidents := strings.TrimSpace(values["accidents"]); accidents != "" {
		if val, err := strconv.Atoi(accidents); err == nil {
			input.Motor.Accidents = val
		}
	}

	if premium := strings.TrimSpace(values["premium"]); premium != "" {
		if val, err := strconv.ParseFloat(premium, 64); err == nil {
			input.Motor.Premium = val
		}
	}

	return input, nil
}

// buildUpdateInput creates UpdatePolicyInput from form values.
func (v *MotorPolicyView) buildUpdateInput() *service.UpdatePolicyInput {
	values := v.form.GetAllValues()

	input := &service.UpdatePolicyInput{
		Motor: &service.UpdateMotorInput{},
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

	// Motor-specific fields
	if val := strings.TrimSpace(values["car_make"]); val != "" {
		input.Motor.Make = &val
	}
	if val := strings.TrimSpace(values["car_model"]); val != "" {
		input.Motor.Model = &val
	}
	if val := strings.TrimSpace(values["registration"]); val != "" {
		input.Motor.RegNumber = &val
	}
	if val := strings.TrimSpace(values["car_colour"]); val != "" {
		input.Motor.Colour = &val
	}

	if carValue := strings.TrimSpace(values["car_value"]); carValue != "" {
		if val, err := strconv.ParseFloat(carValue, 64); err == nil {
			input.Motor.Value = &val
		}
	}

	if cc := strings.TrimSpace(values["cc"]); cc != "" {
		if val, err := strconv.Atoi(cc); err == nil {
			input.Motor.CC = &val
		}
	}

	if manufactured := strings.TrimSpace(values["manufactured"]); manufactured != "" {
		if t, err := time.Parse("2006-01-02", manufactured); err == nil {
			input.Motor.Manufactured = &t
		}
	}

	if accidents := strings.TrimSpace(values["accidents"]); accidents != "" {
		if val, err := strconv.Atoi(accidents); err == nil {
			input.Motor.Accidents = &val
		}
	}

	if premium := strings.TrimSpace(values["premium"]); premium != "" {
		if val, err := strconv.ParseFloat(premium, 64); err == nil {
			input.Motor.Premium = &val
		}
	}

	return input
}

// populateFormFromPolicy fills form fields with policy data.
func (v *MotorPolicyView) populateFormFromPolicy(p *models.Policy) {
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

	// Motor-specific fields
	if p.Motor != nil {
		v.form.SetValue("car_make", p.Motor.GetMake())
		v.form.SetValue("car_model", p.Motor.GetModel())
		v.form.SetValue("car_value", formatNumeric(p.Motor.GetValue()))
		v.form.SetValue("registration", p.Motor.GetRegNumber())
		v.form.SetValue("car_colour", p.Motor.GetColour())
		v.form.SetValue("cc", formatInt(p.Motor.GetCC()))

		if p.Motor.Manufactured.Valid {
			v.form.SetValue("manufactured", p.Motor.Manufactured.Time.Format("2006-01-02"))
		} else {
			v.form.SetValue("manufactured", "")
		}

		v.form.SetValue("accidents", formatInt(p.Motor.GetAccidents()))
		v.form.SetValue("premium", formatNumeric(p.Motor.GetPremium()))
	}
}

// HandleKey handles key events specific to the motor policy view.
func (v *MotorPolicyView) HandleKey(event *tcell.EventKey) *tcell.EventKey {
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
func (v *MotorPolicyView) SetFocus(app *tview.Application) {
	v.app = app
	v.screen.SetFocus(app)
}

// Clear resets all form fields and the menu selection.
func (v *MotorPolicyView) Clear() {
	v.form.Clear()
	v.menu.Clear()
	v.ClearError()
}

// GetPolicyNumber returns the current policy number value.
func (v *MotorPolicyView) GetPolicyNumber() string {
	return v.form.GetValue("policy_num")
}

// SetPolicyNumber sets the policy number field.
func (v *MotorPolicyView) SetPolicyNumber(num string) {
	v.form.SetValue("policy_num", components.FormatPolicyNum(num))
}

// SetCustomerNumber sets the customer number field.
func (v *MotorPolicyView) SetCustomerNumber(num string) {
	v.form.SetValue("customer_num", components.FormatCustomerNum(num))
}

// formatNumeric formats a float64 as a string, returning empty string for 0.
func formatNumeric(val float64) string {
	if val == 0 {
		return ""
	}
	return fmt.Sprintf("%.0f", val)
}

// formatInt formats an int as a string, returning empty string for 0.
func formatInt(val int) string {
	if val == 0 {
		return ""
	}
	return strconv.Itoa(val)
}
