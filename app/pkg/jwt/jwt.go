package jwt

import (
	"github.com/WeCanHearYou/wechy/app/models"
	"github.com/WeCanHearYou/wechy/app/pkg/env"
	jwtgo "github.com/dgrijalva/jwt-go"
)

var jwtSecret = env.MustGet("JWT_SECRET")

//Encode creates new JWT tokens with given claims
func Encode(claims *models.WechyClaims) (string, error) {
	jwtToken := jwtgo.NewWithClaims(jwtgo.GetSigningMethod("HS256"), claims)
	return jwtToken.SignedString([]byte(jwtSecret))
}

//Decode extract claims from JWT tokens
func Decode(token string) (*models.WechyClaims, error) {
	claims := &models.WechyClaims{}
	jwtToken, err := jwtgo.ParseWithClaims(token, claims, func(t *jwtgo.Token) (interface{}, error) {
		return []byte(jwtSecret), nil
	})

	if err == nil && jwtToken.Valid {
		return claims, nil
	}

	return nil, err
}
