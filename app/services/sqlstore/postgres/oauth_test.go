package postgres_test

import (
	"testing"

	"github.com/getfider/fider/app/models/cmd"
	"github.com/getfider/fider/app/models/dto"
	"github.com/getfider/fider/app/pkg/bus"
	. "github.com/getfider/fider/app/pkg/assert"
)

// stampRow is used to read security_stamp directly from the DB.
type stampRow struct {
	ID            int    `db:"id"`
	SecurityStamp string `db:"security_stamp"`
}

// getSecurityStamps returns a map of userID → security_stamp for all users in the tenant.
func getSecurityStamps(tenantID int) map[int]string {
	var rows []*stampRow
	_ = trx.Select(&rows, "SELECT id, security_stamp FROM users WHERE tenant_id = $1", tenantID)
	result := make(map[int]string)
	for _, r := range rows {
		result[r.ID] = r.SecurityStamp
	}
	return result
}

// newOAuthConfig builds a minimal SaveCustomOAuthConfig suitable for INSERT (ID = 0).
// Callers set JSONUserRolesPath and AllowedRoles to control whether role access is active.
func newOAuthConfig(rolesPath, allowedRoles string) *cmd.SaveCustomOAuthConfig {
	return &cmd.SaveCustomOAuthConfig{
		Logo:              &dto.ImageUpload{},
		Provider:          "_TEST_ROLES",
		DisplayName:       "Test Roles Provider",
		ClientID:          "client-id",
		ClientSecret:      "client-secret",
		AuthorizeURL:      "http://provider/authorize",
		TokenURL:          "http://provider/token",
		ProfileURL:        "http://provider/profile",
		JSONUserIDPath:    "id",
		JSONUserNamePath:  "name",
		JSONUserEmailPath: "email",
		JSONUserRolesPath: rolesPath,
		AllowedRoles:      allowedRoles,
	}
}

// TestOAuthConfig_SecurityStamp_ActiveToActive checks that when allowed_roles is
// changed between two non-empty values (and json_user_roles_path is set), security
// stamps are rotated for all users except the one making the change.
func TestOAuthConfig_SecurityStamp_ActiveToActive(t *testing.T) {
	SetupDatabaseTest(t)
	defer TeardownDatabaseTest()

	// Insert initial active config (AllowedRoles="admin", JSONUserRolesPath="roles")
	insertCmd := newOAuthConfig("roles", "admin")
	err := bus.Dispatch(jonSnowCtx, insertCmd)
	Expect(err).IsNil()

	stampsBefore := getSecurityStamps(demoTenant.ID)

	// Update: change AllowedRoles to a different non-empty value — still active.
	updateCmd := newOAuthConfig("roles", "superadmin")
	updateCmd.ID = insertCmd.ID
	err = bus.Dispatch(jonSnowCtx, updateCmd)
	Expect(err).IsNil()

	stampsAfter := getSecurityStamps(demoTenant.ID)

	// Other users' stamps must have been rotated.
	Expect(stampsAfter[aryaStark.ID]).NotEquals(stampsBefore[aryaStark.ID])
	Expect(stampsAfter[sansaStark.ID]).NotEquals(stampsBefore[sansaStark.ID])

	// The acting user's stamp must NOT have changed.
	Expect(stampsAfter[jonSnow.ID]).Equals(stampsBefore[jonSnow.ID])
}

// TestOAuthConfig_SecurityStamp_InactiveToActive checks that when a provider
// transitions from inactive (empty AllowedRoles) to active, security stamps are
// rotated for all users except the one making the change.
func TestOAuthConfig_SecurityStamp_InactiveToActive(t *testing.T) {
	SetupDatabaseTest(t)
	defer TeardownDatabaseTest()

	// Insert initial inactive config: JSONUserRolesPath set but AllowedRoles empty.
	insertCmd := newOAuthConfig("roles", "")
	err := bus.Dispatch(jonSnowCtx, insertCmd)
	Expect(err).IsNil()

	stampsBefore := getSecurityStamps(demoTenant.ID)

	// Update: add AllowedRoles — transitions to active.
	updateCmd := newOAuthConfig("roles", "admin")
	updateCmd.ID = insertCmd.ID
	err = bus.Dispatch(jonSnowCtx, updateCmd)
	Expect(err).IsNil()

	stampsAfter := getSecurityStamps(demoTenant.ID)

	// Other users' stamps must have been rotated.
	Expect(stampsAfter[aryaStark.ID]).NotEquals(stampsBefore[aryaStark.ID])
	Expect(stampsAfter[sansaStark.ID]).NotEquals(stampsBefore[sansaStark.ID])

	// The acting user's stamp must NOT have changed.
	Expect(stampsAfter[jonSnow.ID]).Equals(stampsBefore[jonSnow.ID])
}

// TestOAuthConfig_SecurityStamp_ActiveToInactive_ClearAllowedRoles checks that
// clearing AllowedRoles (deactivating role access) does NOT rotate security stamps,
// because all users are allowed to log in when restrictions are lifted.
func TestOAuthConfig_SecurityStamp_ActiveToInactive_ClearAllowedRoles(t *testing.T) {
	SetupDatabaseTest(t)
	defer TeardownDatabaseTest()

	// Insert initial active config.
	insertCmd := newOAuthConfig("roles", "admin")
	err := bus.Dispatch(jonSnowCtx, insertCmd)
	Expect(err).IsNil()

	stampsBefore := getSecurityStamps(demoTenant.ID)

	// Update: clear AllowedRoles — transitions to inactive (JSONUserRolesPath still set).
	updateCmd := newOAuthConfig("roles", "")
	updateCmd.ID = insertCmd.ID
	err = bus.Dispatch(jonSnowCtx, updateCmd)
	Expect(err).IsNil()

	stampsAfter := getSecurityStamps(demoTenant.ID)

	// No stamps should have changed.
	Expect(stampsAfter[jonSnow.ID]).Equals(stampsBefore[jonSnow.ID])
	Expect(stampsAfter[aryaStark.ID]).Equals(stampsBefore[aryaStark.ID])
	Expect(stampsAfter[sansaStark.ID]).Equals(stampsBefore[sansaStark.ID])
}

// TestOAuthConfig_SecurityStamp_ActiveToInactive_ClearRolesPath checks that
// when json_user_roles_path is cleared (making role evaluation impossible) while
// AllowedRoles also changes, security stamps are NOT rotated.
func TestOAuthConfig_SecurityStamp_ActiveToInactive_ClearRolesPath(t *testing.T) {
	SetupDatabaseTest(t)
	defer TeardownDatabaseTest()

	// Insert initial active config.
	insertCmd := newOAuthConfig("roles", "admin")
	err := bus.Dispatch(jonSnowCtx, insertCmd)
	Expect(err).IsNil()

	stampsBefore := getSecurityStamps(demoTenant.ID)

	// Update: clear JSONUserRolesPath and change AllowedRoles — inactive because no path.
	updateCmd := newOAuthConfig("", "superadmin")
	updateCmd.ID = insertCmd.ID
	err = bus.Dispatch(jonSnowCtx, updateCmd)
	Expect(err).IsNil()

	stampsAfter := getSecurityStamps(demoTenant.ID)

	// No stamps should have changed.
	Expect(stampsAfter[jonSnow.ID]).Equals(stampsBefore[jonSnow.ID])
	Expect(stampsAfter[aryaStark.ID]).Equals(stampsBefore[aryaStark.ID])
	Expect(stampsAfter[sansaStark.ID]).Equals(stampsBefore[sansaStark.ID])
}

// TestOAuthConfig_SecurityStamp_ActiveToInactive_ClearBoth checks that clearing
// both AllowedRoles and json_user_roles_path does NOT rotate security stamps.
func TestOAuthConfig_SecurityStamp_ActiveToInactive_ClearBoth(t *testing.T) {
	SetupDatabaseTest(t)
	defer TeardownDatabaseTest()

	// Insert initial active config.
	insertCmd := newOAuthConfig("roles", "admin")
	err := bus.Dispatch(jonSnowCtx, insertCmd)
	Expect(err).IsNil()

	stampsBefore := getSecurityStamps(demoTenant.ID)

	// Update: clear both fields — fully inactive.
	updateCmd := newOAuthConfig("", "")
	updateCmd.ID = insertCmd.ID
	err = bus.Dispatch(jonSnowCtx, updateCmd)
	Expect(err).IsNil()

	stampsAfter := getSecurityStamps(demoTenant.ID)

	// No stamps should have changed.
	Expect(stampsAfter[jonSnow.ID]).Equals(stampsBefore[jonSnow.ID])
	Expect(stampsAfter[aryaStark.ID]).Equals(stampsBefore[aryaStark.ID])
	Expect(stampsAfter[sansaStark.ID]).Equals(stampsBefore[sansaStark.ID])
}

// TestOAuthConfig_SecurityStamp_InactiveToActive_ViaRolesPath checks that when
// AllowedRoles is already set but json_user_roles_path is empty (so access control
// is inactive), adding json_user_roles_path (without changing AllowedRoles) DOES
// rotate security stamps, because the provider just became active.
func TestOAuthConfig_SecurityStamp_InactiveToActive_ViaRolesPath(t *testing.T) {
	SetupDatabaseTest(t)
	defer TeardownDatabaseTest()

	// Insert initial inactive config: AllowedRoles set but JSONUserRolesPath empty.
	insertCmd := newOAuthConfig("", "admin")
	err := bus.Dispatch(jonSnowCtx, insertCmd)
	Expect(err).IsNil()

	stampsBefore := getSecurityStamps(demoTenant.ID)

	// Update: add JSONUserRolesPath only — AllowedRoles stays "admin", transitions to active.
	updateCmd := newOAuthConfig("roles", "admin")
	updateCmd.ID = insertCmd.ID
	err = bus.Dispatch(jonSnowCtx, updateCmd)
	Expect(err).IsNil()

	stampsAfter := getSecurityStamps(demoTenant.ID)

	// Other users' stamps must have been rotated.
	Expect(stampsAfter[aryaStark.ID]).NotEquals(stampsBefore[aryaStark.ID])
	Expect(stampsAfter[sansaStark.ID]).NotEquals(stampsBefore[sansaStark.ID])

	// The acting user's stamp must NOT have changed.
	Expect(stampsAfter[jonSnow.ID]).Equals(stampsBefore[jonSnow.ID])
}

// TestOAuthConfig_SecurityStamp_InactiveToInactive checks that updating a provider
// that has role access disabled (no AllowedRoles, no JSONUserRolesPath) without
// enabling it does NOT rotate security stamps.
func TestOAuthConfig_SecurityStamp_InactiveToInactive(t *testing.T) {
	SetupDatabaseTest(t)
	defer TeardownDatabaseTest()

	// Insert initial inactive config — no role fields set at all.
	insertCmd := newOAuthConfig("", "")
	err := bus.Dispatch(jonSnowCtx, insertCmd)
	Expect(err).IsNil()

	stampsBefore := getSecurityStamps(demoTenant.ID)

	// Update: change an unrelated field — role access remains inactive.
	updateCmd := newOAuthConfig("", "")
	updateCmd.ID = insertCmd.ID
	updateCmd.DisplayName = "Updated Display Name"
	err = bus.Dispatch(jonSnowCtx, updateCmd)
	Expect(err).IsNil()

	stampsAfter := getSecurityStamps(demoTenant.ID)

	// No stamps should have changed.
	Expect(stampsAfter[jonSnow.ID]).Equals(stampsBefore[jonSnow.ID])
	Expect(stampsAfter[aryaStark.ID]).Equals(stampsBefore[aryaStark.ID])
	Expect(stampsAfter[sansaStark.ID]).Equals(stampsBefore[sansaStark.ID])
}



