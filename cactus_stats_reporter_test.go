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

package bark_test

import (
	"testing"
	"time"

	"github.com/uber-common/bark"
	"github.com/uber-common/bark/mocks"
)

func TestBarkCactusStatsReporter(t *testing.T) {

	statter := &mocks.Statter{}
	barkClient := bark.NewStatsReporterFromCactus(statter)

	statter.On("Inc", "foo", int64(7), float32(1.0)).Return(nil)
	barkClient.IncCounter("foo", bark.Tags{"tag": "val"}, 7)

	statter.On("TimingDuration", "bar", time.Duration(10), float32(1.0)).Return(nil)
	barkClient.RecordTimer("bar", bark.Tags{"tag": "val"}, time.Duration(10))

	statter.On("Gauge", "baz", int64(123), float32(1.0)).Return(nil)
	barkClient.UpdateGauge("baz", bark.Tags{"tag": "val"}, 123)

	statter.AssertCalled(t, "Inc", "foo", int64(7), float32(1.0))
	statter.AssertCalled(t, "TimingDuration", "bar", time.Duration(10), float32(1.0))
	statter.AssertCalled(t, "Gauge", "baz", int64(123), float32(1.0))
}
