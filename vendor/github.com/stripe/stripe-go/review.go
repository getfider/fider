package stripe

import "encoding/json"

// ReviewReasonType describes the reason why the review is open or closed.
type ReviewReasonType string

// List of values that ReviewReasonType can take.
const (
	ReviewReasonApproved        ReviewReasonType = "approved"
	ReviewReasonDisputed        ReviewReasonType = "disputed"
	ReviewReasonManual          ReviewReasonType = "manual"
	ReviewReasonRefunded        ReviewReasonType = "refunded"
	ReviewReasonRefundedAsFraud ReviewReasonType = "refunded_as_fraud"
	ReviewReasonRule            ReviewReasonType = "rule"
)

// ReviewParams is the set of parameters that can be used when approving a review.
type ReviewParams struct {
	Params `form:"*"`
}

// ReviewApproveParams is the set of parameters that can be used when approving a review.
type ReviewApproveParams struct {
	Params `form:"*"`
}

// ReviewListParams is the set of parameters that can be used when listing reviews.
type ReviewListParams struct {
	ListParams   `form:"*"`
	Created      *int64            `form:"created"`
	CreatedRange *RangeQueryParams `form:"created"`
}

// Review is the resource representing a Radar review.
// For more details see https://stripe.com/docs/api#reviews.
type Review struct {
	Charge        *Charge          `json:"charge"`
	Created       int64            `json:"created"`
	ID            string           `json:"id"`
	Livemode      bool             `json:"livemode"`
	Object        string           `json:"object"`
	Open          bool             `json:"open"`
	PaymentIntent *PaymentIntent   `json:"payment_intent"`
	Reason        ReviewReasonType `json:"reason"`
}

// ReviewList is a list of reviews as retrieved from a list endpoint.
type ReviewList struct {
	ListMeta
	Data []*Review `json:"data"`
}

// UnmarshalJSON handles deserialization of a Review.
// This custom unmarshaling is needed because the resulting
// property may be an id or the full struct if it was expanded.
func (r *Review) UnmarshalJSON(data []byte) error {
	if id, ok := ParseID(data); ok {
		r.ID = id
		return nil
	}

	type review Review
	var v review
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	*r = Review(v)
	return nil
}
