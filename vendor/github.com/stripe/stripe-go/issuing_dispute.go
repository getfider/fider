package stripe

import "encoding/json"

// IssuingDisputeReason is the list of possible values for status on an issuing dispute.
type IssuingDisputeReason string

// List of values that IssuingDisputeReason can take.
const (
	IssuingDisputeReasonFraudulent IssuingDisputeReason = "fraudulent"
	IssuingDisputeReasonOther      IssuingDisputeReason = "other"
)

// IssuingDisputeStatus is the list of possible values for status on an issuing dispute.
type IssuingDisputeStatus string

// List of values that IssuingDisputeStatus can take.
const (
	IssuingDisputeStatusLost        IssuingDisputeStatus = "lost"
	IssuingDisputeStatusUnderReview IssuingDisputeStatus = "under_review"
	IssuingDisputeStatusUnsubmitted IssuingDisputeStatus = "unsubmitted"
	IssuingDisputeStatusWon         IssuingDisputeStatus = "won"
)

// IssuingDisputeEvidenceFraudulentParams is the subset of parameters that can be sent as evidence for an issuing dispute
// with the reason set as fraudulent.
type IssuingDisputeEvidenceFraudulentParams struct {
	DisputeExplanation *string `form:"dispute_explanation"`
	UncategorizedFile  *string `form:"uncategorized_file"`
}

// IssuingDisputeEvidenceOtherParams is the subset of parameters that can be sent as evidence for an issuing dispute
// with the reason set as other.
type IssuingDisputeEvidenceOtherParams struct {
	DisputeExplanation *string `form:"dispute_explanation"`
	UncategorizedFile  *string `form:"uncategorized_file"`
}

// IssuingDisputeEvidenceParams is the set of parameters that can be sent as evidence for an issuing dispute.
type IssuingDisputeEvidenceParams struct {
	Fraudulent *IssuingDisputeEvidenceFraudulentParams `form:"fraudulent"`
	Other      *IssuingDisputeEvidenceOtherParams      `form:"other"`
}

// IssuingDisputeParams is the set of parameters that can be used when creating or updating an issuing dispute.
type IssuingDisputeParams struct {
	Params              `form:"*"`
	Amount              *int64                        `form:"amount"`
	Evidence            *IssuingDisputeEvidenceParams `form:"evidence"`
	Reason              *string                       `form:"reason"`
	DisputedTransaction *string                       `form:"disputed_transaction"`
}

// IssuingDisputeListParams is the set of parameters that can be used when listing issuing dispute.
type IssuingDisputeListParams struct {
	ListParams   `form:"*"`
	Created      *int64            `form:"created"`
	CreatedRange *RangeQueryParams `form:"created"`
	Transaction  *string           `form:"transaction"`
}

// IssuingDisputeEvidenceFraudulent is the resource representing the evidence hash on an issuing dispute
// with the reason set as fraudulent.
type IssuingDisputeEvidenceFraudulent struct {
	DisputeExplanation string `json:"dispute_explanation"`
	UncategorizedFile  *File  `json:"uncategorized_file"`
}

// IssuingDisputeEvidenceOther is the resource representing the evidence hash on an issuing dispute
// with the reason set as other.
type IssuingDisputeEvidenceOther struct {
	DisputeExplanation string `json:"dispute_explanation"`
	UncategorizedFile  *File  `json:"uncategorized_file"`
}

// IssuingDisputeEvidence is the resource representing evidence on an issuing dispute.
type IssuingDisputeEvidence struct {
	Fraudulent *IssuingDisputeEvidenceFraudulent `json:"fraudulent"`
	Other      *IssuingDisputeEvidenceOther      `json:"other"`
}

// IssuingDispute is the resource representing an issuing dispute.
type IssuingDispute struct {
	Amount      int64                   `json:"amount"`
	Created     int64                   `json:"created"`
	Currency    Currency                `json:"currency"`
	Evidence    *IssuingDisputeEvidence `json:"evidence"`
	ID          string                  `json:"id"`
	Livemode    bool                    `json:"livemode"`
	Metadata    map[string]string       `json:"metadata"`
	Object      string                  `json:"object"`
	Reason      IssuingDisputeReason    `json:"reason"`
	Status      IssuingDisputeStatus    `json:"status"`
	Transaction *IssuingTransaction     `json:"transaction"`
}

// IssuingDisputeList is a list of issuing disputes as retrieved from a list endpoint.
type IssuingDisputeList struct {
	ListMeta
	Data []*IssuingDispute `json:"data"`
}

// UnmarshalJSON handles deserialization of an IssuingDispute.
// This custom unmarshaling is needed because the resulting
// property may be an id or the full struct if it was expanded.
func (i *IssuingDispute) UnmarshalJSON(data []byte) error {
	if id, ok := ParseID(data); ok {
		i.ID = id
		return nil
	}

	type issuingDispute IssuingDispute
	var v issuingDispute
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	*i = IssuingDispute(v)
	return nil
}
