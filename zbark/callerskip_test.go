package zbark_test

import (
	"fmt"
	"math/rand"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/uber-common/bark"
	"github.com/uber-common/bark/internal/callerskiphelper"
	"github.com/uber-common/bark/zbark"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"
)

// TestCallerSkip ensures file and line are reported correctly by converted
// loggers.
func TestCallerSkip(t *testing.T) {
	const thisFile = "callerskip_test.go"
	const helperFile = "callerskiphelper.go"

	observableCore, observedLogs := observer.New(zap.DebugLevel)

	// assertCallerFile makes assertions about logs from the observable core.
	assertCallerFile := func(t *testing.T, expectedFile, msg string) {
		logEntries := observedLogs.TakeAll()
		assert.Len(t, logEntries, 1, "expected one log message")
		assert.Contains(t, logEntries[0].Message, msg, "unexpected message")
		callerFile := path.Base(logEntries[0].Caller.File)
		assert.Equal(t, expectedFile, callerFile, "incorrect file")
	}

	testZapper := func(t *testing.T, name string, z *zap.Logger) {
		t.Run(fmt.Sprintf("%s logger same file", name), func(t *testing.T) {
			msg := mkRandomString()
			z.Info(msg)
			assertCallerFile(t, thisFile, msg)
		})
		t.Run(fmt.Sprintf("%s logger using helper", name), func(t *testing.T) {
			msg := mkRandomString()
			callerskiphelper.LogWithZapper(z, msg)
			assertCallerFile(t, helperFile, msg)
		})
	}

	testBarker := func(t *testing.T, name string, b bark.Logger) {
		t.Run(fmt.Sprintf("%s logger same file", name), func(t *testing.T) {
			msg := mkRandomString()
			b.Info(msg)
			assertCallerFile(t, thisFile, msg)
		})
		t.Run(fmt.Sprintf("%s logger using helper", name), func(t *testing.T) {
			msg := mkRandomString()
			callerskiphelper.LogWithBarker(b, msg)
			assertCallerFile(t, helperFile, msg)
		})
	}

	z0 := zap.New(observableCore, zap.AddCaller(), zap.AddStacktrace(zap.DebugLevel))
	b1 := zbark.Barkify(z0)

	// Ensure Barkify doesn't modify its argument.
	testZapper(t, "original zap", z0)
	testBarker(t, "barkified", b1)

	z2 := zbark.Zapify(b1)
	testZapper(t, "re-zapified", z2)

	b3 := zbark.Barkify(z2)
	testBarker(t, "re-barkified", b3)

	z4 := zbark.Zapify(b3)
	testZapper(t, "re-re-zapified", z4)
}

func mkRandomString() string {
	return fmt.Sprintf("test %d", rand.Int())
}
