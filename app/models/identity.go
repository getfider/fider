package models

import (
	"time"

	"github.com/getfider/fider/app/models/entities"
	"github.com/getfider/fider/app/models/enum"
	"github.com/getfider/fider/app/pkg/rand"
)

//Upload represents a file that has been uploaded to Fider
type Upload struct {
	ContentType string `db:"content_type"`
	Size        int    `db:"size"`
	Content     []byte `db:"file"`
}

//ImageUpload is the input model used to upload/remove an image
type ImageUpload struct {
	BlobKey string           `json:"bkey"`
	Upload  *ImageUploadData `json:"upload"`
	Remove  bool             `json:"remove"`
}

//ImageUploadData is the input model used to upload a new logo
type ImageUploadData struct {
	FileName    string `json:"fileName"`
	ContentType string `json:"contentType"`
	Content     []byte `json:"content"`
}

//UserInvitation is the model used to register an invite sent to an user
type UserInvitation struct {
	Email           string
	VerificationKey string
}

//GetEmail returns the invited user's email
func (e *UserInvitation) GetEmail() string {
	return e.Email
}

//GetName returns empty for this kind of process
func (e *UserInvitation) GetName() string {
	return ""
}

//GetUser returns the current user performing this action
func (e *UserInvitation) GetUser() *entities.User {
	return nil
}

//GetKind returns EmailVerificationKindUserInvitation
func (e *UserInvitation) GetKind() enum.EmailVerificationKind {
	return enum.EmailVerificationKindUserInvitation
}

//NewEmailVerification is used to register a new email verification process
type NewEmailVerification interface {
	GetEmail() string
	GetName() string
	GetUser() *entities.User
	GetKind() enum.EmailVerificationKind
}

//EmailVerification is the model used by email verification process
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

// GenerateSecretKey returns a 64 chars key
func GenerateSecretKey() string {
	return rand.String(64)
}
