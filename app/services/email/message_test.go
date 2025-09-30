package email_test

import (
	"testing"

	"github.com/getfider/fider/app/services/email"
	. "github.com/getfider/fider/app/pkg/assert"
)

func TestEncodeSubject_ASCIIOnly(t *testing.T) {
	RegisterT(t)
	
	subject := "Hello World"
	encoded := email.EncodeSubject(subject)
	
	Expect(encoded).Equals("Hello World")
}

func TestEncodeSubject_WithNonASCII(t *testing.T) {
	RegisterT(t)
	
	subject := "Test ä ö ü"
	encoded := email.EncodeSubject(subject)
	
	// The encoded result should be in RFC 2047 format
	Expect(encoded).ContainsSubstring("=?utf-8?q?")
	Expect(encoded).ContainsSubstring("?=")
}

func TestEncodeSubject_WithEmoji(t *testing.T) {
	RegisterT(t)
	
	subject := "Hello 👋 World 🌍"
	encoded := email.EncodeSubject(subject)
	
	// The encoded result should be in RFC 2047 format
	Expect(encoded).ContainsSubstring("=?utf-8?q?")
	Expect(encoded).ContainsSubstring("?=")
}

func TestEncodeSubject_Empty(t *testing.T) {
	RegisterT(t)
	
	subject := ""
	encoded := email.EncodeSubject(subject)
	
	Expect(encoded).Equals("")
}

func TestEncodeSubject_WithUmlautsAndSpecialChars(t *testing.T) {
	RegisterT(t)
	
	subject := "Öffentliche Ankündigung - Größere Änderungen"
	encoded := email.EncodeSubject(subject)
	
	// The encoded result should be in RFC 2047 format
	Expect(encoded).ContainsSubstring("=?utf-8?q?")
	Expect(encoded).ContainsSubstring("?=")
}
