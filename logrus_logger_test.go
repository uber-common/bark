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

package bark_test

import (
	"bytes"
	"encoding/json"
	"os"
	"os/exec"
	"testing"

	"github.com/Sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/uber-common/bark"
)

// Create a logrus logger that writes its out output to a buffer for inspection
func getLogrusLogger() (*logrus.Logger, *bytes.Buffer) {
	var logrusLogger *logrus.Logger = logrus.New()

	buffer := &bytes.Buffer{}
	logrusLogger.Out = buffer
	logrusLogger.Formatter = new(logrus.JSONFormatter)
	logrusLogger.Level = logrus.DebugLevel

	return logrusLogger, buffer
}

// Create a bark wrapper for a logrus logger backed by a buffer
func getBarkLogger() (bark.Logger, *bytes.Buffer) {
	logrusLogger, buffer := getLogrusLogger()
	return bark.NewLoggerFromLogrus(logrusLogger), buffer
}

// Extract map of keys and values from raw json data in buffer
func parseLogBytes(logBytes []byte) map[string]interface{} {
	var unmarshalledData interface{}
	json.Unmarshal(logBytes, &unmarshalledData)

	if unmarshalledData == nil {
		return nil
	}

	return unmarshalledData.(map[string]interface{})
}

// Validate bark output against logrus output
func validateOutput(t *testing.T, barkBytes []byte, logrusBytes []byte) {
	barkMap := parseLogBytes(barkBytes)
	require.NotNil(t, barkMap, "Should be able to parse bark output as JSON")

	logrusMap := parseLogBytes(logrusBytes)
	require.NotNil(t, logrusMap, "Should be able to parse logrus output as JSON")

	// Make sure we're checking at least the fields we expect
	minFields := []string{"time", "level", "msg"}
	for _, key := range minFields {
		_, ok := logrusMap[key]
		require.True(t, ok, "Logrus missing required field: %s", key)
	}

	// Make sure bark output has everything logrus does
	for key, logrusValue := range logrusMap {
		barkValue, ok := barkMap[key]
		require.True(t, ok)

		// Can't mock time to logrus, so have to skip it
		if key != "time" {
			require.Equal(t, logrusValue, barkValue, "Field of outputs of output should match logrus")
		}
	}
}

func TestInfo(t *testing.T) {
	logAndValidate(t, func(barkLogger bark.Logger, logrusLogger *logrus.Logger) {
		barkLogger.Info("info1", "info2")
		logrusLogger.Info("info1", "info2")
	})
}

func TestInfof(t *testing.T) {
	logAndValidate(t, func(barkLogger bark.Logger, logrusLogger *logrus.Logger) {
		barkLogger.Infof("infof1%s", "infof2")
		logrusLogger.Infof("infof1%s", "infof2")
	})
}

func TestError(t *testing.T) {
	logAndValidate(t, func(barkLogger bark.Logger, logrusLogger *logrus.Logger) {
		barkLogger.Error("error1", "error2")
		logrusLogger.Error("error1", "error2")
	})
}

func TestErrorf(t *testing.T) {
	logAndValidate(t, func(barkLogger bark.Logger, logrusLogger *logrus.Logger) {
		barkLogger.Errorf("errorf1%s", "errorf2")
		logrusLogger.Errorf("errorf1%s", "errorf2")
	})
}

func TestWarn(t *testing.T) {
	logAndValidate(t, func(barkLogger bark.Logger, logrusLogger *logrus.Logger) {
		barkLogger.Warn("warn1", "warn2")
		logrusLogger.Warn("warn1", "warn2")
	})
}

func TestWarnf(t *testing.T) {
	logAndValidate(t, func(barkLogger bark.Logger, logrusLogger *logrus.Logger) {
		barkLogger.Warnf("warnf1%s", "warnf2")
		logrusLogger.Warnf("warnf1%s", "warnf2")
	})
}

func TestWithField(t *testing.T) {
	logAndValidate(t, func(barkLogger bark.Logger, logrusLogger *logrus.Logger) {
		barkLogger.WithField("key", "value").Info("withfield")
		logrusLogger.WithField("key", "value").Info("withfield")
	})
}

func TestWithFields(t *testing.T) {
	logAndValidate(t, func(barkLogger bark.Logger, logrusLogger *logrus.Logger) {
		barkLogger.WithFields(bark.Fields{"k1": "v1", "k2": "v2"}).Info("withfields")
		logrusLogger.WithFields(logrus.Fields{"k1": "v1", "k2": "v2"}).Info("withfields")
	})
}

func TestGetFields(t *testing.T) {
	var logger bark.Logger

	// Plain logger
	logger = bark.NewLoggerFromLogrus(logrus.New())
	require.Equal(t, logger.Fields(), bark.Fields(nil))

	// Add nil, don't crash
	logger = bark.NewLoggerFromLogrus(logrus.New())
	logger = logger.WithFields(nil)
	require.Equal(t, logger.Fields(), bark.Fields(nil))

	// One field added
	logger = bark.NewLoggerFromLogrus(logrus.New())
	logger = logger.WithField("foo", "bar")
	require.Equal(t, logger.Fields(), bark.Fields{"foo": "bar"})

	// Two fields added at once
	logger = bark.NewLoggerFromLogrus(logrus.New())
	logger = logger.WithFields(bark.Fields{"foo": "bar", "baz": "bump"})
	require.Equal(t, logger.Fields(), bark.Fields{"foo": "bar", "baz": "bump"})

	// One then one
	logger = bark.NewLoggerFromLogrus(logrus.New())
	logger = logger.WithField("foo", "bar")
	logger = logger.WithField("baz", "bump")
	require.Equal(t, logger.Fields(), bark.Fields{"foo": "bar", "baz": "bump"})

	// Two then one
	logger = bark.NewLoggerFromLogrus(logrus.New())
	logger = logger.WithFields(bark.Fields{"foo": "bar", "baz": "bump"})
	logger = logger.WithField("x", "y")
	require.Equal(t, logger.Fields(), bark.Fields{"foo": "bar", "baz": "bump", "x": "y"})

	// One then two
	logger = bark.NewLoggerFromLogrus(logrus.New())
	logger = logger.WithField("x", "y")
	logger = logger.WithFields(bark.Fields{"foo": "bar", "baz": "bump"})
	require.Equal(t, logger.Fields(), bark.Fields{"foo": "bar", "baz": "bump", "x": "y"})

	// Two then two
	logger = bark.NewLoggerFromLogrus(logrus.New())
	logger = logger.WithFields(bark.Fields{"foo": "bar", "baz": "bump"})
	logger = logger.WithFields(bark.Fields{"a": "b", "c": "d"})
	require.Equal(t, logger.Fields(), bark.Fields{"foo": "bar", "baz": "bump", "a": "b", "c": "d"})

	// Add empty map
	logger = bark.NewLoggerFromLogrus(logrus.New())
	logger = logger.WithFields(bark.Fields{})
	require.Equal(t, logger.Fields(), bark.Fields{})
}

func doPanic(t *testing.T, panicker func(...interface{})) {
	defer func() {
		if r := recover(); r != nil {
			return
		}
	}()

	panicker("panic", "panic")
	t.Error("Expected to panic but execution did not stop")
}

func TestPanic(t *testing.T) {
	logAndValidate(t, func(barkLogger bark.Logger, logrusLogger *logrus.Logger) {
		doPanic(t, barkLogger.Panic)
		doPanic(t, logrusLogger.Panic)
	})
}

func doPanicf(t *testing.T, panicf func(string, ...interface{})) {
	defer func() {
		if r := recover(); r != nil {
			return
		}
	}()

	panicf("panicf%s", "panicf")
	t.Error("Expected to panic but execution did not stop")
}

func TestPanicf(t *testing.T) {
	logAndValidate(t, func(barkLogger bark.Logger, logrusLogger *logrus.Logger) {
		doPanicf(t, barkLogger.Panicf)
		doPanicf(t, logrusLogger.Panicf)
	})
}

// Main test driver: create loggers backed by buffers, drive operations on both, compare results
func logAndValidate(t *testing.T, driver func(barkLogger bark.Logger, logrusLogger *logrus.Logger)) {
	barkLogger, barkBuffer := getBarkLogger()
	logrusLogger, logrusBuffer := getLogrusLogger()

	driver(barkLogger, logrusLogger)
	validateOutput(t, barkBuffer.Bytes(), logrusBuffer.Bytes())
}

// Fatal tests are special: have to execute a command and confirm that it exits
func execFatalTool(t *testing.T, how string) []byte {
	helperBinary := "./testhelp/fatal"

	_, err := os.Stat(helperBinary)
	require.NoError(t, err, "Helper binary should exist: run 'make test'!")

	cmd := exec.Command(helperBinary, how)
	output, err := cmd.CombinedOutput()
	if assert.Error(t, err, "Process should exit with an error") {
		_, ok := err.(*exec.ExitError)
		assert.True(t, ok, "Process should exit with an ExitError")
	}

	require.Contains(t, string(output), "fatal error")
	return output
}

func TestFatal(t *testing.T) {
	logrusStderr := execFatalTool(t, "logrus.Fatal")
	barkStderr := execFatalTool(t, "bark.Fatal")
	validateOutput(t, barkStderr, logrusStderr)
}

func TestFatalf(t *testing.T) {
	logrusStderr := execFatalTool(t, "logrus.Fatalf")
	barkStderr := execFatalTool(t, "bark.Fatalf")
	validateOutput(t, barkStderr, logrusStderr)
}
