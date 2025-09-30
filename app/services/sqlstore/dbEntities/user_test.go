package dbEntities_test

import (
	"context"
	"database/sql"
	"net/url"
	"testing"

	"github.com/getfider/fider/app"
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
		IsVerified:    sql.NullBool{Bool: true, Valid: true},
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

	if !entityUser.IsVerified {
		t.Error("Expected IsVerified to be true")
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

func TestUserToModel_Nil(t *testing.T) {
	var dbUser *dbEntities.User
	entityUser := dbUser.ToModel(context.Background())

	if entityUser != nil {
		t.Error("Expected ToModel on nil to return nil")
	}
}
