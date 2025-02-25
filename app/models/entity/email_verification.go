package entity

import (
	"time"

	"github.com/Spicy-Bush/fider-tarkov-community/app/models/enum"
	"github.com/Spicy-Bush/fider-tarkov-community/app/pkg/rand"
)

// EmailVerification is the model used by email verification process
type EmailVerification struct {
	Email      string
	Name       string
	Key        string
	UserID     int
	Kind       enum.EmailVerificationKind
	CreatedAt  time.Time
	ExpiresAt  time.Time
	VerifiedAt *time.Time
}

// GenerateEmailVerificationKey returns a 64 chars key
func GenerateEmailVerificationKey() string {
	return rand.String(64)
}
