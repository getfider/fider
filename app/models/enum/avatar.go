package enum

//AvatarType are the possible types of a user avatar
type AvatarType int

var (
	//AvatarTypeLetter is the default avatar type for users
	AvatarTypeLetter AvatarType = 1
	//AvatarTypeGravatar fetches avatar from gravatar (if available)
	AvatarTypeGravatar AvatarType = 2
	//AvatarTypeCustom uses a user uploaded avatar
	AvatarTypeCustom AvatarType = 3
)

var avatarTypesIDs = map[AvatarType]string{
	AvatarTypeLetter:   "letter",
	AvatarTypeGravatar: "gravatar",
	AvatarTypeCustom:   "custom",
}

var avatarTypesName = map[string]AvatarType{
	"letter":   AvatarTypeLetter,
	"gravatar": AvatarTypeGravatar,
	"custom":   AvatarTypeCustom,
}

// String returns the string version of the avatar type
func (t AvatarType) String() string {
	return avatarTypesIDs[t]
}

// MarshalText returns the Text version of the avatar type
func (t AvatarType) MarshalText() ([]byte, error) {
	return []byte(avatarTypesIDs[t]), nil
}

// UnmarshalText parse string into a avatar type
func (t *AvatarType) UnmarshalText(text []byte) error {
	*t = avatarTypesName[string(text)]
	return nil
}
