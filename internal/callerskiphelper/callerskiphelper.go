// Provides an extra level in stack trace to accompany logs.
package callerskiphelper

import (
	"fmt"

	"github.com/uber-common/bark"
	"go.uber.org/zap"
)

// LogWithBarker emits a log entry using a bark logger.
func LogWithBarker(b bark.Logger, msg string) {
	b.Error(fmt.Sprintf("logged from helper: %s", msg))
}

// LogWithZapper emits a log entry using a zap logger.
func LogWithZapper(z *zap.Logger, msg string) {
	z.Error(fmt.Sprintf("logged from helper: %s", msg))
}
