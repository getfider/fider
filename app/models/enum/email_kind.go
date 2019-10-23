package enum

//EmailVerificationKind specifies which kind of process is being verified by email
type EmailVerificationKind int16

const (
	//EmailVerificationKindSignIn is the sign in by email process
	EmailVerificationKindSignIn EmailVerificationKind = 1
	//EmailVerificationKindSignUp is the sign up (create tenant) by name and email process
	EmailVerificationKindSignUp EmailVerificationKind = 2
	//EmailVerificationKindChangeEmail is the change user email process
	EmailVerificationKindChangeEmail EmailVerificationKind = 3
	//EmailVerificationKindUserInvitation is the sign in invitation sent to an user
	EmailVerificationKindUserInvitation EmailVerificationKind = 4
)
