package dbEntities_test

import (
	"context"
	"database/sql"
	"net/url"
	"testing"

	"github.com/getfider/fider/app"
	"github.com/getfider/fider/app/models/entity"
	"github.com/getfider/fider/app/models/enum"
	"github.com/getfider/fider/app/pkg/web"
	"github.com/getfider/fider/app/services/sqlstore/dbEntities"
)

func TestUserToModel(t *testing.T) {
	// Create a proper context with web.Request
	u, _ := url.Parse("http://test.fider.io")
	req := web.Request{URL: u}
	ctx := context.WithValue(context.Background(), app.RequestCtxKey, req)

	// Create a test dbEntities.User
	dbUser := &dbEntities.User{
		ID:            sql.NullInt64{Int64: 1, Valid: true},
		Name:          sql.NullString{String: "John Doe", Valid: true},
		Email:         sql.NullString{String: "john@example.com", Valid: true},
		Role:          sql.NullInt64{Int64: int64(enum.RoleAdministrator), Valid: true},
		Status:        sql.NullInt64{Int64: int64(enum.UserActive), Valid: true},
		AvatarType:    sql.NullInt64{Int64: int64(enum.AvatarTypeGravatar), Valid: true},
		AvatarBlobKey: sql.NullString{String: "", Valid: true},
		IsTrusted:     sql.NullBool{Bool: true, Valid: true},
		Providers: []*dbEntities.UserProvider{
			{
				Name: sql.NullString{String: "google", Valid: true},
				UID:  sql.NullString{String: "123456", Valid: true},
			},
		},
	}

	// Convert to entity.User
	entityUser := dbUser.ToModel(ctx)

	// Verify conversion
	if entityUser == nil {
		t.Fatal("ToModel returned nil")
	}

	if entityUser.ID != 1 {
		t.Errorf("Expected ID 1, got %d", entityUser.ID)
	}

	if entityUser.Name != "John Doe" {
		t.Errorf("Expected Name 'John Doe', got '%s'", entityUser.Name)
	}

	if entityUser.Email != "john@example.com" {
		t.Errorf("Expected Email 'john@example.com', got '%s'", entityUser.Email)
	}

	if entityUser.Role != enum.RoleAdministrator {
		t.Errorf("Expected Role Administrator, got %v", entityUser.Role)
	}

	if entityUser.Status != enum.UserActive {
		t.Errorf("Expected Status Active, got %v", entityUser.Status)
	}

	if !entityUser.IsTrusted {
		t.Error("Expected IsTrusted to be true")
	}

	if len(entityUser.Providers) != 1 {
		t.Fatalf("Expected 1 provider, got %d", len(entityUser.Providers))
	}

	if entityUser.Providers[0].Name != "google" {
		t.Errorf("Expected provider name 'google', got '%s'", entityUser.Providers[0].Name)
	}

	if entityUser.Providers[0].UID != "123456" {
		t.Errorf("Expected provider UID '123456', got '%s'", entityUser.Providers[0].UID)
	}
}

// TestUserToModel_PreservesPersistedTenant verifies that a user loaded from the database
// keeps the tenant persisted on its own row even when a *different* tenant is present in
// the request context. Adopting the context tenant here would make the downstream
// tenant-equality authorization check tautological (cross-tenant token replay).
func TestUserToModel_PreservesPersistedTenant(t *testing.T) {
	u, _ := url.Parse("http://victim.fider.io")
	req := web.Request{URL: u}
	ctx := context.WithValue(context.Background(), app.RequestCtxKey, req)
	// Request is scoped to the victim tenant (id 2)...
	ctx = context.WithValue(ctx, app.TenantCtxKey, &entity.Tenant{ID: 2, Name: "Victim"})

	// ...but the user row is persisted under the attacker tenant (id 1).
	dbUser := &dbEntities.User{
		ID:     sql.NullInt64{Int64: 1, Valid: true},
		Name:   sql.NullString{String: "Attacker", Valid: true},
		Role:   sql.NullInt64{Int64: int64(enum.RoleAdministrator), Valid: true},
		Status: sql.NullInt64{Int64: int64(enum.UserActive), Valid: true},
		Tenant: &dbEntities.Tenant{ID: 1},
	}

	entityUser := dbUser.ToModel(ctx)

	if entityUser.Tenant == nil {
		t.Fatal("expected Tenant to be populated from the persisted association")
	}
	if entityUser.Tenant.ID != 1 {
		t.Errorf("expected persisted Tenant ID 1, got %d", entityUser.Tenant.ID)
	}
}

// TestUserToModel_SameTenantUsesRequestTenant verifies that when the request tenant matches
// the persisted tenant, the fully-populated request tenant is used for presentation data.
func TestUserToModel_SameTenantUsesRequestTenant(t *testing.T) {
	u, _ := url.Parse("http://demo.fider.io")
	req := web.Request{URL: u}
	ctx := context.WithValue(context.Background(), app.RequestCtxKey, req)
	ctx = context.WithValue(ctx, app.TenantCtxKey, &entity.Tenant{ID: 1, Name: "Demo"})

	dbUser := &dbEntities.User{
		ID:     sql.NullInt64{Int64: 1, Valid: true},
		Name:   sql.NullString{String: "Jon Snow", Valid: true},
		Role:   sql.NullInt64{Int64: int64(enum.RoleAdministrator), Valid: true},
		Status: sql.NullInt64{Int64: int64(enum.UserActive), Valid: true},
		Tenant: &dbEntities.Tenant{ID: 1},
	}

	entityUser := dbUser.ToModel(ctx)

	if entityUser.Tenant == nil || entityUser.Tenant.ID != 1 {
		t.Fatalf("expected Tenant ID 1, got %v", entityUser.Tenant)
	}
	if entityUser.Tenant.Name != "Demo" {
		t.Errorf("expected fully-populated request tenant (Name 'Demo'), got '%s'", entityUser.Tenant.Name)
	}
}

func TestUserToModel_Nil(t *testing.T) {
	var dbUser *dbEntities.User
	entityUser := dbUser.ToModel(context.Background())

	if entityUser != nil {
		t.Error("Expected ToModel on nil to return nil")
	}
}
