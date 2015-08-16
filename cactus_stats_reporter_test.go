package bark_test

import (
	"testing"
	"net"
	"github.com/stretchr/testify/assert"
	"github.com/cactus/go-statsd-client/statsd"
	"github.com/uber/bark"
	"time"
)

func TestBarkCactusStatsReporter(t *testing.T) {
	conn, err := net.Listen("tcp", "127.0.0.1:0")
	assert.NoError(t, err)

	defer conn.Close()

	cactusStatter, err := statsd.New(conn.(*net.TCPListener).Addr().String(), "barktest")
	assert.NoError(t, err)

	barkClient := bark.NewStatsReporterFromCactus(cactusStatter)
	barkClient.IncCounter("foo", map[string]string{"tag":"val"}, 1)
	barkClient.RecordTimer("bar", map[string]string{"tag":"val"}, time.Duration(10))
	barkClient.UpdateGauge("baz", map[string]string{"tag":"val"}, 100)
}