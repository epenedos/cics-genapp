package views

import (
	"github.com/cicsdev/genapp/internal/ui/components"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// MotorPolicyView implements the motor policy screen (SSMAPP1 equivalent).
// Placeholder implementation - full implementation in a later step.
type MotorPolicyView struct {
	*BaseView
	menu       *components.Menu
	form       *components.Form
	onNavigate func(screen string)
}

// NewMotorPolicyView creates a new motor policy view.
func NewMotorPolicyView() *MotorPolicyView {
	v := &MotorPolicyView{
		BaseView: NewBaseView("motor", "SSP1", "General Insurance Motor Policy Menu"),
	}

	// Create the menu
	v.menu = components.PolicyMenu()

	// Create the form with motor policy fields
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
		Label:     "Car Make",
		Name:      "car_make",
		MaxLength: 20,
		FieldType: components.FieldTypeText,
		Editable:  true,
	})
	v.form.AddField(&components.FormField{
		Label:     "Car Model",
		Name:      "car_model",
		MaxLength: 20,
		FieldType: components.FieldTypeText,
		Editable:  true,
	})
	v.form.AddField(&components.FormField{
		Label:        "Car Value",
		Name:         "car_value",
		MaxLength:    6,
		FieldType:    components.FieldTypeNumeric,
		Editable:     true,
		RightJustify: true,
		ZeroPad:      true,
	})
	v.form.AddField(&components.FormField{
		Label:     "Registration",
		Name:      "registration",
		MaxLength: 7,
		FieldType: components.FieldTypeText,
		Editable:  true,
	})
	v.form.AddField(&components.FormField{
		Label:     "Car Colour",
		Name:      "car_colour",
		MaxLength: 8,
		FieldType: components.FieldTypeText,
		Editable:  true,
	})
	v.form.AddField(&components.FormField{
		Label:        "CC",
		Name:         "cc",
		MaxLength:    8,
		FieldType:    components.FieldTypeNumeric,
		Editable:     true,
		RightJustify: true,
		ZeroPad:      true,
	})
	v.form.AddField(&components.FormField{
		Label:     "Manufacture Date",
		Name:      "manufactured",
		MaxLength: 10,
		FieldType: components.FieldTypeDate,
		Editable:  true,
	})
	v.form.AddField(&components.FormField{
		Label:        "No. Accidents",
		Name:         "accidents",
		MaxLength:    6,
		FieldType:    components.FieldTypeNumeric,
		Editable:     true,
		RightJustify: true,
		ZeroPad:      true,
	})
	v.form.AddField(&components.FormField{
		Label:        "Policy Premium",
		Name:         "premium",
		MaxLength:    6,
		FieldType:    components.FieldTypeNumeric,
		Editable:     true,
		RightJustify: true,
		ZeroPad:      true,
	})

	v.screen.SetMenu(v.menu)
	v.screen.SetForm(v.form)
	v.SetOnSubmit(v.handleSubmit)

	return v
}

func (v *MotorPolicyView) handleSubmit() {
	v.ShowSuccess("Motor policy operation (placeholder)")
}

func (v *MotorPolicyView) SetOnNavigate(fn func(screen string)) {
	v.onNavigate = fn
}

func (v *MotorPolicyView) HandleKey(event *tcell.EventKey) *tcell.EventKey {
	return v.BaseView.HandleKey(event)
}

func (v *MotorPolicyView) SetFocus(app *tview.Application) {
	v.app = app
	v.form.SetFocus(app)
}

// EndowmentPolicyView implements the endowment policy screen (SSMAPP2 equivalent).
type EndowmentPolicyView struct {
	*BaseView
	menu       *components.Menu
	form       *components.Form
	onNavigate func(screen string)
}

// NewEndowmentPolicyView creates a new endowment policy view.
func NewEndowmentPolicyView() *EndowmentPolicyView {
	v := &EndowmentPolicyView{
		BaseView: NewBaseView("endowment", "SSP2", "General Insurance Endowment Policy Menu"),
	}

	v.menu = components.PolicyMenu()

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
		Label:     "Term",
		Name:      "term",
		MaxLength: 2,
		FieldType: components.FieldTypeNumeric,
		Editable:  true,
	})
	v.form.AddField(&components.FormField{
		Label:        "Sum Assured",
		Name:         "sum_assured",
		MaxLength:    6,
		FieldType:    components.FieldTypeNumeric,
		Editable:     true,
		RightJustify: true,
		ZeroPad:      true,
	})
	v.form.AddField(&components.FormField{
		Label:     "Life Assured",
		Name:      "life_assured",
		MaxLength: 25,
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

	v.screen.SetMenu(v.menu)
	v.screen.SetForm(v.form)
	v.SetOnSubmit(v.handleSubmit)

	return v
}

func (v *EndowmentPolicyView) handleSubmit() {
	v.ShowSuccess("Endowment policy operation (placeholder)")
}

func (v *EndowmentPolicyView) SetOnNavigate(fn func(screen string)) {
	v.onNavigate = fn
}

func (v *EndowmentPolicyView) HandleKey(event *tcell.EventKey) *tcell.EventKey {
	return v.BaseView.HandleKey(event)
}

func (v *EndowmentPolicyView) SetFocus(app *tview.Application) {
	v.app = app
	v.form.SetFocus(app)
}

// HousePolicyView implements the house policy screen (SSMAPP3 equivalent).
type HousePolicyView struct {
	*BaseView
	menu       *components.Menu
	form       *components.Form
	onNavigate func(screen string)
}

// NewHousePolicyView creates a new house policy view.
func NewHousePolicyView() *HousePolicyView {
	v := &HousePolicyView{
		BaseView: NewBaseView("house", "SSP3", "General Insurance House Policy Menu"),
	}

	v.menu = components.PolicyMenu()

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
		ZeroPad:      true,
	})
	v.form.AddField(&components.FormField{
		Label:        "House Value",
		Name:         "house_value",
		MaxLength:    8,
		FieldType:    components.FieldTypeNumeric,
		Editable:     true,
		RightJustify: true,
		ZeroPad:      true,
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

	v.screen.SetMenu(v.menu)
	v.screen.SetForm(v.form)
	v.SetOnSubmit(v.handleSubmit)

	return v
}

func (v *HousePolicyView) handleSubmit() {
	v.ShowSuccess("House policy operation (placeholder)")
}

func (v *HousePolicyView) SetOnNavigate(fn func(screen string)) {
	v.onNavigate = fn
}

func (v *HousePolicyView) HandleKey(event *tcell.EventKey) *tcell.EventKey {
	return v.BaseView.HandleKey(event)
}

func (v *HousePolicyView) SetFocus(app *tview.Application) {
	v.app = app
	v.form.SetFocus(app)
}

// CommercialPolicyView implements the commercial policy screen (SSMAPP4 equivalent).
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
	return v.BaseView.HandleKey(event)
}

func (v *CommercialPolicyView) SetFocus(app *tview.Application) {
	v.app = app
	v.form.SetFocus(app)
}

// ClaimView implements the claim screen (SSMAPP5 equivalent).
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
	return v.BaseView.HandleKey(event)
}

func (v *ClaimView) SetFocus(app *tview.Application) {
	v.app = app
	v.form.SetFocus(app)
}
