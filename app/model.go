package app

import jwt "github.com/dgrijalva/jwt-go"

//Tenant represents a tenant
type Tenant struct {
	ID     int64  `json:"id"`
	Name   string `json:"name"`
	Domain string `json:"domain"`
}

//User represents an user inside our application
type User struct {
	ID        int64           `json:"id"`
	Name      string          `json:"name"`
	Email     string          `json:"email"`
	Providers []*UserProvider `json:"providers"`
}

//UserProvider represents the relashionship between an User and an Authentication provide
type UserProvider struct {
	Name string `json:"name"`
	UID  string `json:"uid"`
}

//WechyClaims represents what goes into JWT tokens
type WechyClaims struct {
	UserID    int64  `json:"user/id"`
	UserName  string `json:"user/name"`
	UserEmail string `json:"user/email"`
	jwt.StandardClaims
}
