package jwt_test

import (
	"testing"

	"github.com/WeCanHearYou/wechy/app/models"
	"github.com/WeCanHearYou/wechy/app/pkg/jwt"
	. "github.com/onsi/gomega"
)

func TestJWT_Encode(t *testing.T) {
	RegisterTestingT(t)

	claims := &models.WechyClaims{
		UserID:    424,
		UserName:  "Jon Snow",
		UserEmail: "jon.snow@got.com",
	}

	token, err := jwt.Encode(claims)
	Expect(token).NotTo(BeNil())
	Expect(err).To(BeNil())
}

func TestJWT_Decode(t *testing.T) {
	RegisterTestingT(t)

	claims := &models.WechyClaims{
		UserID:    424,
		UserName:  "Jon Snow",
		UserEmail: "jon.snow@got.com",
	}

	token, _ := jwt.Encode(claims)

	decoded, err := jwt.Decode(token)
	Expect(err).To(BeNil())
	Expect(decoded.UserID).To(Equal(claims.UserID))
	Expect(decoded.UserName).To(Equal(claims.UserName))
	Expect(decoded.UserEmail).To(Equal(claims.UserEmail))
}

func TestJWT_DecodeChangedToken(t *testing.T) {
	RegisterTestingT(t)

	claims := &models.WechyClaims{
		UserID:    424,
		UserName:  "Jon Snow",
		UserEmail: "jon.snow@got.com",
	}

	token, _ := jwt.Encode(claims)

	decoded, err := jwt.Decode(token + "foo")
	Expect(err).ToNot(BeNil())
	Expect(decoded).To(BeNil())
}
