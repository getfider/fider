package models

import (
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/getfider/fider/app/pkg/uuid"
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
	IsPrivate      bool   `json:"isPrivate"`
}

var (
	//TenantActive is the default status for most tenants
	TenantActive = 1
	//TenantInactive is used for signup via email that requires user confirmation
	TenantInactive = 2
)

//User represents an user inside our application
type User struct {
	ID        int             `json:"id"`
	Name      string          `json:"name"`
	Email     string          `json:"-"`
	Tenant    *Tenant         `json:"-"`
	Role      Role            `json:"role"`
	Providers []*UserProvider `json:"-"`
}

//Role is the role of a user inside a tenant
type Role int

const (
	//RoleVisitor is the basic role for every user
	RoleVisitor Role = 1
	//RoleCollaborator has limited access to administrative console
	RoleCollaborator Role = 2
	//RoleAdministrator has full access to administrative console
	RoleAdministrator Role = 3
)

//EmailVerificationKind specifies which kind of process is being verified by email
type EmailVerificationKind int16

const (
	//EmailVerificationKindSignIn is the sign in by email process
	EmailVerificationKindSignIn EmailVerificationKind = 1
	//EmailVerificationKindSignUp is the sign up (create tenant) by name and email process
	EmailVerificationKindSignUp EmailVerificationKind = 2
	//EmailVerificationKindChangeEmail is the change user email process
	EmailVerificationKindChangeEmail EmailVerificationKind = 3
	//EmailVerificationKindUserInvitation is the sign in invitation sent to an user
	EmailVerificationKindUserInvitation EmailVerificationKind = 4
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

// IsCollaborator returns true if user has special permissions
func (u *User) IsCollaborator() bool {
	return u.Role == RoleCollaborator || u.Role == RoleAdministrator
}

// IsAdministrator returns true if user is administrator
func (u *User) IsAdministrator() bool {
	return u.Role == RoleAdministrator
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
	Email           string `json:"email" format:"lower"`
	VerificationKey string
	TenantName      string `json:"tenantName"`
	Subdomain       string `json:"subdomain" format:"lower"`
	UserClaims      *OAuthClaims
}

//GetEmail returns the email being verified
func (e *CreateTenant) GetEmail() string {
	return e.Email
}

//GetName returns the name of the email owner
func (e *CreateTenant) GetName() string {
	return e.Name
}

//GetUser returns the current user performing this action
func (e *CreateTenant) GetUser() *User {
	return nil
}

//GetKind returns EmailVerificationKindSignUp
func (e *CreateTenant) GetKind() EmailVerificationKind {
	return EmailVerificationKindSignUp
}

//UpdateTenantSettings is the input model used to update tenant general settings
type UpdateTenantSettings struct {
	Title          string `json:"title"`
	Invitation     string `json:"invitation"`
	WelcomeMessage string `json:"welcomeMessage"`
	CNAME          string `json:"cname" format:"lower"`
}

//UpdateTenantPrivacy is the input model used to update tenant privacy settings
type UpdateTenantPrivacy struct {
	IsPrivate bool `json:"isPrivate"`
}

//SignInByEmail is the input model when user request to sign in by email
type SignInByEmail struct {
	Email           string `json:"email" format:"lower"`
	VerificationKey string
}

//GetEmail returns the email being verified
func (e *SignInByEmail) GetEmail() string {
	return e.Email
}

//GetName returns empty for this kind of process
func (e *SignInByEmail) GetName() string {
	return ""
}

//GetUser returns the current user performing this action
func (e *SignInByEmail) GetUser() *User {
	return nil
}

//GetKind returns EmailVerificationKindSignIn
func (e *SignInByEmail) GetKind() EmailVerificationKind {
	return EmailVerificationKindSignIn
}

//ChangeUserEmail is the input model used to change current user's email
type ChangeUserEmail struct {
	Email           string `json:"email" format:"lower"`
	VerificationKey string
	Requestor       *User
}

//GetEmail returns the email being verified
func (e *ChangeUserEmail) GetEmail() string {
	return e.Email
}

//GetName returns empty for this kind of process
func (e *ChangeUserEmail) GetName() string {
	return ""
}

//GetUser returns the current user performing this action
func (e *ChangeUserEmail) GetUser() *User {
	return e.Requestor
}

//GetKind returns EmailVerificationKindSignIn
func (e *ChangeUserEmail) GetKind() EmailVerificationKind {
	return EmailVerificationKindChangeEmail
}

//UserInvitation is the model used to register an invite sent to an user
type UserInvitation struct {
	Email           string
	VerificationKey string
}

//GetEmail returns the invited user's email
func (e *UserInvitation) GetEmail() string {
	return e.Email
}

//GetName returns empty for this kind of process
func (e *UserInvitation) GetName() string {
	return ""
}

//GetUser returns the current user performing this action
func (e *UserInvitation) GetUser() *User {
	return nil
}

//GetKind returns EmailVerificationKindUserInvitation
func (e *UserInvitation) GetKind() EmailVerificationKind {
	return EmailVerificationKindUserInvitation
}

//NewEmailVerification is used to register a new email verification process
type NewEmailVerification interface {
	GetEmail() string
	GetName() string
	GetUser() *User
	GetKind() EmailVerificationKind
}

//EmailVerification is the model used by email verification process
type EmailVerification struct {
	Email      string
	Name       string
	Key        string
	UserID     int
	Kind       EmailVerificationKind
	CreatedOn  time.Time
	ExpiresOn  time.Time
	VerifiedOn *time.Time
}

// CompleteProfile is the model used to complete user profile during email sign in
type CompleteProfile struct {
	Key   string `json:"key"`
	Name  string `json:"name"`
	Email string
}

// UpdateUserSettings is the model used to update user's settings
type UpdateUserSettings struct {
	Name     string            `json:"name"`
	Settings map[string]string `json:"settings"`
}

// ChangeUserRole is the input model change role of an user
type ChangeUserRole struct {
	UserID int  `route:"user_id"`
	Role   Role `json:"role"`
}

// InviteUsers is used to invite new users into Fider
type InviteUsers struct {
	Subject    string   `json:"subject"`
	Message    string   `json:"message"`
	Recipients []string `json:"recipients" format:"lower"`
}

// GenerateVerificationKey used on email verifications
func GenerateVerificationKey() string {
	return strings.Replace(uuid.NewV4().String(), "-", "", 4)
}
