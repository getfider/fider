package jwt

import (
	"fmt"
	"time"

	"github.com/getfider/fider/app/pkg/env"
	"github.com/getfider/fider/app/pkg/errors"
	jwtgo "github.com/golang-jwt/jwt/v4"
)

var jwtSecret = env.Config.JWTSecret

// Metadata is the basic JWT information
type Metadata = jwtgo.RegisteredClaims

func Time(t time.Time) *jwtgo.NumericDate {
	return jwtgo.NewNumericDate(t)
}

const (
	//FiderClaimsOriginUI is assigned to Fider claims when the Auth Token is generated through the UI
	FiderClaimsOriginUI = "ui"
	//FiderClaimsOriginAPI is assigned to Fider claims when the Auth Token is generated through the API
	FiderClaimsOriginAPI = "api"
)

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

// Encode creates new JWT token with given claims
func Encode(claims jwtgo.Claims) (string, error) {
	jwtToken := jwtgo.NewWithClaims(jwtgo.GetSigningMethod("HS256"), claims)
	token, err := jwtToken.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", errors.Wrap(err, "failed to encode the requested claims")
	}
	return token, nil
}

// DecodeFiderClaims extract claims from JWT tokens
func DecodeFiderClaims(token string) (*FiderClaims, error) {
	claims := &FiderClaims{}
	err := decode(token, claims)
	if err != nil {
		return nil, errors.Wrap(err, "failed to decode Fider claims")
	}
	return claims, nil
}

// DecodeOAuthClaims extract OAuthClaims from given JWT token
func DecodeOAuthClaims(token string) (*OAuthClaims, error) {
	claims := &OAuthClaims{}
	err := decode(token, claims)
	if err != nil {
		return nil, errors.Wrap(err, "failed to decode OAuth claims")
	}
	return claims, nil
}

func decode(token string, claims jwtgo.Claims) error {
	jwtToken, err := jwtgo.ParseWithClaims(token, claims, func(t *jwtgo.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwtgo.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", t.Header["alg"])
		}

		return []byte(jwtSecret), nil
	})

	if err == nil {
		err = claims.Valid()
	}

	if err == nil && jwtToken.Valid {
		return nil
	}

	return err
}
