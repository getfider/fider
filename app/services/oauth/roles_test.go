package oauth

import (
	"testing"
)

func TestExtractRolesFromJSON_ArrayOfObjects(t *testing.T) {
	// IServ-style response
	jsonBody := `{
		"email": "test@example.com",
		"roles": [
			{
				"uuid": "c9b85a22-1a6c-46e1-99eb-36cb88001fda",
				"id": "ROLE_TEACHER",
				"displayName": "Lehrer"
			},
			{
				"uuid": "f66aee04-c335-4299-9cfe-ca7176cc0213",
				"id": "ROLE_ADMIN",
				"displayName": "Administrator"
			}
		]
	}`

	roles := extractRolesFromJSON(jsonBody, "roles[].id")

	if len(roles) != 2 {
		t.Errorf("Expected 2 roles, got %d", len(roles))
	}

	if roles[0] != "ROLE_TEACHER" {
		t.Errorf("Expected first role to be ROLE_TEACHER, got %s", roles[0])
	}

	if roles[1] != "ROLE_ADMIN" {
		t.Errorf("Expected second role to be ROLE_ADMIN, got %s", roles[1])
	}
}

func TestExtractRolesFromJSON_ArrayOfStrings(t *testing.T) {
	jsonBody := `{
		"email": "test@example.com",
		"roles": ["ROLE_ADMIN", "ROLE_USER"]
	}`

	roles := extractRolesFromJSON(jsonBody, "roles")

	if len(roles) != 2 {
		t.Errorf("Expected 2 roles, got %d", len(roles))
	}

	if roles[0] != "ROLE_ADMIN" {
		t.Errorf("Expected first role to be ROLE_ADMIN, got %s", roles[0])
	}
}

func TestExtractRolesFromJSON_NestedPath(t *testing.T) {
	jsonBody := `{
		"user": {
			"profile": {
				"roles": ["ROLE_ADMIN"]
			}
		}
	}`

	roles := extractRolesFromJSON(jsonBody, "user.profile.roles")

	if len(roles) != 1 {
		t.Errorf("Expected 1 role, got %d", len(roles))
	}

	if roles[0] != "ROLE_ADMIN" {
		t.Errorf("Expected role to be ROLE_ADMIN, got %s", roles[0])
	}
}

