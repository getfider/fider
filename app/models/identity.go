package models

import (
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

//Tenant represents a tenant
type Tenant struct {
	ID             int    `json:"id"`
	Name           string `json:"name"`
	Subdomain      string `json:"subdomain"`
	Invitation     string `json:"invitation"`
	WelcomeMessage string `json:"welcomeMessage"`
	CNAME          string `json:"cname"`
	Status         int    `json:"-"`
}

var (
	//TenantActive is the default status for most tenants
	TenantActive = 1
	//TenantInactive is used for signup via e-mail that requires user confirmation
	TenantInactive = 2
)

//User represents an user inside our application
type User struct {
	ID        int             `json:"id"`
	Name      string          `json:"name"`
	Email     string          `json:"-"`
	Gravatar  string          `json:"gravatar"`
	Tenant    *Tenant         `json:"-"`
	Role      int             `json:"role"`
	Providers []*UserProvider `json:"-"`
}

var (
	//RoleVisitor is the basic role for every user
	RoleVisitor = 1
	//RoleMember has limited access to administrative console
	RoleMember = 2
	//RoleAdministrator has full access to administrative console
	RoleAdministrator = 3
)

//HasProvider returns true if current user has registered with given provider
func (u *User) HasProvider(provider string) bool {
	for _, p := range u.Providers {
		if p.Name == provider {
			return true
		}
	}
	return false
}

//IsStaff returns true if user has special permissions
func (u *User) IsStaff() bool {
	return u.Role >= RoleMember
}

//UserProvider represents the relashionship between an User and an Authentication provide
type UserProvider struct {
	Name string
	UID  string
}

//FiderClaims represents what goes into JWT tokens
type FiderClaims struct {
	UserID    int    `json:"user/id"`
	UserName  string `json:"user/name"`
	UserEmail string `json:"user/email"`
	jwt.StandardClaims
}

//OAuthClaims represents what goes into temporary OAuth JWT tokens
type OAuthClaims struct {
	OAuthID       string `json:"oauth/id"`
	OAuthProvider string `json:"oauth/provider"`
	OAuthName     string `json:"oauth/name"`
	OAuthEmail    string `json:"oauth/email"`
	jwt.StandardClaims
}

//CreateTenant is the input model used to create a tenant
type CreateTenant struct {
	Token           string `json:"token"`
	Name            string `json:"name"`
	Email           string `json:"email"`
	VerificationKey string
	TenantName      string `json:"tenantName"`
	Subdomain       string `json:"subdomain"`
	UserClaims      *OAuthClaims
}

//UpdateTenantSettings is the input model used to update tenant settings
type UpdateTenantSettings struct {
	Title          string `json:"title"`
	Invitation     string `json:"invitation"`
	WelcomeMessage string `json:"welcomeMessage"`
	UserClaims     *OAuthClaims
}

//SignInByEmail is the input model when user request to sign in by email
type SignInByEmail struct {
	Email           string `json:"email"`
	VerificationKey string
}

//SignInRequest is the model used by e-mail verification process
type SignInRequest struct {
	Email      string
	Name       string
	Key        string
	CreatedOn  time.Time
	ExpiresOn  time.Time
	VerifiedOn *time.Time
}

// CompleteProfile is the model used to complete user profile during e-mail sign in
type CompleteProfile struct {
	Key   string `json:"key"`
	Name  string `json:"name"`
	Email string
}

// UpdateUserSettings is the model used to update user's settings
type UpdateUserSettings struct {
	Name string `json:"name"`
}
