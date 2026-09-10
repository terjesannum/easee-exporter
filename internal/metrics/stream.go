package metrics

import "github.com/prometheus/client_golang/prometheus"

// StreamSource is the part of the stream this collector reports on.
type StreamSource interface {
	Connected() bool
	Reconnects() int64
}

// StreamCollector exposes stream health. Without it a stream that is up but
// silent is indistinguishable from a fleet of idle chargers, because both
// look like metrics that simply stop changing.
type StreamCollector struct {
	stream     StreamSource
	connected  *prometheus.Desc
	reconnects *prometheus.Desc
}

func NewStreamCollector(stream StreamSource) *StreamCollector {
	return &StreamCollector{
		stream: stream,
		connected: prometheus.NewDesc(
			"easee_exporter_stream_connected",
			"Whether the observation stream is currently connected",
			nil,
			nil,
		),
		reconnects: prometheus.NewDesc(
			"easee_exporter_stream_connects_total",
			"Number of times the observation stream has connected, including the first",
			nil,
			nil,
		),
	}
}

func (c *StreamCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.connected
	ch <- c.reconnects
}

func (c *StreamCollector) Collect(ch chan<- prometheus.Metric) {
	connected := 0.0
	if c.stream.Connected() {
		connected = 1
	}
	ch <- prometheus.MustNewConstMetric(c.connected, prometheus.GaugeValue, connected)
	ch <- prometheus.MustNewConstMetric(c.reconnects, prometheus.CounterValue, float64(c.stream.Reconnects()))
}
