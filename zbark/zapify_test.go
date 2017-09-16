// Copyright (c) 2017 Uber Technologies, Inc.
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

package zbark_test

import (
	"bytes"
	"encoding/json"
	"testing"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/uber-common/bark"
	"github.com/uber-common/bark/zbark"
)

func newTestZapper() (*zap.Logger, *bytes.Buffer) {
	buf := bytes.NewBuffer(nil)
	barkLogger := &logrus.Logger{
		Out:   buf,
		Hooks: make(logrus.LevelHooks),
		Formatter: &logrus.JSONFormatter{
			DisableTimestamp: true,
		},
		Level: logrus.DebugLevel,
	}
	return zbark.Zapify(bark.NewLoggerFromLogrus(barkLogger)), buf
}

func assertJSON(t testing.TB, expected map[string]interface{}, buf *bytes.Buffer) {
	line := bytes.TrimSpace(buf.Bytes())
	msg := make(map[string]interface{})
	if err := json.Unmarshal(line, &msg); err != nil {
		t.Fatalf("can't unmarshal JSON log: %s", string(line))
	}
	assert.Equal(t, expected, msg, "unexpected log message")
}

func TestZapLoggerWith(t *testing.T) {
	log, buf := newTestZapper()
	log.With(zap.String("foo", "bar")).Info("hello", zap.String("baz", "quux"))
	assertJSON(t, map[string]interface{}{
		"foo":   "bar",
		"baz":   "quux",
		"msg":   "hello",
		"level": "info",
	}, buf)
}

func TestZapLoggerLogging(t *testing.T) {
	const msg = "hello"

	log, buf := newTestZapper()
	assertLogged := func(expected string) {
		assertJSON(t, map[string]interface{}{
			"msg":   msg,
			"level": expected,
		}, buf)
	}

	tests := []struct {
		f     func(string, ...zapcore.Field)
		level zapcore.Level
		want  string
	}{
		{log.Debug, zapcore.DebugLevel, "debug"},
		{log.Info, zapcore.InfoLevel, "info"},
		{log.Warn, zapcore.WarnLevel, "warning"},
		{log.Error, zapcore.ErrorLevel, "error"},
		{log.DPanic, zapcore.DPanicLevel, "error"},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			tt.f(msg)
			assertLogged(tt.want)
			buf.Reset()

			if ce := log.Check(tt.level, msg); ce != nil {
				ce.Write()
			}
			assertLogged(tt.want)
			buf.Reset()
		})
	}
}
