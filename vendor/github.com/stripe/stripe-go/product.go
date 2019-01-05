package stripe

import "encoding/json"

// ProductType is the type of a product.
type ProductType string

// List of values that ProductType can take.
const (
	ProductTypeGood    ProductType = "good"
	ProductTypeService ProductType = "service"
)

// PackageDimensionsParams represents the set of parameters for the the dimension of a
// product or a SKU from the perspective of shipping .
type PackageDimensionsParams struct {
	Height *float64 `form:"height"`
	Length *float64 `form:"length"`
	Weight *float64 `form:"weight"`
	Width  *float64 `form:"width"`
}

// ProductParams is the set of parameters that can be used when creating or updating a product.
type ProductParams struct {
	Params              `form:"*"`
	Active              *bool                    `form:"active"`
	Attributes          []*string                `form:"attributes"`
	Caption             *string                  `form:"caption"`
	DeactivateOn        []*string                `form:"deactivate_on"`
	Description         *string                  `form:"description"`
	ID                  *string                  `form:"id"`
	Images              []*string                `form:"images"`
	Name                *string                  `form:"name"`
	PackageDimensions   *PackageDimensionsParams `form:"package_dimensions"`
	Shippable           *bool                    `form:"shippable"`
	StatementDescriptor *string                  `form:"statement_descriptor"`
	Type                *string                  `form:"type"`
	UnitLabel           *string                  `form:"unit_label"`
	URL                 *string                  `form:"url"`
}

// PackageDimensions represents the dimension of a product or a SKU from the
// perspective of shipping.
type PackageDimensions struct {
	Height float64 `json:"height"`
	Length float64 `json:"length"`
	Weight float64 `json:"weight"`
	Width  float64 `json:"width"`
}

// Product is the resource representing a Stripe product.
// For more details see https://stripe.com/docs/api#products.
type Product struct {
	Active              bool               `json:"active"`
	Attributes          []string           `json:"attributes"`
	Caption             string             `json:"caption"`
	Created             int64              `json:"created"`
	DeactivateOn        []string           `json:"deactivate_on"`
	Description         string             `json:"description"`
	ID                  string             `json:"id"`
	Images              []string           `json:"images"`
	Livemode            bool               `json:"livemode"`
	Metadata            map[string]string  `json:"metadata"`
	Name                string             `json:"name"`
	PackageDimensions   *PackageDimensions `json:"package_dimensions"`
	Shippable           bool               `json:"shippable"`
	StatementDescriptor string             `json:"statement_descriptor"`
	Type                ProductType        `json:"type"`
	UnitLabel           string             `json:"unit_label"`
	URL                 string             `json:"url"`
	Updated             int64              `json:"updated"`
}

// ProductList is a list of products as retrieved from a list endpoint.
type ProductList struct {
	ListMeta
	Data []*Product `json:"data"`
}

// ProductListParams is the set of parameters that can be used when listing products.
type ProductListParams struct {
	ListParams `form:"*"`
	Active     *bool     `form:"active"`
	IDs        []*string `form:"ids"`
	Shippable  *bool     `form:"shippable"`
	URL        *string   `form:"url"`
	Type       *string   `form:"type"`
}

// UnmarshalJSON handles deserialization of a Product.
// This custom unmarshaling is needed because the resulting
// property may be an id or the full struct if it was expanded.
func (p *Product) UnmarshalJSON(data []byte) error {
	if id, ok := ParseID(data); ok {
		p.ID = id
		return nil
	}

	type product Product
	var v product
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	*p = Product(v)
	return nil
}
