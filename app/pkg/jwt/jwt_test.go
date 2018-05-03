package jwt_test

import (
	"testing"

	"github.com/getfider/fider/app/models"
	. "github.com/getfider/fider/app/pkg/assert"
	"github.com/getfider/fider/app/pkg/jwt"
)

func TestJWT_Encode(t *testing.T) {
	RegisterT(t)

	claims := &models.FiderClaims{
		UserID:    424,
		UserName:  "Jon Snow",
		UserEmail: "jon.snow@got.com",
	}

	token, err := jwt.Encode(claims)
	Expect(token).IsNotEmpty()
	Expect(err).IsNil()
}

func TestJWT_Decode(t *testing.T) {
	RegisterT(t)

	claims := &models.FiderClaims{
		UserID:    424,
		UserName:  "Jon Snow",
		UserEmail: "jon.snow@got.com",
	}

	token, _ := jwt.Encode(claims)

	decoded, err := jwt.DecodeFiderClaims(token)
	Expect(err).IsNil()
	Expect(decoded.UserID).Equals(claims.UserID)
	Expect(decoded.UserName).Equals(claims.UserName)
	Expect(decoded.UserEmail).Equals(claims.UserEmail)
}

func TestJWT_DecodeOAuthClaims(t *testing.T) {
	RegisterT(t)

	claims := &models.OAuthClaims{
		OAuthID:       "2",
		OAuthEmail:    "jon.snow@got.com",
		OAuthName:     "Jon Snow",
		OAuthProvider: "facebook",
	}

	token, _ := jwt.Encode(claims)

	decoded, err := jwt.DecodeOAuthClaims(token)
	Expect(err).IsNil()
	Expect(decoded.OAuthID).Equals(claims.OAuthID)
	Expect(decoded.OAuthEmail).Equals(claims.OAuthEmail)
	Expect(decoded.OAuthName).Equals(claims.OAuthName)
	Expect(decoded.OAuthProvider).Equals(claims.OAuthProvider)
}

func TestJWT_DecodeChangedToken(t *testing.T) {
	RegisterT(t)

	claims := &models.FiderClaims{
		UserID:    424,
		UserName:  "Jon Snow",
		UserEmail: "jon.snow@got.com",
	}

	token, _ := jwt.Encode(claims)

	decoded, err := jwt.DecodeFiderClaims(token + "foo")
	Expect(err).IsNotNil()
	Expect(decoded).IsNil()
}
