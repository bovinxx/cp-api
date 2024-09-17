package metrics

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

type Metrics struct {
	RequestCounter  prometheus.Counter
	RequestDuration prometheus.Histogram
	TranslatorsUsed *prometheus.GaugeVec
}

func NewMetrics() *Metrics {
	var buckets []float64
	for i := 0.0; i < 10000.0; i += 10.0 {
		buckets = append(buckets, i)
	}

	requestDuration := prometheus.NewHistogram(prometheus.HistogramOpts{
		Name:    "processor_request_dur",
		Help:    "",
		Buckets: buckets,
	})

	requestCounter := prometheus.NewCounter(prometheus.CounterOpts{
		Name: "processor_request_counter",
		Help: "How many requests processed",
	})

	translatorsUsed := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "processor_translators_used",
		Help: "",
	}, []string{"translator"})

	prometheus.MustRegister(requestCounter)
	prometheus.MustRegister(requestDuration)
	prometheus.MustRegister(translatorsUsed)
	return &Metrics{
		RequestCounter:  requestCounter,
		RequestDuration: requestDuration,
		TranslatorsUsed: translatorsUsed,
	}
}

func (m *Metrics) Record(start time.Time, translator string) {
	duration := float64(time.Since(start).Milliseconds())
	m.RequestDuration.Observe(duration)
	m.RequestCounter.Inc()
}
