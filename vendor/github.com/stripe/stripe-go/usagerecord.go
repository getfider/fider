package stripe

// Possible values for the action parameter on usage record creation.
const (
	UsageRecordActionIncrement string = "increment"
	UsageRecordActionSet       string = "set"
)

// UsageRecord represents a usage record.
// See https://stripe.com/docs/api#usage_records
type UsageRecord struct {
	ID               string `json:"id"`
	Livemode         bool   `json:"livemode"`
	Quantity         int64  `json:"quantity"`
	SubscriptionItem string `json:"subscription_item"`
	Timestamp        int64  `json:"timestamp"`
}

// UsageRecordParams create a usage record for a specified subscription item
// and date, and fills it with a quantity.
type UsageRecordParams struct {
	Params           `form:"*"`
	Action           *string `form:"action"`
	Quantity         *int64  `form:"quantity"`
	SubscriptionItem *string `form:"-"` // passed in the URL
	Timestamp        *int64  `form:"timestamp"`
}
