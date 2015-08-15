/*
 * Simple tool to make sure we can log
 */
package bark_test

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

	return bark.NewFromLogrus(logrusLogger), buffer
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

func TestInfo(t *testing.T) {
	logger, buffer := getLogger()

	logger.Info("info", "info")
	checkOutput(t, buffer, map[string]string{"level": "info", "msg": "infoinfo"})
}

func TestInfof(t *testing.T) {
	logger, buffer := getLogger()

	logger.Infof("infof%s", "infof")
	checkOutput(t, buffer, map[string]string{"level": "info", "msg": "infofinfof"})
}

func TestWarn(t *testing.T) {
	logger, buffer := getLogger()

	logger.Warn("warn", "warn")
	checkOutput(t, buffer, map[string]string{"level": "warning", "msg": "warnwarn"})
}