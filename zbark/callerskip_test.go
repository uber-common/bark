package zbark_test

import (
	"fmt"
	"math/rand"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/uber-common/bark/zbark"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"
)

// TestCallerSkip ensures file and line are reported correctly by converted
// loggers.
func TestCallerSkip(t *testing.T) {
	const expectedFile = "callerskip_test.go"

	observableCore, observedLogs := observer.New(zap.DebugLevel)

	// assertCallerFile makes assertions about logs from the observable core.
	assertCallerFile := func(t *testing.T, msg string) {
		logEntries := observedLogs.TakeAll()
		assert.Len(t, logEntries, 1, "expected one log message")
		assert.Equal(t, msg, logEntries[0].Message, "unexpected message")
		callerFile := path.Base(logEntries[0].Caller.File)
		assert.Equal(t, expectedFile, callerFile, "incorrect file")
	}

	z0 := zap.New(observableCore, zap.AddCaller(), zap.AddStacktrace(zap.DebugLevel))
	b1 := zbark.Barkify(z0)

	t.Run("original zap logger", func(t *testing.T) {
		// Ensure Barkify doesn't modify its argument.
		msg := mkRandomString()
		z0.Info(msg)
		assertCallerFile(t, msg)
	})

	t.Run("barkified logger", func(t *testing.T) {
		msg := mkRandomString()
		b1.Info(msg)
		assertCallerFile(t, msg)
	})

	z2 := zbark.Zapify(b1)

	t.Run("re-zapified logger", func(t *testing.T) {
		msg := mkRandomString()
		z2.Info(msg)
		assertCallerFile(t, msg)
	})

	b3 := zbark.Barkify(z2)

	t.Run("re-barkified logger", func(t *testing.T) {
		msg := mkRandomString()
		b3.Info(msg)
		assertCallerFile(t, msg)
	})

	z4 := zbark.Zapify(b3)

	t.Run("re-re-zapified logger", func(t *testing.T) {
		msg := mkRandomString()
		z4.Info(msg)
		assertCallerFile(t, msg)
	})
}

func mkRandomString() string {
	return fmt.Sprintf("test %d", rand.Int())
}
