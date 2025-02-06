package middlewares_test

import (
	"net/http"
	"testing"

	"github.com/Spicy-Bush/fider-tarkov-community/app/metrics"
	"github.com/Spicy-Bush/fider-tarkov-community/app/middlewares"
	. "github.com/Spicy-Bush/fider-tarkov-community/app/pkg/assert"
	"github.com/Spicy-Bush/fider-tarkov-community/app/pkg/env"
	"github.com/Spicy-Bush/fider-tarkov-community/app/pkg/mock"
	"github.com/Spicy-Bush/fider-tarkov-community/app/pkg/web"
	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
)

func TestInstrumentation_Disabled(t *testing.T) {
	RegisterT(t)

	server := mock.NewServer()
	env.Config.Metrics.Enabled = false
	server.Use(middlewares.Instrumentation())
	handler := func(c *web.Context) error {
		return c.NoContent(http.StatusOK)
	}

	status, _ := server.Execute(handler)

	Expect(status).Equals(http.StatusOK)
	Expect(getCounterValue(metrics.HttpRequests, "200", "GET /")).Equals(float64(0))
}

func TestInstrumentation_Enabled(t *testing.T) {
	RegisterT(t)

	server := mock.NewServer()
	server.Use(middlewares.Instrumentation())
	handler := func(c *web.Context) error {
		return c.NoContent(http.StatusOK)
	}

	status, _ := server.Execute(handler)

	Expect(status).Equals(http.StatusOK)
	Expect(getCounterValue(metrics.HttpRequests, "200", "GET /")).Equals(float64(1))
}

func getCounterValue(metric *prometheus.CounterVec, lvs ...string) float64 {
	var m = &dto.Metric{}
	if err := metric.WithLabelValues(lvs...).Write(m); err != nil {
		return 0
	}
	return m.Counter.GetValue()
}
