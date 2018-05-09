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

	"github.com/uber-common/bark"

	"go.uber.org/zap"
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
	l.SugaredLogger = l.SugaredLogger.With(key, value) // safe to change because we pass-by-value
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
		zapFields = append(zapFields, zap.Any(k, barkFields[k]))
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
