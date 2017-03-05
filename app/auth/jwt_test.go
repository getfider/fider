package auth_test

import (
	"testing"

	"github.com/WeCanHearYou/wchy/app/auth"
	. "github.com/onsi/gomega"
)

func TestJWT_Encode(t *testing.T) {
	RegisterTestingT(t)

	claims := &auth.WchyClaims{
		UserID:    424,
		UserName:  "Jon Snow",
		UserEmail: "jon.snow@got.com",
	}

	token, err := auth.Encode(claims)
	Expect(token).NotTo(BeNil())
	Expect(err).To(BeNil())
}

func TestJWT_Decode(t *testing.T) {
	RegisterTestingT(t)

	claims := &auth.WchyClaims{
		UserID:    424,
		UserName:  "Jon Snow",
		UserEmail: "jon.snow@got.com",
	}

	token, _ := auth.Encode(claims)

	decoded, err := auth.Decode(token)
	Expect(err).To(BeNil())
	Expect(decoded.UserID).To(Equal(claims.UserID))
	Expect(decoded.UserName).To(Equal(claims.UserName))
	Expect(decoded.UserEmail).To(Equal(claims.UserEmail))
}
