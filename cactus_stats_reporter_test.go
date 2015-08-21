package bark_test

// Copyright (c) 2015 Uber Technologies, Inc.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

import (
	"net"
	"testing"
	"time"

	"github.com/cactus/go-statsd-client/statsd"
	"github.com/stretchr/testify/assert"
	"github.com/uber/bark"
)

func TestBarkCactusStatsReporter(t *testing.T) {
	conn, err := net.Listen("tcp", "127.0.0.1:0")
	assert.NoError(t, err)

	defer conn.Close()

	cactusStatter, err := statsd.New(conn.(*net.TCPListener).Addr().String(), "barktest")
	assert.NoError(t, err)

	barkClient := bark.NewStatsReporterFromCactus(cactusStatter)
	barkClient.IncCounter("foo", map[string]string{"tag": "val"}, 1)
	barkClient.RecordTimer("bar", map[string]string{"tag": "val"}, time.Duration(10))
	barkClient.UpdateGauge("baz", map[string]string{"tag": "val"}, 100)
}
