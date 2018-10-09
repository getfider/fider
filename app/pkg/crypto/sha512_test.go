package crypto_test

import (
	"testing"

	. "github.com/getfider/fider/app/pkg/assert"
	"github.com/getfider/fider/app/pkg/crypto"
)

func TestSHA512Hash(t *testing.T) {
	RegisterT(t)

	hash := crypto.SHA512("Fider")

	Expect(hash).Equals("262d21f30715f2b226264844c4ab2a934a4c0241321f77bebbca191e172df93da71c939c56fcb4bbdd8895fa8c496882d38e3ce66d9d4e3dee5bacde01e73988")
}
