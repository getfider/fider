package metrics

import (
	"github.com/getfider/fider/app/pkg/env"
	"github.com/prometheus/client_golang/prometheus"
)

var TotalTenants = prometheus.NewCounter(
	prometheus.CounterOpts{
		Name: "fider_tenants_total",
		Help: "Number of Fider tenants.",
	},
)

var TotalPosts = prometheus.NewCounter(
	prometheus.CounterOpts{
		Name: "fider_posts_total",
		Help: "Number of Fider posts.",
	},
)

var TotalComments = prometheus.NewCounter(
	prometheus.CounterOpts{
		Name: "fider_comments_total",
		Help: "Number of Fider comments.",
	},
)

var TotalVotes = prometheus.NewCounter(
	prometheus.CounterOpts{
		Name: "fider_votes_total",
		Help: "Number of Fider votes.",
	},
)

var fiderInfo = prometheus.NewGauge(
	prometheus.GaugeOpts{
		Name:        "fider_info",
		Help:        "Information about Fider environment.",
		ConstLabels: prometheus.Labels{"version": env.Version()},
	},
)

func init() {
	fiderInfo.Inc()
	prometheus.MustRegister(TotalTenants, TotalPosts, TotalComments, TotalVotes, fiderInfo)
}
