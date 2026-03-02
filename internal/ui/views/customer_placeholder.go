package views

import (
	"github.com/cicsdev/genapp/internal/ui/components"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// CustomerView implements the customer menu screen (SSMAPC1 equivalent).
// This is a placeholder that establishes the layout - full implementation
// will connect to services in a later step.
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

	// Validate form if needed
	if valid, errMsg := v.form.Validate(); !valid {
		v.ShowError(errMsg)
		return
	}

	// Get form values
	values := v.form.GetAllValues()
	customerNum := values["customer_num"]

	// Process based on option
	switch option.Key {
	case "1": // Inquiry
		if customerNum == "" {
			v.ShowError("Customer Number required for inquiry")
			return
		}
		// TODO: Call customer service to get customer
		v.ShowSuccess("Inquiry completed (placeholder)")

	case "2": // Add
		// TODO: Call customer service to add customer
		v.ShowSuccess("Customer added (placeholder)")

	case "4": // Update
		if customerNum == "" {
			v.ShowError("Customer Number required for update")
			return
		}
		// TODO: Call customer service to update customer
		v.ShowSuccess("Customer updated (placeholder)")
	}
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
