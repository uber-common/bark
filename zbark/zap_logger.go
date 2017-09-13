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

// New builds a Bark logger from a Zap logger.
func New(l *zap.Logger) bark.Logger {
	return barkZapLogger{l.Sugar()}
}

type barkZapLogger struct{ *zap.SugaredLogger }

func (l barkZapLogger) WithField(key string, value interface{}) bark.Logger {
	l.SugaredLogger = l.SugaredLogger.With(key, value) // safe to change because we pass-by-value
	return l
}

func (l barkZapLogger) WithFields(keyValues bark.LogFields) bark.Logger {
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

func (l barkZapLogger) WithError(err error) bark.Logger {
	l.SugaredLogger = l.SugaredLogger.With(zap.Error(err)) // safe to change because we pass-by-value
	return l
}

func (l barkZapLogger) Fields() bark.Fields {
	l.SugaredLogger.Warn("Fields() call to bark logger is not supported by Zap")
	return nil
}
