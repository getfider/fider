package stripe

import "encoding/json"

// SigmaScheduledQueryRunStatus is the possible values for status for a scheduled query run.
type SigmaScheduledQueryRunStatus string

// List of values that SigmaScheduledQueryRunStatus can take.
const (
	SigmaScheduledQueryRunStatusCanceled  SigmaScheduledQueryRunStatus = "canceled"
	SigmaScheduledQueryRunStatusCompleted SigmaScheduledQueryRunStatus = "completed"
	SigmaScheduledQueryRunStatusFailed    SigmaScheduledQueryRunStatus = "failed"
	SigmaScheduledQueryRunStatusTimedOut  SigmaScheduledQueryRunStatus = "timed_out"
)

// SigmaScheduledQueryRunParams is the set of parameters that can be used when updating a scheduled query run.
type SigmaScheduledQueryRunParams struct {
	Params `form:"*"`
}

// SigmaScheduledQueryRunListParams is the set of parameters that can be used when listing scheduled query runs.
type SigmaScheduledQueryRunListParams struct {
	ListParams `form:"*"`
}

// SigmaScheduledQueryRun is the resource representing a scheduled query run.
type SigmaScheduledQueryRun struct {
	Created              int64                        `json:"created"`
	DataLoadTime         int64                        `json:"data_load_time"`
	Error                string                       `json:"error"`
	File                 *File                        `json:"file"`
	ID                   string                       `json:"id"`
	Livemode             bool                         `json:"livemode"`
	Object               string                       `json:"object"`
	ResultAvailableUntil int64                        `json:"result_available_until"`
	SQL                  string                       `json:"sql"`
	Status               SigmaScheduledQueryRunStatus `json:"status"`
	Query                string                       `json:"query"`
}

// SigmaScheduledQueryRunList is a list of scheduled query runs as retrieved from a list endpoint.
type SigmaScheduledQueryRunList struct {
	ListMeta
	Data []*SigmaScheduledQueryRun `json:"data"`
}

// UnmarshalJSON handles deserialization of an SigmaScheduledQueryRun.
// This custom unmarshaling is needed because the resulting
// property may be an id or the full struct if it was expanded.
func (i *SigmaScheduledQueryRun) UnmarshalJSON(data []byte) error {
	if id, ok := ParseID(data); ok {
		i.ID = id
		return nil
	}

	type sigmaScheduledQueryRun SigmaScheduledQueryRun
	var v sigmaScheduledQueryRun
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	*i = SigmaScheduledQueryRun(v)
	return nil
}
