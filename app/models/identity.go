package models

import jwt "github.com/dgrijalva/jwt-go"

//Tenant represents a tenant
type Tenant struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Subdomain string `json:"subdomain"`
	CNAME     string `json:"cname"`
}

//User represents an user inside our application
type User struct {
	ID             int             `json:"id"`
	Name           string          `json:"name"`
	Email          string          `json:"email"`
	Tenant         *Tenant         `json:"tenant"`
	Role           int             `json:"role"`
	Providers      []*UserProvider `json:"providers"`
	SupportedIdeas []int           `json:"supportedIdeas"`
}

var (
	//RoleVisitor is the basic role for every user
	RoleVisitor = 1
	//RoleMember has limited access to administrative console
	RoleMember = 2
	//RoleAdministrator has full access to administrative console
	RoleAdministrator = 3
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
	Name string `json:"name"`
	UID  string `json:"uid"`
}

//FiderClaims represents what goes into JWT tokens
type FiderClaims struct {
	UserID    int    `json:"user/id"`
	UserName  string `json:"user/name"`
	UserEmail string `json:"user/email"`
	jwt.StandardClaims
}
