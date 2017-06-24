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

//DecodeFiderClaims extract claims from JWT tokens
func DecodeFiderClaims(token string) (*models.FiderClaims, error) {
	claims := &models.FiderClaims{}
	err := decode(token, claims)
	if err == nil {
		return claims, nil
	}
	return nil, err
}

//DecodeOAuthClaims extract OAuthClaims from given JWT token
func DecodeOAuthClaims(token string) (*models.OAuthClaims, error) {
	claims := &models.OAuthClaims{}
	err := decode(token, claims)
	if err == nil {
		return claims, nil
	}
	return nil, err
}

func decode(token string, claims jwtgo.Claims) error {
	jwtToken, err := jwtgo.ParseWithClaims(token, claims, func(t *jwtgo.Token) (interface{}, error) {
		return []byte(jwtSecret), nil
	})

	if err == nil && jwtToken.Valid {
		return nil
	}
	return err
}
