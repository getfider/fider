package stripe

import "encoding/json"

// SKUInventoryType describe's the possible value for inventory type
type SKUInventoryType string

// List of values that SKUInventoryType can take.
const (
	SKUInventoryTypeBucket   SKUInventoryType = "bucket"
	SKUInventoryTypeFinite   SKUInventoryType = "finite"
	SKUInventoryTypeInfinite SKUInventoryType = "infinite"
)

// SKUInventoryValue describe's the possible value for inventory value
type SKUInventoryValue string

// List of values that SKUInventoryValue can take.
const (
	SKUInventoryValueInStock    SKUInventoryValue = "in_stock"
	SKUInventoryValueLimited    SKUInventoryValue = "limited"
	SKUInventoryValueOutOfStock SKUInventoryValue = "out_of_stock"
)

// InventoryParams is the set of parameters allowed as inventory on a SKU.
type InventoryParams struct {
	Quantity *int64  `form:"quantity"`
	Type     *string `form:"type"`
	Value    *string `form:"value"`
}

// SKUParams is the set of parameters allowed on SKU creation or update.
type SKUParams struct {
	Params            `form:"*"`
	Active            *bool                    `form:"active"`
	Attributes        map[string]string        `form:"attributes"`
	Currency          *string                  `form:"currency"`
	Description       *string                  `form:"description"`
	ID                *string                  `form:"id"`
	Image             *string                  `form:"image"`
	Inventory         *InventoryParams         `form:"inventory"`
	PackageDimensions *PackageDimensionsParams `form:"package_dimensions"`
	Price             *int64                   `form:"price"`
	Product           *string                  `form:"product"`
}

// Inventory represents the inventory options of a SKU.
type Inventory struct {
	Quantity int64             `json:"quantity"`
	Type     SKUInventoryType  `json:"type"`
	Value    SKUInventoryValue `json:"value"`
}

// SKU is the resource representing a SKU.
// For more details see https://stripe.com/docs/api#skus.
type SKU struct {
	Active            bool               `json:"active"`
	Attributes        map[string]string  `json:"attributes"`
	Created           int64              `json:"created"`
	Currency          Currency           `json:"currency"`
	Description       string             `json:"description"`
	ID                string             `json:"id"`
	Image             string             `json:"image"`
	Inventory         *Inventory         `json:"inventory"`
	Livemode          bool               `json:"livemode"`
	Metadata          map[string]string  `json:"metadata"`
	PackageDimensions *PackageDimensions `json:"package_dimensions"`
	Price             int64              `json:"price"`
	Product           *Product           `json:"product"`
	Updated           int64              `json:"updated"`
}

// SKUList is a list of SKUs as returned from a list endpoint.
type SKUList struct {
	ListMeta
	Data []*SKU `json:"data"`
}

// SKUListParams is the set of parameters that can be used when listing SKUs.
type SKUListParams struct {
	ListParams `form:"*"`
	Active     *bool             `form:"active"`
	Attributes map[string]string `form:"attributes"`
	IDs        []*string         `form:"ids"`
	InStock    *bool             `form:"in_stock"`
	Product    *string           `form:"product"`
}

// UnmarshalJSON handles deserialization of a SKU.
// This custom unmarshaling is needed because the resulting
// property may be an id or the full struct if it was expanded.
func (s *SKU) UnmarshalJSON(data []byte) error {
	if id, ok := ParseID(data); ok {
		s.ID = id
		return nil
	}

	type sku SKU
	var v sku
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	*s = SKU(v)
	return nil
}
