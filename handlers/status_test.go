package handlers

import (
	"database/sql"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/WeCanHearYou/wchy-api/context"
	"github.com/WeCanHearYou/wchy-api/services"
	"github.com/stretchr/testify/assert"
)

func TestStatusHandler(t *testing.T) {
	db, _ := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	ctx := context.WchyContext{
		Health: services.NewPostgresHealthCheckService(db),
	}
	router := GetMainEngine(ctx)

	req, err := http.NewRequest("GET", "/status", nil)
	if err != nil {
		fmt.Println(err)
	}

	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	assert.Equal(t, resp.Code, 200)
}
