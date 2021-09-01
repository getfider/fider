package enum

//Role is the role of a user inside a tenant
type Role int

const (
	//RoleVisitor is the basic role for every user
	RoleVisitor Role = 1
	//RoleCollaborator has limited access to administrative console
	RoleCollaborator Role = 2
	//RoleAdministrator has full access to administrative console
	RoleAdministrator Role = 3
)

var roleIDs = map[Role]string{
	RoleVisitor:       "visitor",
	RoleCollaborator:  "collaborator",
	RoleAdministrator: "administrator",
}

var roleNames = map[string]Role{
	"visitor":       RoleVisitor,
	"collaborator":  RoleCollaborator,
	"administrator": RoleAdministrator,
}

// String returns the string version of the user role
func (role Role) String() string {
	return roleIDs[role]
}

// MarshalText returns the Text version of the user role
func (role Role) MarshalText() ([]byte, error) {
	return []byte(roleIDs[role]), nil
}

// UnmarshalText parse string into a user role
func (role *Role) UnmarshalText(text []byte) error {
	*role = roleNames[string(text)]
	return nil
}
