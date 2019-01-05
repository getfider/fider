package stripe

// ReportRunStatus is the possible values for status on a report run.
type ReportRunStatus string

// List of values that ReportRunStatus can take.
const (
	ReportRunStatusFailed    ReportRunStatus = "failed"
	ReportRunStatusPending   ReportRunStatus = "pending"
	ReportRunStatusSucceeded ReportRunStatus = "succeeded"
)

// ReportRunParametersParams is the set of parameters that can be used when creating a report run.
type ReportRunParametersParams struct {
	ConnectedAccount  *string `form:"connected_account"`
	Currency          *string `form:"currency"`
	IntervalEnd       *int64  `form:"interval_end"`
	IntervalStart     *int64  `form:"interval_start"`
	Payout            *string `form:"payout"`
	ReportingCategory *string `form:"reporting_category"`
}

// ReportRunParams is the set of parameters that can be used when creating a report run.
type ReportRunParams struct {
	Params     `form:"*"`
	Parameters *ReportRunParametersParams `form:"parameters"`
	ReportType *string                    `form:"report_type"`
}

// ReportRunListParams is the set of parameters that can be used when listing report runs.
type ReportRunListParams struct {
	ListParams   `form:"*"`
	Created      *int64            `form:"created"`
	CreatedRange *RangeQueryParams `form:"created"`
}

// ReportRunParameters describes the parameters hash on a report run.
type ReportRunParameters struct {
	ConnectedAccount  string   `json:"connected_account"`
	Currency          Currency `json:"currency"`
	IntervalEnd       int64    `json:"interval_end"`
	IntervalStart     int64    `json:"interval_start"`
	Payout            string   `json:"payout"`
	ReportingCategory string   `json:"reporting_category"`
}

// ReportRun is the resource representing a report run.
type ReportRun struct {
	Created     int64                `json:"created"`
	Error       string               `json:"error"`
	ID          string               `json:"id"`
	Livemode    bool                 `json:"livemode"`
	Object      string               `json:"object"`
	Parameters  *ReportRunParameters `json:"parameters"`
	ReportType  string               `json:"report_type"`
	Result      *File                `json:"result"`
	Status      ReportRunStatus      `json:"status"`
	SucceededAt int64                `json:"succeeded_at"`
}

// ReportRunList is a list of report runs as retrieved from a list endpoint.
type ReportRunList struct {
	ListMeta
	Data []*ReportRun `json:"data"`
}
