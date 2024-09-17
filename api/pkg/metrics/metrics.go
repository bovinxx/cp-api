package metrics

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

type Metrics struct {
	RequestCounter  prometheus.Counter
	RequestDuration *prometheus.HistogramVec
}

func NewMetrics() *Metrics {
	var buckets []float64
	cnt := 0.0
	for i := 0; i < 100; i++ {
		buckets = append(buckets, cnt)
		cnt += 10.0
	}

	requestDuration := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "server_request_dur",
		Help:    "Histogram of response time for handler in seconds",
		Buckets: buckets,
	}, []string{"route"})

	requestCounter := prometheus.NewCounter(prometheus.CounterOpts{
		Name: "server_request_counter",
		Help: "How many HTTP requests processed",
	})
	prometheus.MustRegister(requestCounter)
	prometheus.MustRegister(requestDuration)
	return &Metrics{
		RequestCounter:  requestCounter,
		RequestDuration: requestDuration,
	}
}

func (m *Metrics) Record(start time.Time, handle string) {
	duration := float64(time.Since(start).Milliseconds())
	m.RequestDuration.WithLabelValues(handle).Observe(duration)
	m.RequestCounter.Inc()
}
