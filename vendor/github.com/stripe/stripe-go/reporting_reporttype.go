package stripe

// ReportTypeListParams is the set of parameters that can be used when listing report types.
type ReportTypeListParams struct {
	ListParams `form:"*"`
}

// ReportTypeParams is the set of parameters that can be used when retrieving a report type.
type ReportTypeParams struct {
	Params `form:"*"`
}

// ReportType is the resource representing a report type.
type ReportType struct {
	Created            int64  `json:"created"`
	DataAvailableEnd   int64  `json:"data_available_end"`
	DataAvailableStart int64  `json:"data_available_start"`
	ID                 string `json:"id"`
	Name               string `json:"name"`
	Object             string `json:"object"`
	Updated            int64  `json:"updated"`
	Version            int64  `json:"version"`
}

// ReportTypeList is a list of report types as retrieved from a list endpoint.
type ReportTypeList struct {
	ListMeta
	Data []*ReportType `json:"data"`
}
