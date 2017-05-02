package models

import jwt "github.com/dgrijalva/jwt-go"
import "database/sql"

//Tenant represents a tenant
type Tenant struct {
	ID        int            `json:"id" db:"id"`
	Name      string         `json:"name" db:"name"`
	Subdomain string         `json:"subdomain" db:"subdomain"`
	CNAME     sql.NullString `json:"cname" db:"cname"`
}

//User represents an user inside our application
type User struct {
	ID             int             `json:"id" db:"id"`
	Name           string          `json:"name" db:"name"`
	Email          string          `json:"email" db:"email"`
	Tenant         *Tenant         `json:"tenant" db:"tenant"`
	Role           Role            `json:"role" db:"role"`
	Providers      []*UserProvider `json:"providers"`
	SupportedIdeas []int           `json:"supportedIdeas"`
}

//Role represents what the user can do inside its tenant
type Role int

var (
	//RoleVisitor is the basic role for every user
	RoleVisitor = Role(1)
	//RoleMember has limited access to administrative console
	RoleMember = Role(2)
	//RoleAdministrator has full access to administrative console
	RoleAdministrator = Role(3)
)

//HasProvider returns true if current user has given provider registered
func (u *User) HasProvider(provider string) bool {
	for _, p := range u.Providers {
		if p.Name == provider {
			return true
		}
	}
	return false
}

//UserProvider represents the relashionship between an User and an Authentication provide
type UserProvider struct {
	Name string `json:"name" db:"provider"`
	UID  string `json:"uid" db:"provider_uid"`
}

//WechyClaims represents what goes into JWT tokens
type WechyClaims struct {
	UserID    int    `json:"user/id"`
	UserName  string `json:"user/name"`
	UserEmail string `json:"user/email"`
	jwt.StandardClaims
}
