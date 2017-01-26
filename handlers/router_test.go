package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"

	"github.com/WeCanHearYou/wchy-api/context"
	"github.com/WeCanHearYou/wchy-api/models"
	"github.com/WeCanHearYou/wchy-api/services"
	"github.com/jmoiron/jsonq"
)

func makeRequest(method, url string) (int, *jsonq.JsonQuery) {
	ctx := context.WchyContext{
		Health: &services.InMemoryHealthCheckService{Status: false},
		Tenant: &services.InMemoryTenantService{Tenants: []*models.Tenant{
			&models.Tenant{ID: 1, Name: "Orange Inc.", Domain: "orange"},
			&models.Tenant{ID: 2, Name: "The Triathlon Shop", Domain: "trishop"},
		}},
		Settings: context.WchySettings{
			BuildTime: "today",
		},
	}
	router := GetMainEngine(ctx)

	request, _ := http.NewRequest(method, url, nil)
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	var query *jsonq.JsonQuery
	if response.Result().StatusCode == 200 {
		var data interface{}
		decoder := json.NewDecoder(response.Body)
		decoder.Decode(&data)
		query = jsonq.NewQuery(data)
	}

	return response.Result().StatusCode, query
}
