package entity_test

import (
	"encoding/json"
	"testing"

	"github.com/getfider/fider/app/models/entity"
	. "github.com/getfider/fider/app/pkg/assert"
)

func TestUserWithEmail_MarshalJSON(t *testing.T) {

	RegisterT(t)
	user := entity.UserWithEmail{
		User: &entity.User{
			ID:     1,
			Name:   "John Doe",
			Email:  "johndoe@example.com",
			Role:   1,
			Status: 1,
		},
	}

	expectedJSON := `{"id":1,"name":"John Doe","role":"visitor","status":"active","email":"johndoe@example.com"}`

	jsonData, err := json.Marshal(user)
	if err != nil {
		t.Errorf("Failed to marshal user to JSON: %v", err)
	}

	Expect(string(jsonData)).Equals(expectedJSON)

}

func TestUser_MarshalJSON(t *testing.T) {

	RegisterT(t)
	user := entity.User{
		ID:     1,
		Name:   "John Doe",
		Email:  "johndoe@example.com",
		Role:   1,
		Status: 1,
	}

	expectedJSON := `{"id":1,"name":"John Doe","role":"visitor","status":"active"}`

	jsonData, err := json.Marshal(user)
	if err != nil {
		t.Errorf("Failed to marshal user to JSON: %v", err)
	}

	Expect(string(jsonData)).Equals(expectedJSON)

}
