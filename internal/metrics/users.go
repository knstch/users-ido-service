package metrics

import (
	"sync"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	registerOnce sync.Once

	usersCreatedTotal = prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: "users_service",
		Subsystem: "users",
		Name:      "created_total",
		Help:      "Total number of newly created users",
	})
)

func register() {
	registerOnce.Do(func() {
		prometheus.MustRegister(usersCreatedTotal)
	})
}

// IncUsersCreated increments the "users created" Prometheus counter.
func IncUsersCreated() {
	register()
	usersCreatedTotal.Inc()
}

