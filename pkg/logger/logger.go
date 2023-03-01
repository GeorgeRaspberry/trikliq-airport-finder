package logger

import (
	"trikliq-airport-finder/pkg/logger/encoders"
	"trikliq-airport-finder/pkg/logger/filters"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	Log *zap.Logger
)

func init() {
	encoders.Setup()
	Log, _ = Setup()

	defer Log.Sync()
}

// configure will return instance of zap logger configuration, configured to be verbose or to use JSON formatting
func Setup() (logger *zap.Logger, err error) {
	config := zap.Config{
		Level:             zap.NewAtomicLevelAt(zapcore.DebugLevel),
		Development:       false,
		DisableCaller:     false,
		DisableStacktrace: false,
		Sampling:          nil,
		Encoding:          "secureConsole",
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey:     "message",
			LevelKey:       "level",
			TimeKey:        "time",
			NameKey:        "logger",
			CallerKey:      "go",
			StacktraceKey:  "trace",
			LineEnding:     "\n",
			EncodeLevel:    zapcore.CapitalColorLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.StringDurationEncoder,
			EncodeCaller: func(caller zapcore.EntryCaller, enc zapcore.PrimitiveArrayEncoder) {
				callerName := caller.TrimmedPath()
				callerName = filters.MinWidth(callerName, " ", 20)
				enc.AppendString(callerName)
			},
			EncodeName: zapcore.FullNameEncoder,
		},
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: nil,
		InitialFields:    nil,
	}

	return config.Build()
}
