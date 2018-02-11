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
	"errors"
	"path"
	"testing"

	"github.com/uber-common/bark"
	"github.com/uber-common/bark/zbark"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
)

func newTestBarker() (bark.Logger, *observer.ObservedLogs) {
	core, logs := observer.New(zap.LevelEnablerFunc(func(zapcore.Level) bool {
		return true
	}))
	return zbark.Barkify(zap.New(core, zap.AddCaller())), logs
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

		callerFile := path.Base(entry.Caller.File)
		assert.Equal(t, "barkify_test.go", callerFile, "caller was not test file")
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
