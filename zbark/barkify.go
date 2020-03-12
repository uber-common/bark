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

package zbark

import (
	"sort"
	"time"

	"github.com/uber-common/bark"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Barkify wraps a zap logger in a compatibility layer so that it satisfies
// the bark.Logger interface. Note that the wrapper always returns nil from
// the Fields method, since zap doesn't support this functionality.
func Barkify(l *zap.Logger) bark.Logger {
	if z, ok := l.Core().(*zapper); ok {
		return z.l
	}
	return barker{l.WithOptions(zap.AddCallerSkip(_barkifyCallerSkip)).Sugar()}
}

type barker struct{ *zap.SugaredLogger }

func (l barker) WithField(key string, value interface{}) bark.Logger {
	l.SugaredLogger = l.SugaredLogger.With(toZapField(key, value)) // safe to change because we pass-by-value
	return l
}

func (l barker) WithFields(keyValues bark.LogFields) bark.Logger {
	barkFields := keyValues.Fields()

	// Deterministic ordering of fields.
	keys := make([]string, 0, len(barkFields))
	for k := range barkFields {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	zapFields := make([]interface{}, 0, len(barkFields))
	for _, k := range keys {
		zapFields = append(zapFields, toZapField(k, barkFields[k]))
	}

	l.SugaredLogger = l.SugaredLogger.With(zapFields...) // safe to change because we pass-by-value
	return l
}

func (l barker) WithError(err error) bark.Logger {
	l.SugaredLogger = l.SugaredLogger.With(zap.Error(err)) // safe to change because we pass-by-value
	return l
}

func (l barker) Fields() bark.Fields {
	// Zap has already marshaled the accumulated logger context to []byte, so we
	// can't reconstruct the original objects. To satisfy this interface, just
	// return nil.
	l.SugaredLogger.Warn("zap-to-bark compatibility wrapper does not support Fields method")
	return nil
}

// toZapField converts a logrus field to a zap field.
//
// This relies on Zap's field constructors for fields when possible but falls
// back to zap.Reflect for all other fields.
//
// In the interest of conforming to the output of logrus, this does not use
// zap.Any. That's because logrus eventually uses encoding/json to encode
// everything -- same as zap.Reflect.
func toZapField(key string, v interface{}) zap.Field {
	switch val := v.(type) {

	// Types that implement ObjectMarshaler or ArrayMarshaler are
	// explicitly specifying how they want to be logged with Zap. We don't
	// care about parity with logrus in that case.
	case zapcore.ObjectMarshaler:
		return zap.Object(key, val)
	case zapcore.ArrayMarshaler:
		return zap.Array(key, val)

	case bool:
		return zap.Bool(key, val)
	case *bool:
		return zap.Boolp(key, val)
	case []bool:
		return zap.Bools(key, val)

	case float32:
		return zap.Float32(key, val)
	case *float32:
		return zap.Float32p(key, val)
	case []float32:
		return zap.Float32s(key, val)

	case float64:
		return zap.Float64(key, val)
	case *float64:
		return zap.Float64p(key, val)
	case []float64:
		return zap.Float64s(key, val)

	case int:
		return zap.Int(key, val)
	case *int:
		return zap.Intp(key, val)
	case []int:
		return zap.Ints(key, val)

	case int8:
		return zap.Int8(key, val)
	case *int8:
		return zap.Int8p(key, val)
	case []int8:
		return zap.Int8s(key, val)

	case int16:
		return zap.Int16(key, val)
	case *int16:
		return zap.Int16p(key, val)
	case []int16:
		return zap.Int16s(key, val)

	case int32:
		return zap.Int32(key, val)
	case *int32:
		return zap.Int32p(key, val)
	case []int32:
		return zap.Int32s(key, val)

	case int64:
		return zap.Int64(key, val)
	case *int64:
		return zap.Int64p(key, val)
	case []int64:
		return zap.Int64s(key, val)

	case string:
		return zap.String(key, val)
	case *string:
		return zap.Stringp(key, val)
	case []string:
		return zap.Strings(key, val)

	case uint:
		return zap.Uint(key, val)
	case *uint:
		return zap.Uintp(key, val)
	case []uint:
		return zap.Uints(key, val)

	case uint8:
		return zap.Uint8(key, val)
	case *uint8:
		return zap.Uint8p(key, val)
	// []uint8 == []byte. See below.

	case uint16:
		return zap.Uint16(key, val)
	case *uint16:
		return zap.Uint16p(key, val)
	case []uint16:
		return zap.Uint16s(key, val)

	case uint32:
		return zap.Uint32(key, val)
	case *uint32:
		return zap.Uint32p(key, val)
	case []uint32:
		return zap.Uint32s(key, val)

	case uint64:
		return zap.Uint64(key, val)
	case *uint64:
		return zap.Uint64p(key, val)
	case []uint64:
		return zap.Uint64s(key, val)

	case []byte:
		return zap.Binary(key, val)

	case time.Time:
		return zap.Time(key, val)
	case *time.Time:
		return zap.Timep(key, val)
	case []time.Time:
		// return zap.Times(key, val)
		// https://github.com/uber-go/zap/issues/798
		// Fall back to zap.Reflect meanwhile.

	// Logrus logs time.Duration as numbers so we should do the
	// same.
	case time.Duration:
		return zap.Int64(key, int64(val))
	case *time.Duration:
		return zap.Int64p(key, (*int64)(val))
	case []time.Duration:
		return zap.Array(key, durationAsIntSlice(val))

	case error:
		return zap.NamedError(key, val)
	}

	// Use zap.Reflect for everything else.
	return zap.Reflect(key, v)
}

type durationAsIntSlice []time.Duration

func (ds durationAsIntSlice) MarshalLogArray(enc zapcore.ArrayEncoder) error {
	for _, d := range ds {
		enc.AppendInt64(int64(d))
	}
	return nil
}
