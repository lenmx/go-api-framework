package xxorm

import (
	"context"
	"encoding/json"
	"go.uber.org/zap"
	"time"
	"project-name/config"
	"project-name/pkg/xlogger"
	"xorm.io/builder"
	"xorm.io/xorm/contexts"
)

type LogHook struct {
	*xlogger.XLogger
	before func(c *contexts.ContextHook) (context.Context, error)
	after  func(c *contexts.ContextHook) error
}

var _ contexts.Hook = &LogHook{}

func (h *LogHook) BeforeProcess(c *contexts.ContextHook) (context.Context, error) {
	return h.before(c)
}

func (h *LogHook) AfterProcess(c *contexts.ContextHook) error {
	return h.after(c)
}

func NewLogHook() *LogHook {
	hook := &LogHook{
		before: before,
		after:  after,
	}

	return hook
}

func before(c *contexts.ContextHook) (context.Context, error) {
	return c.Ctx, nil
}

func after(c *contexts.ContextHook) error {
	var (
		sql   string
		_args []byte
		args  string
		err   error
	)

	ts := string(c.ExecuteTime.Milliseconds()) + "ms"

	sql, err = builder.ConvertToBoundSQL(c.SQL, c.Args)
	if err != nil {
		sql = c.SQL
	}

	_args, err = json.Marshal(c.Args)
	if err != nil {
		args = ""
	} else {
		args = string(_args)
	}

	// 异常
	if c.Err != nil {
		xlogger.DbLogger.Error("",
			zap.NamedError("err", c.Err),
			zap.String("SQL", sql),
			zap.String("args", args),
			zap.String("exec_time", ts),
		)

		return nil
	}

	// 慢日志
	slowLogTimeDuration, _ := time.ParseDuration(config.G_config.Db.SlowLogTime)
	if err == nil && slowLogTimeDuration > 0 {
		if c.ExecuteTime >= slowLogTimeDuration {
			xlogger.DbLogger.Warn("",
				zap.String("SQL", sql),
				zap.String("args", args),
				zap.String("exec_time", ts),
			)
			return nil
		}
	}

	// 普通日志
	xlogger.DbLogger.Info("",
		zap.String("SQL", sql),
		zap.String("args", args),
		zap.String("exec_time", ts),
	)
	return nil
}
