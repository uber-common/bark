package bark_test

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
	"testing"
	"bytes"
	"encoding/json"
	"github.com/uber/bark"
	"github.com/Sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func getLogger() (bark.Logger, *bytes.Buffer) {
	var logrusLogger *logrus.Logger = logrus.New()

	buffer := &bytes.Buffer{}
	logrusLogger.Out = buffer
	logrusLogger.Formatter = new(logrus.JSONFormatter)
	logrusLogger.Level = logrus.DebugLevel

	return bark.NewLoggerFromLogrus(logrusLogger), buffer
}

func checkOutput(t *testing.T, buffer *bytes.Buffer, expectedFields map[string]string) {
	output, _ := buffer.ReadBytes('\n')

	var outputInterface interface{}
	json.Unmarshal(output, &outputInterface)
	outputMap := outputInterface.(map[string]interface{})

	for key, value := range(expectedFields) {
		assert.Equal(t, outputMap[key], value)
	}
}

type testFunc func(l bark.Logger)
type fields map[string]string

func TestBasicLogging(t *testing.T) {
	var testCases = []struct {
		doLogging testFunc
		expectedFields map[string]string
	}{
		{func(l bark.Logger) { l.Info("info", "info") }, fields{"level": "info", "msg": "infoinfo"}},
		{func(l bark.Logger) { l.Infof("infof%s", "infof") }, fields{"level": "info", "msg": "infofinfof"}},
		{func(l bark.Logger) { l.Warn("warn", "warn") }, fields{"level": "warning", "msg": "warnwarn"}},
		{func(l bark.Logger) { l.Warnf("warnf%s", "warnf") }, fields{"level": "warning", "msg": "warnfwarnf"}},
		{func(l bark.Logger) { l.Error("error", "error") }, fields{"level": "error", "msg": "errorerror"}},
		{func(l bark.Logger) { l.Errorf("errorf%s", "errorf") }, fields{"level": "error", "msg": "errorferrorf"}},
		{func(l bark.Logger) { l.WithField("foo", "bar").Info("Info") }, fields{"level": "info", "msg": "Info", "foo":"bar"}},
		{
			func(l bark.Logger) {
				l.WithFields(bark.Fields{"someField":"someValue"}).Info("Yay")
			},
			fields{"level": "info", "msg": "Yay", "someField": "someValue"},
		},
	}

	for _, testCase := range testCases {
		logger, buffer := getLogger()
		testCase.doLogging(logger)
		checkOutput(t, buffer, testCase.expectedFields)
	}

}

func TestPanic(t *testing.T) {
	logger, buffer := getLogger()

	defer func() {
		if r := recover(); r != nil {
			checkOutput(t, buffer, fields{"level": "panic", "msg": "panic"})
		}
	}()

	logger.Panic("panic")
	t.Fail()
}

func TestPanicf(t *testing.T) {
	logger, buffer := getLogger()

	defer func() {
		if r := recover(); r != nil {
			checkOutput(t, buffer, fields{"level": "panic", "msg": "panicpanic"})
		}
	}()

	logger.Panicf("panic%s", "panic")
	t.Fail()
}