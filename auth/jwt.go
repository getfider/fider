package auth

import (
	"github.com/WeCanHearYou/wchy/env"
	jwt "github.com/dgrijalva/jwt-go"
)

//WchyClaims represents what goes into JWT tokens
type WchyClaims struct {
	UserID    int64  `json:"user/id"`
	UserName  string `json:"user/name"`
	UserEmail string `json:"user/email"`
	jwt.StandardClaims
}

var jwtSecret = env.MustGet("JWT_SECRET")

//Encode creates new JWT tokens with given claims
func Encode(claims *WchyClaims) (string, error) {
	jwtToken := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), claims)
	return jwtToken.SignedString([]byte(jwtSecret))
}

//Decode extract claims from JWT tokens
func Decode(token string) (*WchyClaims, error) {
	claims := &WchyClaims{}
	jwtToken, err := jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(jwtSecret), nil
	})

	if err == nil && jwtToken.Valid {
		return claims, nil
	}

	return nil, err
}
