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
	"errors"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/uber-common/bark"
	"github.com/uber-common/bark/zbark"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
)

func newTestBarker() (bark.Logger, *observer.ObservedLogs) {
	core, logs := observer.New(zap.LevelEnablerFunc(func(zapcore.Level) bool {
		return true
	}))
	return zbark.Barkify(zap.New(core)), logs
}

func TestBarkLoggerWith(t *testing.T) {
	t.Run("WithField", func(t *testing.T) {
		l, logs := newTestBarker()
		l = l.WithField("foo", "bar")

		l.Info("hello")
		require.Equal(t, 1, logs.Len(), "message count did not match")

		entry := logs.All()[0]
		require.Len(t, entry.Context, 1, "context size did not match")
		field := entry.Context[0]

		assert.Equal(t, zapcore.StringType, field.Type, "field type did not match")
		assert.Equal(t, "foo", field.Key, "field name did not match")
		assert.Equal(t, "bar", field.String, "field value did not match")
	})

	t.Run("WithFields", func(t *testing.T) {
		l, logs := newTestBarker()
		l = l.WithFields(bark.Fields{
			"foo": "bar",
			"baz": int8(42),
			"qux": errors.New("great sadness"),
		})

		l.Info("hello")
		require.Equal(t, 1, logs.Len(), "message count did not match")

		entry := logs.All()[0]
		require.Len(t, entry.Context, 3, "context size did not match")

		// We ensure deterministic ordering of fields by sorting them.
		foo, baz, qux := entry.Context[1], entry.Context[0], entry.Context[2]

		assert.Equal(t, "foo", foo.Key, "field name did not match")
		assert.Equal(t, zapcore.StringType, foo.Type, "field type did not match")
		assert.Equal(t, "bar", foo.String, "field value did not match")

		assert.Equal(t, "baz", baz.Key, "field name did not match")
		assert.Equal(t, zapcore.Int8Type, baz.Type, "field type did not match")
		assert.EqualValues(t, 42, baz.Integer, "field value did not match")

		assert.Equal(t, "qux", qux.Key, "field name did not match")
		assert.Equal(t, zapcore.ErrorType, qux.Type, "field type did not match")
		assert.Equal(t, errors.New("great sadness"), qux.Interface, "field value did not match")
	})

	t.Run("WithError", func(t *testing.T) {
		l, logs := newTestBarker()
		l = l.WithError(errors.New("great sadness"))

		l.Info("hello")
		require.Equal(t, 1, logs.Len(), "message count did not match")

		entry := logs.All()[0]
		require.Len(t, entry.Context, 1, "context size did not match")
		field := entry.Context[0]

		assert.Equal(t, zapcore.ErrorType, field.Type, "field type did not match")
		assert.Equal(t, "error", field.Key, "field name did not match")
		assert.Equal(t, errors.New("great sadness"), field.Interface, "field value did not match")
	})

	t.Run("Fields", func(t *testing.T) {
		l, logs := newTestBarker()
		l = l.WithField("foo", "bar").WithError(errors.New("great sadness"))
		assert.Empty(t, l.Fields(), "expected empty field list")

		require.Equal(t, 1, logs.Len(), "message count did not match")

		entry := logs.All()[0]
		assert.Equal(t,
			"zap-to-bark compatibility wrapper does not support Fields method", entry.Message,
			"message did not match")
		assert.Equal(t, zapcore.WarnLevel, entry.Level, "message level did not match")
	})
}

func TestZapLogrusParity(t *testing.T) {
	var (
		boolv     bool          = true
		float64v  float64       = 1.23
		float32v  float32       = 4.5
		intv      int           = 42
		int8v     int8          = 6
		int16v    int16         = 55
		int32v    int32         = 200
		int64v    int64         = 1234
		stringv   string        = "hello world"
		uintv     uint          = 128
		uint8v    uint8         = 7
		uint16v   uint16        = 56
		uint32v   uint32        = 201
		uint64v   uint64        = 1235
		timev     time.Time     = time.Now()
		durationv time.Duration = 42 * time.Minute
	)

	tests := []struct {
		desc  string
		field interface{}
	}{
		{desc: "nil", field: nil},
		{desc: "bool", field: boolv},
		{desc: "bool/ptr", field: &boolv},
		{desc: "bool/slice", field: []bool{true, false, true}},
		{desc: "float32", field: float32v},
		{desc: "float32/ptr", field: &float32v},
		{desc: "float32/slice", field: []float32{1.2, 2.3, 3.4}},
		{desc: "float64", field: float64v},
		{desc: "float64/ptr", field: &float64v},
		{desc: "float64/slice", field: []float64{1.2, 2.3, 3.4}},
		{desc: "int", field: intv},
		{desc: "int/ptr", field: &intv},
		{desc: "int/slice", field: []int{1, 2, 3}},
		{desc: "int8", field: int8v},
		{desc: "int8/ptr", field: &int8v},
		{desc: "int8/slice", field: []int8{1, 2, 3}},
		{desc: "int16", field: int16v},
		{desc: "int16/ptr", field: &int16v},
		{desc: "int16/slice", field: []int16{4, 5, 6}},
		{desc: "int32", field: int32v},
		{desc: "int32/ptr", field: &int32v},
		{desc: "int32/slice", field: []int32{7, 8, 9}},
		{desc: "int64", field: int64v},
		{desc: "int64/ptr", field: &int64v},
		{desc: "int64/slice", field: []int64{10, 11, 12}},
		{desc: "string", field: stringv},
		{desc: "string/ptr", field: &stringv},
		{desc: "string/slice", field: []string{"a", "b", "c"}},
		{desc: "uint", field: uintv},
		{desc: "uint/ptr", field: &uintv},
		{desc: "uint/slice", field: []uint{1, 2, 3}},
		{desc: "uint8", field: uint8v},
		{desc: "uint8/ptr", field: &uint8v},
		{desc: "uint16", field: uint16v},
		{desc: "uint16/ptr", field: &uint16v},
		{desc: "uint16/slice", field: []uint16{4, 5, 6}},
		{desc: "uint32", field: uint32v},
		{desc: "uint32/ptr", field: &uint32v},
		{desc: "uint32/slice", field: []uint32{7, 8, 9}},
		{desc: "uint64", field: uint64v},
		{desc: "uint64/ptr", field: &uint64v},
		{desc: "uint64/slice", field: []uint64{10, 11, 12}},
		{desc: "byte/slice", field: []byte{1, 2, 3}},
		{desc: "time", field: timev},
		{desc: "time/ptr", field: &timev},
		{desc: "time/slice", field: []time.Time{timev, timev.Add(time.Second), timev.Add(time.Hour)}},
		{desc: "duration", field: durationv},
		{desc: "duration/ptr", field: &durationv},
		{desc: "duration/slice", field: []time.Duration{time.Second, time.Minute, time.Hour}},
		{desc: "error", field: errors.New("great sadness")},
		{desc: "slice/any", field: []interface{}{"a", 1, true}},
		{desc: "slice/empty", field: emptyZapArray{}},
		{desc: "struct", field: struct{ SomeField string }{"foo"}},
		{desc: "struct/empty", field: emptyZapStruct{}},
		{desc: "struct/stringer", field: stringerStruct{SomeField: "foo"}},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			gotLogrus := logrusLog(t, tt.field)
			t.Logf("logrus = %#v", gotLogrus)

			gotZap := zapLog(t, tt.field)
			t.Logf("zap = %#v", gotZap)

			assert.Equal(t, gotLogrus, gotZap,
				"output from logrus and zap does not match")
		})
	}
}

func logrusLog(t *testing.T, field interface{}) interface{} {
	var buff bytes.Buffer
	logger := logrus.New()
	logger.SetOutput(&buff)
	logger.SetFormatter(&logrus.JSONFormatter{})

	got, err := logAndParseField(field, &buff, bark.NewLoggerFromLogrus(logger))
	require.NoError(t, err, "unable to parse logrus output: %s", buff.String())
	return got
}

func zapLog(t *testing.T, field interface{}) interface{} {
	encoderCfg := zap.NewDevelopmentEncoderConfig()
	encoderCfg.EncodeTime = zapcore.RFC3339NanoTimeEncoder

	var buff bytes.Buffer
	logger := zap.New(zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderCfg),
		zapcore.AddSync(&buff),
		zapcore.InfoLevel,
	))

	got, err := logAndParseField(field, &buff, zbark.Barkify(logger))
	require.NoError(t, err, "unable to parse Zap output: %s", buff.String())
	return got
}

// Logs the provided field to the provided Bark logger and returns the
// JSON object that results from that field.
//
// The logger must be configured to log JSON to the buffer.
func logAndParseField(field interface{}, buff *bytes.Buffer, logger bark.Logger) (interface{}, error) {
	logger.WithField("field", field).Info("msg")
	var msg struct {
		Field interface{} `json:"field"`
	}
	err := json.Unmarshal(buff.Bytes(), &msg)
	return msg.Field, err
}

type stringerStruct struct{ SomeField string }

func (s stringerStruct) String() string {
	return "stringerStruct.String()"
}

type emptyZapStruct struct{}

func (emptyZapStruct) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	return nil
}

type emptyZapArray struct{}

func (emptyZapStruct) MarshalLogArray(enc zapcore.ArrayEncoder) error {
	return nil
}
