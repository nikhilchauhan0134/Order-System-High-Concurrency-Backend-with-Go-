package metrics

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// define a map to hold counters dynamically
	counters   = map[string]prometheus.Counter{}
	histograms = map[string]prometheus.Histogram{}
)

// IncrementCounter increments a Prometheus counter by 1
func IncrementCounter(name string) {
	counter, ok := counters[name]
	if !ok {
		counter = promauto.NewCounter(prometheus.CounterOpts{
			Name: name,
			Help: name + " counter",
		})
		counters[name] = counter
	}
	counter.Inc()
}

// ObserveDuration records a duration for a metric
func ObserveDuration(name string, duration float64) {
	hist, ok := histograms[name]
	if !ok {
		hist = promauto.NewHistogram(prometheus.HistogramOpts{
			Name:    name,
			Help:    name + " duration histogram",
			Buckets: prometheus.DefBuckets,
		})
		histograms[name] = hist
	}
	hist.Observe(duration)
}

// Timer helper to measure duration easily
func Timer(name string) func() {
	start := time.Now()
	return func() {
		ObserveDuration(name, time.Since(start).Seconds())
	}
}
