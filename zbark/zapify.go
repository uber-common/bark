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
	"fmt"

	"github.com/uber-common/bark"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Zapify wraps a bark logger in a compatibility layer to produce a
// *zap.Logger. Note that the wrapper always treats zap's DPanicLevel as an
// error (even in production).
func Zapify(l bark.Logger) *zap.Logger {
	return zap.New(&zapper{l})
}

type zapper struct {
	l bark.Logger
}

func (z *zapper) Enabled(lvl zapcore.Level) bool {
	return true
}

func (z *zapper) With(fs []zapcore.Field) zapcore.Core {
	return &zapper{z.with(fs)}
}

func (z *zapper) Check(ent zapcore.Entry, ce *zapcore.CheckedEntry) *zapcore.CheckedEntry {
	return ce.AddCore(ent, z)
}

func (z *zapper) Write(ent zapcore.Entry, fs []zapcore.Field) error {
	logger := z.with(fs)

	var logFunc func(string, ...interface{})
	switch ent.Level {
	case zapcore.DebugLevel:
		logFunc = logger.Debugf
	case zapcore.InfoLevel:
		logFunc = logger.Infof
	case zapcore.WarnLevel:
		logFunc = logger.Warnf
	case zapcore.ErrorLevel, zapcore.DPanicLevel:
		logFunc = logger.Errorf
	case zapcore.PanicLevel:
		logFunc = logger.Panicf
	case zapcore.FatalLevel:
		logFunc = logger.Fatalf
	default:
		return fmt.Errorf("bark-to-zap compatibility wrapper got unknown level %v", ent.Level)
	}
	// The underlying bark logger timestamps the entry.
	logFunc(ent.Message)
	return nil
}

func (z *zapper) Sync() error {
	return nil
}

func (z *zapper) with(fs []zapcore.Field) bark.Logger {
	me := zapcore.NewMapObjectEncoder()
	for _, f := range fs {
		f.AddTo(me)
	}
	return z.l.WithFields(bark.Fields(me.Fields))
}
