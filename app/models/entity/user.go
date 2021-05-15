package entity

import "github.com/getfider/fider/app/models/enum"

//User represents an user inside our application
type User struct {
	ID            int             `json:"id"`
	Name          string          `json:"name"`
	Tenant        *Tenant         `json:"-"`
	Email         string          `json:"-"`
	Role          enum.Role       `json:"role"`
	Providers     []*UserProvider `json:"-"`
	AvatarBlobKey string          `json:"-"`
	AvatarType    enum.AvatarType `json:"-"`
	AvatarURL     string          `json:"avatarURL,omitempty"`
	Status        enum.UserStatus `json:"status"`
}

//HasProvider returns true if current user has registered with given provider
func (u *User) HasProvider(provider string) bool {
	for _, p := range u.Providers {
		if p.Name == provider {
			return true
		}
	}
	return false
}

// IsCollaborator returns true if user has special permissions
func (u *User) IsCollaborator() bool {
	return u.Role == enum.RoleCollaborator || u.Role == enum.RoleAdministrator
}

// IsAdministrator returns true if user is administrator
func (u *User) IsAdministrator() bool {
	return u.Role == enum.RoleAdministrator
}

//UserProvider represents the relationship between an User and an Authentication provide
type UserProvider struct {
	Name string
	UID  string
}
