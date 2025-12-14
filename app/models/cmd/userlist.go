package cmd

import "github.com/getfider/fider/app/models/enum"

type UserListCreateCompany struct {
	Name       string
	TenantId   int
	SignedUpAt string
	Plan       enum.Plan
	Subdomain  string
	UserId     int
	UserEmail  string
	UserName   string
}

type UserListUpdateCompany struct {
	TenantId int
	Name     string
	Plan     enum.Plan
}

type UserListUpdateUser struct {
	Id       int
	TenantId int
	Email    string
	Name     string
}

type UserListHandleRoleChange struct {
	Id   int
	Role enum.Role
}
