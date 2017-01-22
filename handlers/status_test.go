package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/WeCanHearYou/wchy-api/context"
	"github.com/WeCanHearYou/wchy-api/services"
	"github.com/WeCanHearYou/wchy-api/util"
	"github.com/stretchr/testify/assert"
)

func TestStatusHandler(t *testing.T) {
	ctx := context.WchyContext{
		Health: services.NewInMemoryHealthCheckService(false),
		Settings: context.WchySettings{
			BuildTime: "today",
		},
	}
	router := GetMainEngine(ctx)

	requeest, _ := http.NewRequest("GET", "/status", nil)
	response := httptest.NewRecorder()
	router.ServeHTTP(response, requeest)
	data := util.NewJSONObject(response.Body.Bytes())

	assert.Equal(t, 200, response.Result().StatusCode)
	assert.Equal(t, "today", data.GetString("build"))
}
