package handlers

// Internal test file for hasAllowedRole — uses package handlers (not handlers_test)
// so the unexported function is accessible without exporting it.

import (
	"testing"

	"github.com/getfider/fider/app/pkg/env"
)

// setAllowedRoles sets env.Config.OAuth.AllowedRoles for the duration of a test
// and automatically restores the original value when the test finishes.
func setAllowedRoles(t *testing.T, value string) {
	t.Helper()
	original := env.Config.OAuth.AllowedRoles
	env.Config.OAuth.AllowedRoles = value
	t.Cleanup(func() { env.Config.OAuth.AllowedRoles = original })
}

// --- empty / unconfigured OAUTH_ALLOWED_ROLES ---

func TestHasAllowedRole_EmptyConfig_AllowsAll(t *testing.T) {
	setAllowedRoles(t, "")
	if !hasAllowedRole([]string{"ROLE_ADMIN"}, "roles") {
		t.Error("expected true when OAUTH_ALLOWED_ROLES is empty")
	}
}

func TestHasAllowedRole_EmptyConfig_AllowsEmptyRoles(t *testing.T) {
	setAllowedRoles(t, "")
	if !hasAllowedRole([]string{}, "roles") {
		t.Error("expected true when OAUTH_ALLOWED_ROLES is empty and user has no roles")
	}
}

func TestHasAllowedRole_WhitespaceOnlyConfig_AllowsAll(t *testing.T) {
	setAllowedRoles(t, "   ")
	if !hasAllowedRole([]string{"ROLE_ADMIN"}, "roles") {
		t.Error("expected true when OAUTH_ALLOWED_ROLES is whitespace only")
	}
}

// --- provider has no JSONUserRolesPath (skip check) ---

func TestHasAllowedRole_NoRolesPath_SkipsCheck(t *testing.T) {
	setAllowedRoles(t, "ROLE_ADMIN")
	// Provider without a roles path must always be allowed through,
	// regardless of whether the user carries matching roles or not.
	if !hasAllowedRole([]string{}, "") {
		t.Error("expected true when provider has no JSONUserRolesPath (check should be skipped)")
	}
}

func TestHasAllowedRole_WhitespaceRolesPath_SkipsCheck(t *testing.T) {
	setAllowedRoles(t, "ROLE_ADMIN")
	if !hasAllowedRole([]string{}, "   ") {
		t.Error("expected true when JSONUserRolesPath is whitespace only (check should be skipped)")
	}
}

// --- matching role ---

func TestHasAllowedRole_MatchingSingleRole(t *testing.T) {
	setAllowedRoles(t, "ROLE_ADMIN")
	if !hasAllowedRole([]string{"ROLE_ADMIN"}, "roles") {
		t.Error("expected true when user has the required role")
	}
}

func TestHasAllowedRole_MatchingOneOfMultipleAllowed(t *testing.T) {
	setAllowedRoles(t, "ROLE_ADMIN,ROLE_TEACHER")
	if !hasAllowedRole([]string{"ROLE_TEACHER"}, "roles") {
		t.Error("expected true when user has one of the allowed roles")
	}
}

func TestHasAllowedRole_MatchingMultipleAllowed(t *testing.T) {
	setAllowedRoles(t, "ROLE_ADMIN,ROLE_TEACHER")
	if !hasAllowedRole([]string{"ROLE_ADMIN", "ROLE_TEACHER"}, "roles") {
		t.Error("expected true when user has all of the allowed roles")
	}
}

func TestHasAllowedRole_AllowedRolesWithWhitespace(t *testing.T) {
	setAllowedRoles(t, " ROLE_ADMIN , ROLE_TEACHER ")
	if !hasAllowedRole([]string{"ROLE_TEACHER"}, "roles") {
		t.Error("expected true when allowed roles config has surrounding whitespace")
	}
}

// --- no matching role ---

func TestHasAllowedRole_NoMatchingRole(t *testing.T) {
	setAllowedRoles(t, "ROLE_ADMIN")
	if hasAllowedRole([]string{"ROLE_STUDENT"}, "roles") {
		t.Error("expected false when user does not have any allowed role")
	}
}

func TestHasAllowedRole_EmptyUserRoles_WithConfig(t *testing.T) {
	setAllowedRoles(t, "ROLE_ADMIN")
	if hasAllowedRole([]string{}, "roles") {
		t.Error("expected false when user has no roles and OAUTH_ALLOWED_ROLES is set")
	}
}

func TestHasAllowedRole_NilUserRoles_WithConfig(t *testing.T) {
	setAllowedRoles(t, "ROLE_ADMIN")
	if hasAllowedRole(nil, "roles") {
		t.Error("expected false when user roles slice is nil and OAUTH_ALLOWED_ROLES is set")
	}
}

func TestHasAllowedRole_CaseSensitiveNoMatch(t *testing.T) {
	setAllowedRoles(t, "ROLE_ADMIN")
	// Role names are matched case-sensitively.
	if hasAllowedRole([]string{"role_admin"}, "roles") {
		t.Error("expected false: role matching must be case-sensitive")
	}
}

func TestHasAllowedRole_MultipleDefinedNoMatch(t *testing.T) {
    setAllowedRoles(t, "ROLE_ADMIN,ROLE_TEACHER")
    if hasAllowedRole([]string{"ROLE_GUEST"}, "roles") {
        t.Error("expected false when user has multiple roles but none match the allowed roles")
    }
}

func TestHasAllowedRole_MultipleNoMatch(t *testing.T) {
    setAllowedRoles(t, "ROLE_ADMIN,ROLE_TEACHER")
    if hasAllowedRole([]string{"ROLE_STUDENT", "ROLE_GUEST"}, "roles") {
        t.Error("expected false when user has multiple roles but none match the allowed roles")
    }
}

// --- edge cases ---

func TestHasAllowedRole_ConfigWithOnlyCommas_AllowsAll(t *testing.T) {
	// A config of only separators produces no valid role entries → allow all.
	setAllowedRoles(t, ",,, ,")
	if !hasAllowedRole([]string{}, "roles") {
		t.Error("expected true when config contains only separators (no valid roles)")
	}
}


