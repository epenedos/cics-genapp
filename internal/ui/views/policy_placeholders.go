package views

import (
	"github.com/cicsdev/genapp/internal/ui/components"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// CommercialPolicyView implements the commercial policy screen (SSMAPP4 equivalent).
// Placeholder implementation - full implementation in a later step.
type CommercialPolicyView struct {
	*BaseView
	menu       *components.Menu
	form       *components.Form
	onNavigate func(screen string)
}

// NewCommercialPolicyView creates a new commercial policy view.
func NewCommercialPolicyView() *CommercialPolicyView {
	v := &CommercialPolicyView{
		BaseView: NewBaseView("commercial", "SSP4", "General Insurance Commercial Policy Menu"),
	}

	// Commercial has no update option
	v.menu = components.CommercialPolicyMenu()

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
		Label:     "Start date",
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
		Label:     "Address",
		Name:      "address",
		MaxLength: 25,
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
		Label:     "Latitude",
		Name:      "latitude",
		MaxLength: 11,
		FieldType: components.FieldTypeText,
		Editable:  true,
	})
	v.form.AddField(&components.FormField{
		Label:     "Longitude",
		Name:      "longitude",
		MaxLength: 11,
		FieldType: components.FieldTypeText,
		Editable:  true,
	})
	v.form.AddField(&components.FormField{
		Label:     "Customer Name",
		Name:      "customer_name",
		MaxLength: 25,
		FieldType: components.FieldTypeText,
		Editable:  true,
	})
	v.form.AddField(&components.FormField{
		Label:     "Property Type",
		Name:      "property_type",
		MaxLength: 25,
		FieldType: components.FieldTypeText,
		Editable:  true,
	})

	v.screen.SetMenu(v.menu)
	v.screen.SetForm(v.form)
	v.SetOnSubmit(v.handleSubmit)

	return v
}

func (v *CommercialPolicyView) handleSubmit() {
	v.ShowSuccess("Commercial policy operation (placeholder)")
}

func (v *CommercialPolicyView) SetOnNavigate(fn func(screen string)) {
	v.onNavigate = fn
}

func (v *CommercialPolicyView) HandleKey(event *tcell.EventKey) *tcell.EventKey {
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

func (v *CommercialPolicyView) SetFocus(app *tview.Application) {
	v.app = app
	v.form.SetFocus(app)
}

// Clear resets all form fields and the menu selection.
func (v *CommercialPolicyView) Clear() {
	v.form.Clear()
	v.menu.Clear()
	v.ClearError()
}

// SetCustomerNumber sets the customer number field.
func (v *CommercialPolicyView) SetCustomerNumber(num string) {
	v.form.SetValue("customer_num", components.FormatCustomerNum(num))
}

// ClaimView implements the claim screen (SSMAPP5 equivalent).
// Placeholder implementation - full implementation in a later step.
type ClaimView struct {
	*BaseView
	menu       *components.Menu
	form       *components.Form
	onNavigate func(screen string)
}

// NewClaimView creates a new claim view.
func NewClaimView() *ClaimView {
	v := &ClaimView{
		BaseView: NewBaseView("claim", "SSP5", "General Insurance Policy Claim Menu"),
	}

	v.menu = components.ClaimMenu()

	v.form = components.NewForm()
	v.form.AddField(&components.FormField{
		Label:        "Claim Number",
		Name:         "claim_num",
		MaxLength:    10,
		FieldType:    components.FieldTypeNumeric,
		Editable:     true,
		InitialFocus: true,
		RightJustify: true,
		ZeroPad:      true,
	})
	v.form.AddField(&components.FormField{
		Label:        "Policy Number",
		Name:         "policy_num",
		MaxLength:    10,
		FieldType:    components.FieldTypeNumeric,
		Editable:     true,
		RightJustify: true,
		ZeroPad:      true,
	})
	v.form.AddField(&components.FormField{
		Label:        "Customer Number",
		Name:         "customer_num",
		MaxLength:    10,
		FieldType:    components.FieldTypeNumeric,
		Editable:     true,
		RightJustify: true,
		ZeroPad:      true,
	})
	v.form.AddField(&components.FormField{
		Label:     "Claim date",
		Name:      "claim_date",
		MaxLength: 10,
		FieldType: components.FieldTypeDate,
		Editable:  true,
	})
	v.form.AddField(&components.FormField{
		Label:        "Paid",
		Name:         "paid",
		MaxLength:    10,
		FieldType:    components.FieldTypeNumeric,
		Editable:     true,
		RightJustify: true,
		ZeroPad:      true,
	})
	v.form.AddField(&components.FormField{
		Label:        "Value",
		Name:         "value",
		MaxLength:    10,
		FieldType:    components.FieldTypeNumeric,
		Editable:     true,
		RightJustify: true,
		ZeroPad:      true,
	})
	v.form.AddField(&components.FormField{
		Label:     "Cause",
		Name:      "cause",
		MaxLength: 25,
		FieldType: components.FieldTypeText,
		Editable:  true,
	})
	v.form.AddField(&components.FormField{
		Label:     "Observation",
		Name:      "observations",
		MaxLength: 25,
		FieldType: components.FieldTypeText,
		Editable:  true,
	})

	v.screen.SetMenu(v.menu)
	v.screen.SetForm(v.form)
	v.SetOnSubmit(v.handleSubmit)

	return v
}

func (v *ClaimView) handleSubmit() {
	v.ShowSuccess("Claim operation (placeholder)")
}

func (v *ClaimView) SetOnNavigate(fn func(screen string)) {
	v.onNavigate = fn
}

func (v *ClaimView) HandleKey(event *tcell.EventKey) *tcell.EventKey {
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

func (v *ClaimView) SetFocus(app *tview.Application) {
	v.app = app
	v.form.SetFocus(app)
}

// Clear resets all form fields and the menu selection.
func (v *ClaimView) Clear() {
	v.form.Clear()
	v.menu.Clear()
	v.ClearError()
}
