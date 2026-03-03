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

// HousePolicyView implements the house policy screen (SSMAPP3 equivalent).
// Provides house policy inquiry, add, update, and delete operations.
type HousePolicyView struct {
	*BaseView
	menu       *components.Menu
	form       *components.Form
	onNavigate func(screen string)
}

// NewHousePolicyView creates a new house policy view matching SSMAPP3.
func NewHousePolicyView() *HousePolicyView {
	v := &HousePolicyView{
		BaseView: NewBaseView("house", "SSP3", "General Insurance House Policy Menu"),
	}

	// Create the menu matching SSMAPP3
	v.menu = components.PolicyMenu()

	// Create the form with house policy fields
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
		Label:     "Property Type",
		Name:      "property_type",
		MaxLength: 15,
		FieldType: components.FieldTypeText,
		Editable:  true,
	})
	v.form.AddField(&components.FormField{
		Label:        "Bedrooms",
		Name:         "bedrooms",
		MaxLength:    3,
		FieldType:    components.FieldTypeNumeric,
		Editable:     true,
		RightJustify: true,
	})
	v.form.AddField(&components.FormField{
		Label:        "House Value",
		Name:         "house_value",
		MaxLength:    10,
		FieldType:    components.FieldTypeNumeric,
		Editable:     true,
		RightJustify: true,
	})
	v.form.AddField(&components.FormField{
		Label:     "House Name",
		Name:      "house_name",
		MaxLength: 20,
		FieldType: components.FieldTypeText,
		Editable:  true,
	})
	v.form.AddField(&components.FormField{
		Label:     "House Number",
		Name:      "house_number",
		MaxLength: 4,
		FieldType: components.FieldTypeText,
		Editable:  true,
	})
	v.form.AddField(&components.FormField{
		Label:     "Postcode",
		Name:      "postcode",
		MaxLength: 8,
		FieldType: components.FieldTypeText,
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
func (v *HousePolicyView) SetOnNavigate(fn func(screen string)) {
	v.onNavigate = fn
}

// handleSubmit processes the form submission.
func (v *HousePolicyView) handleSubmit() {
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

// handleInquiry retrieves and displays house policy data.
func (v *HousePolicyView) handleInquiry() {
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

	// Verify it's a house policy
	if policy.PolicyType != models.PolicyTypeHouse {
		v.ShowError("Policy is not a house policy")
		return
	}

	// Populate form with policy data
	v.populateFormFromPolicy(policy)
	v.ShowSuccess("House policy inquiry successful")
}

// handleAdd creates a new house policy from form data.
func (v *HousePolicyView) handleAdd() {
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
	v.ShowSuccess("House policy " + result.PolicyNum + " added")
}

// handleUpdate modifies an existing house policy.
func (v *HousePolicyView) handleUpdate() {
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

	v.ShowSuccess("House policy " + policyNum + " updated")
}

// handleDelete removes a house policy.
func (v *HousePolicyView) handleDelete() {
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
	v.ShowSuccess("House policy " + policyNum + " deleted")
}

// getFormattedPolicyNum returns the policy number padded to 10 digits.
func (v *HousePolicyView) getFormattedPolicyNum() string {
	raw := strings.TrimSpace(v.form.GetValue("policy_num"))
	if raw == "" {
		return ""
	}
	return components.FormatPolicyNum(raw)
}

// getFormattedCustomerNum returns the customer number padded to 10 digits.
func (v *HousePolicyView) getFormattedCustomerNum() string {
	raw := strings.TrimSpace(v.form.GetValue("customer_num"))
	if raw == "" {
		return ""
	}
	return components.FormatCustomerNum(raw)
}

// validateAddForm validates fields required for adding a house policy.
func (v *HousePolicyView) validateAddForm() string {
	// Validate date formats
	if err := v.validateDates(); err != "" {
		return err
	}

	return ""
}

// validateUpdateForm validates fields for updating a house policy.
func (v *HousePolicyView) validateUpdateForm() string {
	return v.validateDates()
}

// validateDates validates date format fields.
func (v *HousePolicyView) validateDates() string {
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
func (v *HousePolicyView) buildAddInput() (*service.AddPolicyInput, error) {
	values := v.form.GetAllValues()

	input := &service.AddPolicyInput{
		CustomerNum: components.FormatCustomerNum(values["customer_num"]),
		PolicyType:  models.PolicyTypeHouse,
		House:       &service.AddHouseInput{},
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

	// House-specific fields
	input.House.PropertyType = strings.TrimSpace(values["property_type"])
	input.House.HouseName = strings.TrimSpace(values["house_name"])
	input.House.HouseNumber = strings.TrimSpace(values["house_number"])
	input.House.Postcode = strings.TrimSpace(values["postcode"])

	if bedrooms := strings.TrimSpace(values["bedrooms"]); bedrooms != "" {
		if val, err := strconv.Atoi(bedrooms); err == nil {
			input.House.Bedrooms = val
		}
	}

	if houseValue := strings.TrimSpace(values["house_value"]); houseValue != "" {
		if val, err := strconv.ParseFloat(houseValue, 64); err == nil {
			input.House.Value = val
		}
	}

	return input, nil
}

// buildUpdateInput creates UpdatePolicyInput from form values.
func (v *HousePolicyView) buildUpdateInput() *service.UpdatePolicyInput {
	values := v.form.GetAllValues()

	input := &service.UpdatePolicyInput{
		House: &service.UpdateHouseInput{},
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

	// House-specific fields
	if val := strings.TrimSpace(values["property_type"]); val != "" {
		input.House.PropertyType = &val
	}
	if val := strings.TrimSpace(values["house_name"]); val != "" {
		input.House.HouseName = &val
	}
	if val := strings.TrimSpace(values["house_number"]); val != "" {
		input.House.HouseNumber = &val
	}
	if val := strings.TrimSpace(values["postcode"]); val != "" {
		input.House.Postcode = &val
	}

	if bedrooms := strings.TrimSpace(values["bedrooms"]); bedrooms != "" {
		if val, err := strconv.Atoi(bedrooms); err == nil {
			input.House.Bedrooms = &val
		}
	}

	if houseValue := strings.TrimSpace(values["house_value"]); houseValue != "" {
		if val, err := strconv.ParseFloat(houseValue, 64); err == nil {
			input.House.Value = &val
		}
	}

	return input
}

// populateFormFromPolicy fills form fields with policy data.
func (v *HousePolicyView) populateFormFromPolicy(p *models.Policy) {
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

	// House-specific fields
	if p.House != nil {
		v.form.SetValue("property_type", p.House.GetPropertyType())
		v.form.SetValue("bedrooms", formatInt(p.House.GetBedrooms()))
		v.form.SetValue("house_value", formatNumeric(p.House.GetValue()))
		v.form.SetValue("house_name", p.House.GetHouseName())
		v.form.SetValue("house_number", p.House.GetHouseNumber())
		v.form.SetValue("postcode", p.House.GetPostcode())
	}
}

// HandleKey handles key events specific to the house policy view.
func (v *HousePolicyView) HandleKey(event *tcell.EventKey) *tcell.EventKey {
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
func (v *HousePolicyView) SetFocus(app *tview.Application) {
	v.app = app
	v.screen.SetFocus(app)
}

// Clear resets all form fields and the menu selection.
func (v *HousePolicyView) Clear() {
	v.form.Clear()
	v.menu.Clear()
	v.ClearError()
}

// GetPolicyNumber returns the current policy number value.
func (v *HousePolicyView) GetPolicyNumber() string {
	return v.form.GetValue("policy_num")
}

// SetPolicyNumber sets the policy number field.
func (v *HousePolicyView) SetPolicyNumber(num string) {
	v.form.SetValue("policy_num", components.FormatPolicyNum(num))
}

// SetCustomerNumber sets the customer number field.
func (v *HousePolicyView) SetCustomerNumber(num string) {
	v.form.SetValue("customer_num", components.FormatCustomerNum(num))
}
