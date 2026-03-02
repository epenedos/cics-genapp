// Package views provides the screen implementations for the GENAPP application.
package views

import (
	"github.com/cicsdev/genapp/internal/service"
	"github.com/cicsdev/genapp/internal/ui/components"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// BaseView provides common functionality for all views.
type BaseView struct {
	name     string
	screen   *components.Screen
	app      *tview.Application
	onSubmit func()
	onCancel func()

	// Service references
	customerService *service.CustomerService
	policyService   *service.PolicyService
	counterService  *service.CounterService
}

// NewBaseView creates a new base view.
func NewBaseView(name, screenID, title string) *BaseView {
	return &BaseView{
		name:   name,
		screen: components.NewScreen(screenID, title),
	}
}

// Name returns the view name.
func (v *BaseView) Name() string {
	return v.name
}

// Layout returns the screen layout.
func (v *BaseView) Layout() tview.Primitive {
	return v.screen.Layout()
}

// SetFocus sets focus to the primary input.
func (v *BaseView) SetFocus(app *tview.Application) {
	v.app = app
	v.screen.SetFocus(app)
}

// HandleKey handles key events.
func (v *BaseView) HandleKey(event *tcell.EventKey) *tcell.EventKey {
	return v.screen.HandleKey(event)
}

// Clear resets all form fields.
func (v *BaseView) Clear() {
	v.screen.Clear()
}

// ShowError displays an error message.
func (v *BaseView) ShowError(msg string) {
	v.screen.ShowError(msg)
}

// ClearError clears the error message.
func (v *BaseView) ClearError() {
	v.screen.ClearError()
}

// ShowSuccess displays a success message.
func (v *BaseView) ShowSuccess(msg string) {
	v.screen.ShowSuccess(msg)
}

// Screen returns the underlying screen component.
func (v *BaseView) Screen() *components.Screen {
	return v.screen
}

// SetServices sets the service references for data operations.
func (v *BaseView) SetServices(customer *service.CustomerService, policy *service.PolicyService, counter *service.CounterService) {
	v.customerService = customer
	v.policyService = policy
	v.counterService = counter
}

// SetOnSubmit sets the submit handler.
func (v *BaseView) SetOnSubmit(fn func()) {
	v.onSubmit = fn
	v.screen.SetOnEnter(fn)
}

// SetOnCancel sets the cancel handler.
func (v *BaseView) SetOnCancel(fn func()) {
	v.onCancel = fn
	v.screen.SetOnEscape(fn)
}

// App returns the tview application reference.
func (v *BaseView) App() *tview.Application {
	return v.app
}

// CustomerService returns the customer service.
func (v *BaseView) CustomerService() *service.CustomerService {
	return v.customerService
}

// PolicyService returns the policy service.
func (v *BaseView) PolicyService() *service.PolicyService {
	return v.policyService
}

// CounterService returns the counter service.
func (v *BaseView) CounterService() *service.CounterService {
	return v.counterService
}
