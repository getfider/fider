package stripe

// RadarValueListItemType is the possible values for a type of value list items.
type RadarValueListItemType string

// List of values that RadarValueListItemType can take.
const (
	RadarValueListItemTypeCardBin             RadarValueListItemType = "card_bin"
	RadarValueListItemTypeCardFingerprint     RadarValueListItemType = "card_fingerprint"
	RadarValueListItemTypeCountry             RadarValueListItemType = "country"
	RadarValueListItemTypeEmail               RadarValueListItemType = "email"
	RadarValueListItemTypeIPAddress           RadarValueListItemType = "ip_address"
	RadarValueListItemTypeString              RadarValueListItemType = "string"
	RadarValueListItemTypeCaseSensitiveString RadarValueListItemType = "case_sensitive_string"
)

// RadarValueListParams is the set of parameters that can be used when creating a value list.
type RadarValueListParams struct {
	Params   `form:"*"`
	Alias    *string `form:"alias"`
	ItemType *string `form:"item_type"`
	Name     *string `form:"name"`
}

// RadarValueListListParams is the set of parameters that can be used when listing value lists.
type RadarValueListListParams struct {
	ListParams   `form:"*"`
	Alias        *string           `form:"alias"`
	Contains     *string           `form:"contains"`
	Created      *int64            `form:"created"`
	CreatedRange *RangeQueryParams `form:"created"`
}

// RadarValueList is the resource representing a value list.
type RadarValueList struct {
	Alias     string                  `json:"alias"`
	Created   int64                   `json:"created"`
	CreatedBy string                  `json:"created_by"`
	Deleted   bool                    `json:"deleted"`
	ID        string                  `json:"id"`
	ItemType  RadarValueListItemType  `json:"item_type"`
	ListItems *RadarValueListItemList `json:"list_items"`
	Livemode  bool                    `json:"livemode"`
	Metadata  map[string]string       `json:"metadata"`
	Name      string                  `json:"name"`
	Object    string                  `json:"object"`
	Updated   int64                   `json:"updated"`
	UpdatedBy string                  `json:"updated_by"`
}

// RadarValueListList is a list of value lists as retrieved from a list endpoint.
type RadarValueListList struct {
	ListMeta
	Data []*RadarValueList `json:"data"`
}
