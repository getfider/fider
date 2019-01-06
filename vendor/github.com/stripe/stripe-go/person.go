package stripe

import "encoding/json"

// RelationshipParams is used to set the relationship between an account and a person.
type RelationshipParams struct {
	AccountOpener    *bool    `form:"account_opener"`
	Director         *bool    `form:"director"`
	Executive        *bool    `form:"executive"`
	Owner            *bool    `form:"owner"`
	PercentOwnership *float64 `form:"percent_ownership"`
	Title            *string  `form:"title"`
}

// PersonParams is the set of parameters that can be used when creating or updating a person.
// For more details see https://stripe.com/docs/api#create_person.
type PersonParams struct {
	Params         `form:"*"`
	Account        *string               `form:"-"` // Included in URL
	Address        *AccountAddressParams `form:"address"`
	AddressKana    *AccountAddressParams `form:"address_kana"`
	AddressKanji   *AccountAddressParams `form:"address_kanji"`
	DOB            *DOBParams            `form:"dob"`
	Email          *string               `form:"email"`
	FirstName      *string               `form:"first_name"`
	FirstNameKana  *string               `form:"first_name_kana"`
	FirstNameKanji *string               `form:"first_name_kanji"`
	Gender         *string               `form:"gender"`
	IDNumber       *string               `form:"id_number"`
	LastName       *string               `form:"last_name"`
	LastNameKana   *string               `form:"last_name_kana"`
	LastNameKanji  *string               `form:"last_name_kanji"`
	MaidenName     *string               `form:"maiden_name"`
	Phone          *string               `form:"phone"`
	Relationship   *RelationshipParams   `form:"relationship"`
	SSNLast4       *string               `form:"ssn_last_4"`
}

// RelationshipListParams is used to filter persons by the relationship
type RelationshipListParams struct {
	AccountOpener *bool `form:"account_opener"`
	Director      *bool `form:"director"`
	Executive     *bool `form:"executive"`
	Owner         *bool `form:"owner"`
}

// PersonListParams is the set of parameters that can be used when listing persons.
// For more detail see https://stripe.com/docs/api#list_persons.
type PersonListParams struct {
	ListParams   `form:"*"`
	Account      *string                 `form:"-"` // Included in URL
	Relationship *RelationshipListParams `form:"relationship"`
}

// Relationship represents extra information needed for a Person.
type Relationship struct {
	AccountOpener    bool    `json:"account_opener"`
	Director         bool    `json:"director"`
	Executive        bool    `json:"executive"`
	Owner            bool    `json:"owner"`
	PercentOwnership float64 `json:"percent_ownership"`
	Title            string  `json:"title"`
}

// Requirements represents what's missing to verify a Person.
type Requirements struct {
	CurrentlyDue  []string `json:"currently_due"`
	EventuallyDue []string `json:"eventually_due"`
	PastDue       []string `json:"past_due"`
}

// Person is the resource representing a Stripe person.
// For more details see https://stripe.com/docs/api#persons.
type Person struct {
	Account          string                `json:"account"`
	Address          *AccountAddress       `json:"address"`
	AddressKana      *AccountAddress       `json:"address_kana"`
	AddressKanji     *AccountAddress       `json:"address_kanji"`
	Deleted          bool                  `json:"deleted"`
	DOB              *DOB                  `json:"dob"`
	Email            string                `json:"email"`
	FirstName        string                `json:"first_name"`
	FirstNameKana    string                `json:"first_name_kana"`
	FirstNameKanji   string                `json:"first_name_kanji"`
	Gender           string                `json:"gender"`
	ID               string                `json:"id"`
	IDNumberProvided bool                  `json:"id_number_provided"`
	LastName         string                `json:"last_name"`
	LastNameKana     string                `json:"last_name_kana"`
	LastNameKanji    string                `json:"last_name_kanji"`
	MaidenName       string                `json:"maiden_name"`
	Metadata         map[string]string     `json:"metadata"`
	Object           string                `json:"object"`
	Phone            string                `json:"phone"`
	Relationship     *Relationship         `json:"relationship"`
	Requirements     *Requirements         `json:"requirements"`
	SSNLast4Provided bool                  `json:"ssn_last_4_provided"`
	Verification     *IdentityVerification `json:"verification"`
}

// PersonList is a list of persons as retrieved from a list endpoint.
type PersonList struct {
	ListMeta
	Data []*Person `json:"data"`
}

// UnmarshalJSON handles deserialization of a Person.
// This custom unmarshaling is needed because the resulting
// property may be an id or the full struct if it was expanded.
func (c *Person) UnmarshalJSON(data []byte) error {
	if id, ok := ParseID(data); ok {
		c.ID = id
		return nil
	}

	type person Person
	var v person
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	*c = Person(v)
	return nil
}
