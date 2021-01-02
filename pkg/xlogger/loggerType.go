package xlogger

type logType string

const (
	typeMonitor logType = "monitor"
	typeNormal  logType = "normal"
	typeError   logType = "error"
	typeDb      logType = "db"
)
