package model

//Tenant represents a tenant
type Tenant struct {
	ID     int64  `json:"id"`
	Name   string `json:"name"`
	Domain string `json:"domain"`
}
