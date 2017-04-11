package infra_test

import (
	"testing"

	"github.com/WeCanHearYou/wechy/app"
	"github.com/WeCanHearYou/wechy/app/infra"
	. "github.com/onsi/gomega"
)

func TestJWT_Encode(t *testing.T) {
	RegisterTestingT(t)

	claims := &app.WechyClaims{
		UserID:    424,
		UserName:  "Jon Snow",
		UserEmail: "jon.snow@got.com",
	}

	token, err := infra.Encode(claims)
	Expect(token).NotTo(BeNil())
	Expect(err).To(BeNil())
}

func TestJWT_Decode(t *testing.T) {
	RegisterTestingT(t)

	claims := &app.WechyClaims{
		UserID:    424,
		UserName:  "Jon Snow",
		UserEmail: "jon.snow@got.com",
	}

	token, _ := infra.Encode(claims)

	decoded, err := infra.Decode(token)
	Expect(err).To(BeNil())
	Expect(decoded.UserID).To(Equal(claims.UserID))
	Expect(decoded.UserName).To(Equal(claims.UserName))
	Expect(decoded.UserEmail).To(Equal(claims.UserEmail))
}
