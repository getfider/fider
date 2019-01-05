package stripe

// IssuerFraudType are strings that map to the fraud label category from the issuer.
type IssuerFraudType string

// List of values that IssuerFraudType can take.
const (
	IssuerFraudTypeCardNeverReceived         IssuerFraudType = "card_never_received"
	IssuerFraudTypeFraudulentCardApplication IssuerFraudType = "fraudulent_card_application"
	IssuerFraudTypeMadeWithCounterfeitCard   IssuerFraudType = "made_with_counterfeit_card"
	IssuerFraudTypeMadeWithLostCard          IssuerFraudType = "made_with_lost_card"
	IssuerFraudTypeMadeWithStolenCard        IssuerFraudType = "made_with_stolen_card"
	IssuerFraudTypeMisc                      IssuerFraudType = "misc"
	IssuerFraudTypeUnauthorizedUseOfCard     IssuerFraudType = "unauthorized_use_of_card"
)

// IssuerFraudRecordParams is the set of parameters that can be used when
// retrieving issuer fraud records. For more details see
// https://stripe.com/docs#retrieve_issuer_fraud_records.
type IssuerFraudRecordParams struct {
	Params `form:"*"`
}

// IssuerFraudRecordListParams is the set of parameters that can be used when
// listing issuer fraud records. For more details see
// https://stripe.com/docs#list_issuer_fraud_records.
type IssuerFraudRecordListParams struct {
	ListParams `form:"*"`
	Charge     *string `form:"-"`
}

// IssuerFraudRecordList is a list of issuer fraud records as retrieved from a
// list endpoint.
type IssuerFraudRecordList struct {
	ListMeta
	Values []*IssuerFraudRecord `json:"data"`
}

// IssuerFraudRecord is the resource representing an issuer fraud record. For
// more details see https://stripe.com/docs/api#issuer_fraud_records.
type IssuerFraudRecord struct {
	Actionable        bool            `json:"actionable"`
	Charge            *Charge         `json:"charge"`
	Created           int64           `json:"created"`
	FraudType         IssuerFraudType `json:"fraud_type"`
	HasLiabilityShift bool            `json:"has_liability_shift"`
	ID                string          `json:"id"`
	Livemode          bool            `json:"livemode"`
	PostDate          int64           `json:"post_date"`
}
