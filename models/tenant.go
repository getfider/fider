package models

//Tenant represents a tenant
type Tenant struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Domain string `json:"domain"`
}
