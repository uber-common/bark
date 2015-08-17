package bark

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
	"github.com/Sirupsen/logrus"
)

// Interface provides indirection so Entry and Logger implementations can use exact same methods
type logrusLoggerOrEntry interface {
	Debug(args ...interface{})
	Debugf(format string, args ...interface{})
	Info(args ...interface{})
	Infof(format string, args ...interface{})
	Warn(args ...interface{})
	Warnf(format string, args ...interface{})
	Error(args ...interface{})
	Errorf(format string, args ...interface{})
	Fatal(args ...interface{})
	Fatalf(format string, args ...interface{})
	Panic(args ...interface{})
	Panicf(format string, args ...interface{})
	WithField(key string, value interface{}) *logrus.Entry
	WithFields(keyValues logrus.Fields) *logrus.Entry
}

// The bark-compliant Logger implementation
type barkLogrusLogger struct {
	delegate logrusLoggerOrEntry
}

// The bark-compliant Entry implementation
type barkLogrusEntry struct {
	barkLogrusLogger
	delegate *logrus.Entry

}

// Constructors
func newBarkLogrusLogger(wrappedObject *logrus.Logger) Logger {
	return &barkLogrusLogger{delegate: wrappedObject}
}

func newBarkLogrusEntry(wrappedObject *logrus.Entry) Entry {
	return &barkLogrusEntry{barkLogrusLogger: barkLogrusLogger{delegate: wrappedObject}, delegate: wrappedObject}
}

// All methods thunk to the logrus delegate
func (l *barkLogrusLogger) Debug(args ...interface{}) {
	l.delegate.Debug(args...)
}

func (l *barkLogrusLogger) Debugf(format string, args ...interface{}) {
	l.delegate.Debugf(format, args...)
}

func (l *barkLogrusLogger) Info(args ...interface{}) {
	l.delegate.Info(args...)
}

func (l *barkLogrusLogger) Infof(format string, args ...interface{}) {
	l.delegate.Infof(format, args...)
}

func (l *barkLogrusLogger) Warn(args ...interface{}) {
	l.delegate.Warn(args...)
}

func (l *barkLogrusLogger) Warnf(format string, args ...interface{}) {
	l.delegate.Warnf(format, args...)
}

func (l *barkLogrusLogger) Error(args ...interface{}) {
	l.delegate.Error(args...)
}

func (l *barkLogrusLogger) Errorf(format string, args ...interface{}) {
	l.delegate.Errorf(format, args...)
}

func (l *barkLogrusLogger) Fatal(args ...interface{}) {
	l.delegate.Fatal(args...)
}

func (l *barkLogrusLogger) Fatalf(format string, args ...interface{}) {
	l.delegate.Fatalf(format, args...)
}

func (l *barkLogrusLogger) Panic(args ...interface{}) {
	l.delegate.Panic(args...)
}

func (l *barkLogrusLogger) Panicf(format string, args ...interface{}) {
	l.delegate.Panicf(format, args...)
}

func (e *barkLogrusEntry) Data() map[string]interface{} {
	return e.delegate.Data
}

func (l *barkLogrusLogger) WithField(key string, value interface{}) Entry {
	return newBarkLogrusEntry(l.delegate.WithField(key, value))
}

func (l *barkLogrusLogger) WithFields(keyValues LogFields) Entry {
	return newBarkLogrusEntry(l.delegate.WithFields(logrus.Fields(keyValues.Fields())))
}

