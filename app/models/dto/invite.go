package dto

import (
	"github.com/getfider/fider/app/models/entities"
	"github.com/getfider/fider/app/models/enum"
)

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
