package xlogger

import (
	"go.uber.org/zap"
)

type XLogger struct {
	*zap.Logger
}
