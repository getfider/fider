package jwt

import (
	jwtgo "github.com/dgrijalva/jwt-go"
)

// Metadata is the basic JWT information
type Metadata = jwtgo.StandardClaims

// FiderClaims represents what goes into JWT tokens
type FiderClaims struct {
	UserID    int    `json:"user/id"`
	UserName  string `json:"user/name"`
	UserEmail string `json:"user/email"`
	Origin    string `json:"origin"`
	Metadata
}

// OAuthClaims represents what goes into temporary OAuth JWT tokens
type OAuthClaims struct {
	OAuthID       string `json:"oauth/id"`
	OAuthProvider string `json:"oauth/provider"`
	OAuthName     string `json:"oauth/name"`
	OAuthEmail    string `json:"oauth/email"`
	Metadata
}
