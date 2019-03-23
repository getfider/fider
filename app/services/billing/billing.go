package billing

import (
	"github.com/getfider/fider/app/pkg/bus"
	"github.com/getfider/fider/app/pkg/env"
	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/client"
)

var stripeClient *client.API

func init() {
	bus.Register(Service{})
}

type Service struct{}

func (s Service) Name() string {
	return "Stripe"
}

func (s Service) Category() string {
	return "billing"
}

func (s Service) Enabled() bool {
	return env.IsBillingEnabled()
}

func (s Service) Init() {
	stripe.LogLevel = 0
	stripeClient = &client.API{}
	stripeClient.Init(env.Config.Stripe.SecretKey, nil)

	bus.AddHandler(listPlans)
}

/*
func (c *Client) CancelSubscription() error {
func (c *Client) ClearPaymentInfo() error {
func (c *Client) CreateCustomer(email string) (string, error) {
func (c *Client) DeleteCustomer() error {
func (c *Client) GetPaymentInfo() (*models.PaymentInfo, error) {
func (c *Client) GetPlanByID(countryCode, planID string) (*models.BillingPlan, error) {
func (c *Client) GetUpcomingInvoice() (*models.UpcomingInvoice, error) {
func (c *Client) ListPlans(countryCode string) ([]*models.BillingPlan, error) {
func (c *Client) Subscribe(planID string) error {
func (c *Client) UpdatePaymentInfo(input *models.CreateEditBillingPaymentInfo) error {
*/
