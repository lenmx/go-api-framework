package xlogger

import (
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"runtime"
	"time"
	"project-name/config"
	"project-name/pkg/constvar"
)

var (
	MonitorLogger *XLogger
	NormalLogger  *XLogger
	ErrorLogger   *XLogger
	DbLogger      *XLogger
)

func InitLogger() {
	MonitorLogger = &XLogger{newZap(typeMonitor)}
	NormalLogger = &XLogger{newZap(typeNormal)}
	ErrorLogger = &XLogger{newZap(typeError)}
	DbLogger = &XLogger{newZap(typeDb)}

	runtime.SetFinalizer(MonitorLogger, finalizer)
	runtime.SetFinalizer(NormalLogger, finalizer)
	runtime.SetFinalizer(ErrorLogger, finalizer)
	runtime.SetFinalizer(DbLogger, finalizer)

	return
}

func (logger *XLogger) Debug(message string, fields ...zap.Field) {
	logger.Logger.Debug(message, fields...)
}

func (logger *XLogger) Debugf(message string, args ...interface{}) {
	logger.Logger.Sugar().Debugf(message, args...)
}

func (logger *XLogger) Info(message string, fields ...zap.Field) {
	logger.Logger.Info(message, fields...)
}

func (logger *XLogger) Infof(message string, args ...interface{}) {
	logger.Logger.Sugar().Infof(message, args...)
}

func (logger *XLogger) Warn(message string, fields ...zap.Field) {
	logger.Logger.Warn(message, fields...)
}

func (logger *XLogger) Warnf(message string, args ...interface{}) {
	logger.Logger.Sugar().Warnf(message, args...)
}

func (logger *XLogger) Error(message string, fields ...zap.Field) {
	logger.Logger.Error(message, fields...)
}

func (logger *XLogger) Errorf(message string, args ...interface{}) {
	logger.Logger.Sugar().Errorf(message, args...)
}

func (logger *XLogger) DPanic(message string, fields ...zap.Field) {
	logger.Logger.DPanic(message, fields...)
}

func (logger *XLogger) DPanicf(message string, args ...interface{}) {
	logger.Logger.Sugar().DPanicf(message, args...)
}

func (logger *XLogger) Panic(message string, fields ...zap.Field) {
	logger.Logger.Panic(message, fields...)
}

func (logger *XLogger) Panicf(message string, args ...interface{}) {
	logger.Logger.Sugar().Panicf(message, args...)
}

func (logger *XLogger) Fatal(message string, fields ...zap.Field) {
	logger.Logger.Fatal(message, fields...)
}

func (logger *XLogger) Fatalf(message string, args ...interface{}) {
	logger.Logger.Sugar().Fatalf(message, args...)
}

func newZap(_type logType) (logger *zap.Logger) {
	core := zapcore.NewCore(
		getEncoder(string(_type)),
		getLogWriter(string(_type)),
		zapcore.DebugLevel)

	logger = zap.New(
		core,
		//zap.AddCaller(),
		//zap.AddCallerSkip(9),
		zap.AddStacktrace(zap.ErrorLevel))

	return
}

func getEncoder(logType string) zapcore.Encoder {
	if logType == "" {
		logType = "normal"
	}
	encoderConfig := zapcore.EncoderConfig{
		MessageKey:    "msg",
		LevelKey:      "level",
		TimeKey:       "ts",
		NameKey:       logType,
		CallerKey:     "caller",
		FunctionKey:   "fn",
		StacktraceKey: "stacktrace",
		LineEnding:    zapcore.DefaultLineEnding,
		EncodeLevel:   zapcore.LowercaseLevelEncoder,
		EncodeTime: func(time time.Time, encoder zapcore.PrimitiveArrayEncoder) {
			encoder.AppendString(time.Format(constvar.DateTimeFormatMs))
		},
		EncodeDuration:   zapcore.SecondsDurationEncoder,
		EncodeCaller:     zapcore.FullCallerEncoder,
		EncodeName:       zapcore.FullNameEncoder,
		ConsoleSeparator: "",
	}

	return zapcore.NewJSONEncoder(encoderConfig)
}

func getLogWriter(logType string) zapcore.WriteSyncer {
	var filename string
	if logType == "" {
		logType = "normal"
	}

	filename = "./log/" + logType + ".log"

	lumberJackLogger := &lumberjack.Logger{
		Filename:   filename,
		MaxSize:    config.G_config.Log.MaxSize,
		MaxBackups: config.G_config.Log.MaxBackups,
		MaxAge:     config.G_config.Log.MaxAge,
		Compress:   config.G_config.Log.Compress,
	}

	return zapcore.NewMultiWriteSyncer(
		zapcore.AddSync(lumberJackLogger),
		zapcore.AddSync(os.Stdout),
	)
}

func (logger *XLogger) close() (err error) {
	err = logger.Logger.Sync()
	logger = nil
	return
}

func finalizer(logger *XLogger) {
	var err error
	if err = logger.close(); err != nil {
		panic(err)
	}
}
