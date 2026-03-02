package views

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/cicsdev/genapp/internal/models"
	"github.com/cicsdev/genapp/internal/service"
	"github.com/cicsdev/genapp/internal/ui/components"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// CustomerView implements the customer menu screen (SSMAPC1 equivalent).
// Provides customer inquiry, add, and update operations.
type CustomerView struct {
	*BaseView
	menu       *components.Menu
	form       *components.Form
	onNavigate func(screen string)
}

// NewCustomerView creates a new customer view matching SSMAPC1.
func NewCustomerView() *CustomerView {
	v := &CustomerView{
		BaseView: NewBaseView("customer", "SSC1", "General Insurance Customer Menu"),
	}

	// Create the menu matching SSMAPC1
	v.menu = components.NewMenu()
	v.menu.AddOption("1", "Cust Inquiry", true)
	v.menu.AddOption("2", "Cust Add", true)
	v.menu.AddOption("3", "", false) // Reserved/blank in original
	v.menu.AddOption("4", "Cust Update", true)

	// Create the form with customer fields
	v.form = components.NewForm()
	v.form.AddField(&components.FormField{
		Label:        "Cust Number",
		Name:         "customer_num",
		MaxLength:    10,
		FieldType:    components.FieldTypeNumeric,
		Editable:     true,
		InitialFocus: true,
		RightJustify: true,
		ZeroPad:      true,
	})
	v.form.AddField(&components.FormField{
		Label:     "First Name",
		Name:      "first_name",
		MaxLength: 10,
		FieldType: components.FieldTypeText,
		Editable:  true,
	})
	v.form.AddField(&components.FormField{
		Label:     "Last Name",
		Name:      "last_name",
		MaxLength: 20,
		FieldType: components.FieldTypeText,
		Editable:  true,
	})
	v.form.AddField(&components.FormField{
		Label:     "DOB",
		Name:      "dob",
		MaxLength: 10,
		FieldType: components.FieldTypeDate,
		Editable:  true,
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
	v.form.AddField(&components.FormField{
		Label:     "Phone: Home",
		Name:      "phone_home",
		MaxLength: 20,
		FieldType: components.FieldTypeText,
		Editable:  true,
	})
	v.form.AddField(&components.FormField{
		Label:     "Phone: Mob",
		Name:      "phone_mobile",
		MaxLength: 20,
		FieldType: components.FieldTypeText,
		Editable:  true,
	})
	v.form.AddField(&components.FormField{
		Label:     "Email Addr",
		Name:      "email_address",
		MaxLength: 27,
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
func (v *CustomerView) SetOnNavigate(fn func(screen string)) {
	v.onNavigate = fn
}

// handleSubmit processes the form submission.
func (v *CustomerView) handleSubmit() {
	option, valid := v.menu.ProcessSelection()
	if !valid {
		v.ShowError("Please select an option (1, 2, or 4)")
		return
	}

	v.ClearError()

	// Process based on option
	switch option.Key {
	case "1": // Inquiry
		v.handleInquiry()

	case "2": // Add
		v.handleAdd()

	case "4": // Update
		v.handleUpdate()
	}
}

// handleInquiry retrieves and displays customer data.
func (v *CustomerView) handleInquiry() {
	customerNum := v.getFormattedCustomerNum()
	if customerNum == "" {
		v.ShowError("Customer Number required for inquiry")
		return
	}

	// Check if service is available
	if v.CustomerService() == nil {
		v.ShowError("Service not available")
		return
	}

	// Call the customer service
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	customer, err := v.CustomerService().Get(ctx, customerNum)
	if err != nil {
		if errors.Is(err, service.ErrCustomerNotFound) {
			v.ShowError("Customer not found")
		} else if errors.Is(err, service.ErrInvalidCustomerNumber) {
			v.ShowError("Invalid customer number format")
		} else {
			v.ShowError("Error: " + truncateError(err.Error(), 50))
		}
		return
	}

	// Populate form with customer data
	v.populateFormFromCustomer(customer)
	v.ShowSuccess("Customer inquiry successful")
}

// handleAdd creates a new customer from form data.
func (v *CustomerView) handleAdd() {
	// For add operations, customer number field should be empty or will be generated
	existingNum := strings.TrimSpace(v.form.GetValue("customer_num"))
	if existingNum != "" && existingNum != "0000000000" {
		v.ShowError("Clear customer number for new customer")
		return
	}

	// Validate form for add operation
	if errMsg := v.validateAddForm(); errMsg != "" {
		v.ShowError(errMsg)
		return
	}

	// Check if service is available
	if v.CustomerService() == nil {
		v.ShowError("Service not available")
		return
	}

	// Build the add input from form values
	input, err := v.buildAddInput()
	if err != nil {
		v.ShowError(err.Error())
		return
	}

	// Call the customer service
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := v.CustomerService().Add(ctx, input)
	if err != nil {
		if errors.Is(err, service.ErrValidationFailed) {
			v.ShowError("Validation: " + extractValidationMessage(err))
		} else {
			v.ShowError("Error: " + truncateError(err.Error(), 50))
		}
		return
	}

	// Set the generated customer number
	v.form.SetValue("customer_num", components.FormatCustomerNum(result.CustomerNum))
	v.ShowSuccess("Customer " + result.CustomerNum + " added")
}

// handleUpdate modifies an existing customer.
func (v *CustomerView) handleUpdate() {
	customerNum := v.getFormattedCustomerNum()
	if customerNum == "" {
		v.ShowError("Customer Number required for update")
		return
	}

	// Validate form for update operation
	if errMsg := v.validateUpdateForm(); errMsg != "" {
		v.ShowError(errMsg)
		return
	}

	// Check if service is available
	if v.CustomerService() == nil {
		v.ShowError("Service not available")
		return
	}

	// Build the update input from form values
	input := v.buildUpdateInput()

	// Call the customer service
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := v.CustomerService().Update(ctx, customerNum, input)
	if err != nil {
		if errors.Is(err, service.ErrCustomerNotFound) {
			v.ShowError("Customer not found")
		} else if errors.Is(err, service.ErrInvalidCustomerNumber) {
			v.ShowError("Invalid customer number format")
		} else if errors.Is(err, service.ErrValidationFailed) {
			v.ShowError("Validation: " + extractValidationMessage(err))
		} else {
			v.ShowError("Error: " + truncateError(err.Error(), 50))
		}
		return
	}

	v.ShowSuccess("Customer " + customerNum + " updated")
}

// getFormattedCustomerNum returns the customer number padded to 10 digits.
func (v *CustomerView) getFormattedCustomerNum() string {
	raw := strings.TrimSpace(v.form.GetValue("customer_num"))
	if raw == "" {
		return ""
	}
	return components.FormatCustomerNum(raw)
}

// validateAddForm validates fields required for adding a customer.
func (v *CustomerView) validateAddForm() string {
	lastName := strings.TrimSpace(v.form.GetValue("last_name"))
	if lastName == "" {
		return "Last Name is required"
	}

	// Validate DOB format if provided
	dob := strings.TrimSpace(v.form.GetValue("dob"))
	if dob != "" {
		if _, err := time.Parse("2006-01-02", dob); err != nil {
			return "DOB must be in yyyy-mm-dd format"
		}
	}

	return ""
}

// validateUpdateForm validates fields for updating a customer.
func (v *CustomerView) validateUpdateForm() string {
	// Validate DOB format if provided
	dob := strings.TrimSpace(v.form.GetValue("dob"))
	if dob != "" {
		if _, err := time.Parse("2006-01-02", dob); err != nil {
			return "DOB must be in yyyy-mm-dd format"
		}
	}

	return ""
}

// buildAddInput creates AddCustomerInput from form values.
func (v *CustomerView) buildAddInput() (*service.AddCustomerInput, error) {
	values := v.form.GetAllValues()

	input := &service.AddCustomerInput{
		FirstName:    strings.TrimSpace(values["first_name"]),
		LastName:     strings.TrimSpace(values["last_name"]),
		HouseName:    strings.TrimSpace(values["house_name"]),
		HouseNumber:  strings.TrimSpace(values["house_number"]),
		Postcode:     strings.TrimSpace(values["postcode"]),
		PhoneHome:    strings.TrimSpace(values["phone_home"]),
		PhoneMobile:  strings.TrimSpace(values["phone_mobile"]),
		EmailAddress: strings.TrimSpace(values["email_address"]),
	}

	// Parse date of birth if provided
	dob := strings.TrimSpace(values["dob"])
	if dob != "" {
		t, err := time.Parse("2006-01-02", dob)
		if err != nil {
			return nil, errors.New("DOB must be in yyyy-mm-dd format")
		}
		input.DateOfBirth = &t
	}

	return input, nil
}

// buildUpdateInput creates UpdateCustomerInput from form values.
func (v *CustomerView) buildUpdateInput() *service.UpdateCustomerInput {
	values := v.form.GetAllValues()

	input := &service.UpdateCustomerInput{}

	// Only set fields that have values (to support partial updates)
	if val := strings.TrimSpace(values["first_name"]); val != "" {
		input.FirstName = &val
	}
	if val := strings.TrimSpace(values["last_name"]); val != "" {
		input.LastName = &val
	}
	if val := strings.TrimSpace(values["house_name"]); val != "" {
		input.HouseName = &val
	}
	if val := strings.TrimSpace(values["house_number"]); val != "" {
		input.HouseNumber = &val
	}
	if val := strings.TrimSpace(values["postcode"]); val != "" {
		input.Postcode = &val
	}
	if val := strings.TrimSpace(values["phone_home"]); val != "" {
		input.PhoneHome = &val
	}
	if val := strings.TrimSpace(values["phone_mobile"]); val != "" {
		input.PhoneMobile = &val
	}
	if val := strings.TrimSpace(values["email_address"]); val != "" {
		input.EmailAddress = &val
	}

	// Parse date of birth if provided
	if dob := strings.TrimSpace(values["dob"]); dob != "" {
		if t, err := time.Parse("2006-01-02", dob); err == nil {
			input.DateOfBirth = &t
		}
	}

	return input
}

// populateFormFromCustomer fills form fields with customer data.
func (v *CustomerView) populateFormFromCustomer(c *models.Customer) {
	v.form.SetValue("customer_num", components.FormatCustomerNum(c.CustomerNum))
	v.form.SetValue("first_name", c.GetFirstName())
	v.form.SetValue("last_name", c.GetLastName())

	// Format date of birth
	if c.DateOfBirth.Valid {
		v.form.SetValue("dob", c.DateOfBirth.Time.Format("2006-01-02"))
	} else {
		v.form.SetValue("dob", "")
	}

	v.form.SetValue("house_name", c.GetHouseName())
	v.form.SetValue("house_number", c.GetHouseNumber())
	v.form.SetValue("postcode", c.GetPostcode())
	v.form.SetValue("phone_home", c.GetPhoneHome())
	v.form.SetValue("phone_mobile", c.GetPhoneMobile())
	v.form.SetValue("email_address", c.GetEmailAddress())
}

// HandleKey handles key events specific to the customer view.
func (v *CustomerView) HandleKey(event *tcell.EventKey) *tcell.EventKey {
	// Handle F-keys for navigation to policy screens
	switch event.Key() {
	case tcell.KeyF1:
		if v.onNavigate != nil {
			v.onNavigate("motor")
		}
		return nil
	case tcell.KeyF2:
		if v.onNavigate != nil {
			v.onNavigate("endowment")
		}
		return nil
	case tcell.KeyF4:
		if v.onNavigate != nil {
			v.onNavigate("house")
		}
		return nil
	case tcell.KeyF5:
		if v.onNavigate != nil {
			v.onNavigate("commercial")
		}
		return nil
	}

	return v.BaseView.HandleKey(event)
}

// SetFocus sets focus to the customer number field.
func (v *CustomerView) SetFocus(app *tview.Application) {
	v.app = app
	v.form.SetFocus(app)
}

// Clear resets all form fields and the menu selection.
func (v *CustomerView) Clear() {
	v.form.Clear()
	v.menu.Clear()
	v.ClearError()
}

// GetCustomerNumber returns the current customer number value.
func (v *CustomerView) GetCustomerNumber() string {
	return v.form.GetValue("customer_num")
}

// SetCustomerNumber sets the customer number field.
func (v *CustomerView) SetCustomerNumber(num string) {
	v.form.SetValue("customer_num", components.FormatCustomerNum(num))
}

// truncateError truncates an error message to a maximum length.
func truncateError(msg string, maxLen int) string {
	if len(msg) <= maxLen {
		return msg
	}
	return msg[:maxLen-3] + "..."
}

// extractValidationMessage extracts the validation details from an error.
func extractValidationMessage(err error) string {
	msg := err.Error()
	// Look for the part after "validation failed: "
	if idx := strings.Index(msg, ": "); idx >= 0 {
		return msg[idx+2:]
	}
	return msg
}
