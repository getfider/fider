package handlers

// Internal test file for hasAllowedRole — uses package handlers (not handlers_test)
// so the unexported function is accessible without exporting it.

import (
	"testing"
)

// --- empty / unconfigured allowed roles ---

func TestHasAllowedRole_EmptyConfig_AllowsAll(t *testing.T) {
	if !hasAllowedRole([]string{"ROLE_ADMIN"}, "roles", "") {
		t.Error("expected true when allowedRoles is empty")
	}
}

func TestHasAllowedRole_EmptyConfig_AllowsEmptyRoles(t *testing.T) {
	if !hasAllowedRole([]string{}, "roles", "") {
		t.Error("expected true when allowedRoles is empty and user has no roles")
	}
}

func TestHasAllowedRole_WhitespaceOnlyConfig_AllowsAll(t *testing.T) {
	if !hasAllowedRole([]string{"ROLE_ADMIN"}, "roles", "   ") {
		t.Error("expected true when allowedRoles is whitespace only")
	}
}

// --- provider has no JSONUserRolesPath (skip check) ---

func TestHasAllowedRole_NoRolesPath_SkipsCheck(t *testing.T) {
	// Provider without a roles path must always be allowed through,
	// regardless of whether the user carries matching roles or not.
	if !hasAllowedRole([]string{}, "", "ROLE_ADMIN") {
		t.Error("expected true when provider has no JSONUserRolesPath (check should be skipped)")
	}
}

func TestHasAllowedRole_WhitespaceRolesPath_SkipsCheck(t *testing.T) {
	if !hasAllowedRole([]string{}, "   ", "ROLE_ADMIN") {
		t.Error("expected true when JSONUserRolesPath is whitespace only (check should be skipped)")
	}
}

// --- matching role ---

func TestHasAllowedRole_MatchingSingleRole(t *testing.T) {
	if !hasAllowedRole([]string{"ROLE_ADMIN"}, "roles", "ROLE_ADMIN") {
		t.Error("expected true when user has the required role")
	}
}

func TestHasAllowedRole_MatchingOneOfMultipleAllowed(t *testing.T) {
	if !hasAllowedRole([]string{"ROLE_TEACHER"}, "roles", "ROLE_ADMIN,ROLE_TEACHER") {
		t.Error("expected true when user has one of the allowed roles")
	}
}

func TestHasAllowedRole_MatchingMultipleAllowed(t *testing.T) {
	if !hasAllowedRole([]string{"ROLE_ADMIN", "ROLE_TEACHER"}, "roles", "ROLE_ADMIN,ROLE_TEACHER") {
		t.Error("expected true when user has all of the allowed roles")
	}
}

func TestHasAllowedRole_AllowedRolesWithWhitespace(t *testing.T) {
	if !hasAllowedRole([]string{"ROLE_TEACHER"}, "roles", " ROLE_ADMIN , ROLE_TEACHER ") {
		t.Error("expected true when allowed roles config has surrounding whitespace")
	}
}

// --- no matching role ---

func TestHasAllowedRole_NoMatchingRole(t *testing.T) {
	if hasAllowedRole([]string{"ROLE_STUDENT"}, "roles", "ROLE_ADMIN") {
		t.Error("expected false when user does not have any allowed role")
	}
}

func TestHasAllowedRole_EmptyUserRoles_WithConfig(t *testing.T) {
	if hasAllowedRole([]string{}, "roles", "ROLE_ADMIN") {
		t.Error("expected false when user has no roles and allowedRoles is set")
	}
}

func TestHasAllowedRole_NilUserRoles_WithConfig(t *testing.T) {
	if hasAllowedRole(nil, "roles", "ROLE_ADMIN") {
		t.Error("expected false when user roles slice is nil and allowedRoles is set")
	}
}

func TestHasAllowedRole_CaseSensitiveNoMatch(t *testing.T) {
	// Role names are matched case-sensitively.
	if hasAllowedRole([]string{"role_admin"}, "roles", "ROLE_ADMIN") {
		t.Error("expected false: role matching must be case-sensitive")
	}
}

func TestHasAllowedRole_MultipleDefinedNoMatch(t *testing.T) {
	if hasAllowedRole([]string{"ROLE_GUEST"}, "roles", "ROLE_ADMIN,ROLE_TEACHER") {
		t.Error("expected false when user role does not match any allowed role")
	}
}

func TestHasAllowedRole_MultipleNoMatch(t *testing.T) {
	if hasAllowedRole([]string{"ROLE_STUDENT", "ROLE_GUEST"}, "roles", "ROLE_ADMIN,ROLE_TEACHER") {
		t.Error("expected false when user has multiple roles but none match the allowed roles")
	}
}

// --- edge cases ---

func TestHasAllowedRole_ConfigWithOnlyCommas_AllowsAll(t *testing.T) {
	// A config of only separators produces no valid role entries → allow all.
	if !hasAllowedRole([]string{}, "roles", ",,, ,") {
		t.Error("expected true when config contains only separators (no valid roles)")
	}
}

