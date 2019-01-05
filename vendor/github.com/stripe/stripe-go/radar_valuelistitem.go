package stripe

// RadarValueListItemParams is the set of parameters that can be used when creating a value list item.
type RadarValueListItemParams struct {
	Params         `form:"*"`
	Value          *string `form:"value"`
	RadarValueList *string `form:"value_list"`
}

// RadarValueListItemListParams is the set of parameters that can be used when listing value list items.
type RadarValueListItemListParams struct {
	ListParams     `form:"*"`
	Created        *int64            `form:"created"`
	CreatedRange   *RangeQueryParams `form:"created"`
	RadarValueList *string           `form:"value_list"`
	Value          *string           `form:"value"`
}

// RadarValueListItem is the resource representing a value list item.
type RadarValueListItem struct {
	Created        int64  `json:"created"`
	CreatedBy      string `json:"created_by"`
	Deleted        bool   `json:"deleted"`
	ID             string `json:"id"`
	Livemode       bool   `json:"livemode"`
	Name           string `json:"name"`
	Object         string `json:"object"`
	Value          string `json:"value"`
	RadarValueList string `json:"value_list"`
}

// RadarValueListItemList is a list of value list items as retrieved from a list endpoint.
type RadarValueListItemList struct {
	ListMeta
	Data []*RadarValueListItem `json:"data"`
}
