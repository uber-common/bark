/*
 * "bark" provides an abstraction for loggers used in Uber's
 * Go libraries.  It decouples these libraries slightly from specific
 * logger implementations; for example, the popular open source library
 * "logrus," which offers no interfaces (and thus cannot be, for instance, easily mocked).
 * Users may choose to implement the interface themselves or use the provided logrus
 * wrapper.
 */
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
	"time"
	"github.com/cactus/go-statsd-client/statsd"
)

type Logger interface {
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
	WithField(key string, value interface{}) Entry
	WithFields(keyValues LogFields) Entry
}

type Entry interface {
	Logger
	Data() map[string]interface{}
}

type LogFields interface {
	Fields() map[string]interface{}
}

type Fields map[string]interface{}

func (f Fields) Fields() map[string]interface{} {
	return f
}

// Create a bark-compliant wrapper for a logrus-brand logger
func NewLoggerFromLogrus(wrappedLogger *logrus.Logger) Logger {
	return newBarkLogrusLogger(wrappedLogger)
}

type StatsReporter interface {
	IncCounter(name string, tags map[string]string, value int64)
	UpdateGauge(name string, tags map[string]string, value int64)
	RecordTimer(name string, tags map[string]string, d time.Duration)
}

func NewStatsReporterFromCactus(wrappedStatsd statsd.Statter) StatsReporter {
	return newBarkCactusStatsReporter(wrappedStatsd)
}