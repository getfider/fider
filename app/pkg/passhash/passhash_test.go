package passhash

import (
	"testing"
)

func TestStringString(t *testing.T) {
	plainText := "This is a test."

	hash, err := HashString(plainText)

	if err != nil {
		t.Error(err)
	}

	if !MatchString(hash, plainText) {
		t.Error("Password does not match")
	}
}

func TestByteByte(t *testing.T) {
	plainText := []byte("This is a test.")

	hash, err := HashBytes(plainText)

	if err != nil {
		t.Error(err)
	}

	if !MatchBytes(hash, plainText) {
		t.Error("Password does not match")
	}
}

func TestStringByte(t *testing.T) {
	plainText := "This is a test."

	hash, err := HashString(plainText)

	if err != nil {
		t.Error(err)
	}

	if !MatchBytes([]byte(hash), []byte(plainText)) {
		t.Error("Password does not match")
	}
}

func TestByteString(t *testing.T) {
	plainText := []byte("This is a test.")

	hash, err := HashBytes(plainText)

	if err != nil {
		t.Error(err)
	}

	if !MatchString(string(hash), string(plainText)) {
		t.Error("Password does not match")
	}
}
