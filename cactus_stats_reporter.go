package bark

import (
	"github.com/cactus/go-statsd-client/statsd"
	"time"
)

type barkCactusStatsReporter struct {
	delegate statsd.Statter
}

func newBarkCactusStatsReporter(wrappedObject statsd.Statter) StatsReporter {
	return &barkCactusStatsReporter{delegate: wrappedObject}
}

func (s *barkCactusStatsReporter) IncCounter(name string, tags map[string]string, value int64) {
	s.delegate.Inc(name, value, 1.0)
}

func (s *barkCactusStatsReporter) UpdateGauge(name string, tags map[string]string, value int64) {
	s.delegate.Gauge(name, value, 1.0)
}

func (s *barkCactusStatsReporter) RecordTimer(name string, tags map[string]string, d time.Duration) {
	s.delegate.TimingDuration(name, d, 1.0)
}
