// Package plan provides the /plans APIs
package plan

import (
	"net/http"

	stripe "github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/form"
)

// Client is used to invoke /plans APIs.
type Client struct {
	B   stripe.Backend
	Key string
}

// New creates a new plan.
func New(params *stripe.PlanParams) (*stripe.Plan, error) {
	return getC().New(params)
}

// New creates a new plan.
func (c Client) New(params *stripe.PlanParams) (*stripe.Plan, error) {
	plan := &stripe.Plan{}
	err := c.B.Call(http.MethodPost, "/v1/plans", c.Key, params, plan)
	return plan, err
}

// Get returns the details of a plan.
func Get(id string, params *stripe.PlanParams) (*stripe.Plan, error) {
	return getC().Get(id, params)
}

// Get returns the details of a plan.
func (c Client) Get(id string, params *stripe.PlanParams) (*stripe.Plan, error) {
	path := stripe.FormatURLPath("/v1/plans/%s", id)
	plan := &stripe.Plan{}
	err := c.B.Call(http.MethodGet, path, c.Key, params, plan)
	return plan, err
}

// Update updates a plan's properties.
func Update(id string, params *stripe.PlanParams) (*stripe.Plan, error) {
	return getC().Update(id, params)
}

// Update updates a plan's properties.
func (c Client) Update(id string, params *stripe.PlanParams) (*stripe.Plan, error) {
	path := stripe.FormatURLPath("/v1/plans/%s", id)
	plan := &stripe.Plan{}
	err := c.B.Call(http.MethodPost, path, c.Key, params, plan)
	return plan, err
}

// Del removes a plan.
func Del(id string, params *stripe.PlanParams) (*stripe.Plan, error) {
	return getC().Del(id, params)
}

// Del removes a plan.
func (c Client) Del(id string, params *stripe.PlanParams) (*stripe.Plan, error) {
	path := stripe.FormatURLPath("/v1/plans/%s", id)
	plan := &stripe.Plan{}
	err := c.B.Call(http.MethodDelete, path, c.Key, params, plan)
	return plan, err
}

// List returns a list of plans.
func List(params *stripe.PlanListParams) *Iter {
	return getC().List(params)
}

// List returns a list of plans.
func (c Client) List(listParams *stripe.PlanListParams) *Iter {
	return &Iter{stripe.GetIter(listParams, func(p *stripe.Params, b *form.Values) ([]interface{}, stripe.ListMeta, error) {
		list := &stripe.PlanList{}
		err := c.B.CallRaw(http.MethodGet, "/v1/plans", c.Key, b, p, list)

		ret := make([]interface{}, len(list.Data))
		for i, v := range list.Data {
			ret[i] = v
		}

		return ret, list.ListMeta, err
	})}
}

// Iter is an iterator for plans.
type Iter struct {
	*stripe.Iter
}

// Plan returns the plan which the iterator is currently pointing to.
func (i *Iter) Plan() *stripe.Plan {
	return i.Current().(*stripe.Plan)
}

func getC() Client {
	return Client{stripe.GetBackend(stripe.APIBackend), stripe.Key}
}
