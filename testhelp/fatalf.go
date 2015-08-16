package main

import (
	"github.com/Sirupsen/logrus"
	"github.com/uber/bark"
	"os"
)

func main() {
	var logrusLogger *logrus.Logger = logrus.New()
	logrusLogger.Formatter = new(logrus.JSONFormatter)
	logrusLogger.Level = logrus.DebugLevel
	bark.NewFromLogrus(logrusLogger).Fatalf("halp%s", "pls")
	os.Exit(2)
}
