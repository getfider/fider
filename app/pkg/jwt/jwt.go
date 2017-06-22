package jwt

import (
	jwtgo "github.com/dgrijalva/jwt-go"
	"github.com/getfider/fider/app/models"
	"github.com/getfider/fider/app/pkg/env"
)

var jwtSecret = env.MustGet("JWT_SECRET")

//Encode creates new JWT tokens with given claims
func Encode(claims jwtgo.Claims) (string, error) {
	jwtToken := jwtgo.NewWithClaims(jwtgo.GetSigningMethod("HS256"), claims)
	return jwtToken.SignedString([]byte(jwtSecret))
}

//Decode extract claims from JWT tokens
func Decode(token string) (*models.FiderClaims, error) {
	claims := &models.FiderClaims{}
	jwtToken, err := jwtgo.ParseWithClaims(token, claims, func(t *jwtgo.Token) (interface{}, error) {
		return []byte(jwtSecret), nil
	})

	if err == nil && jwtToken.Valid {
		return claims, nil
	}

	return nil, err
}

//DecodeAsOAuthClaims extract OAuthClaims from given JWT token
func DecodeAsOAuthClaims(token string) (*models.OAuthClaims, error) {
	claims := &models.OAuthClaims{}
	jwtToken, err := jwtgo.ParseWithClaims(token, claims, func(t *jwtgo.Token) (interface{}, error) {
		return []byte(jwtSecret), nil
	})

	if err == nil && jwtToken.Valid {
		return claims, nil
	}

	return nil, err
}
