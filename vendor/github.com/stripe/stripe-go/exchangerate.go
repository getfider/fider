package stripe

// ExchangeRate is the resource representing the currency exchange rates at
// a given time.
type ExchangeRate struct {
	ID    string               `json:"id"`
	Rates map[Currency]float64 `json:"rates"`
}

// ExchangeRateParams is the set of parameters that can be used when retrieving
// exchange rates.
type ExchangeRateParams struct {
	Params `form:"*"`
}

// ExchangeRateList is a list of exchange rates as retrieved from a list endpoint.
type ExchangeRateList struct {
	ListMeta
	Data []*ExchangeRate `json:"data"`
}

// ExchangeRateListParams are the parameters allowed during ExchangeRate listing.
type ExchangeRateListParams struct {
	ListParams `form:"*"`
}
