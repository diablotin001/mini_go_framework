package logger

import (
    "os"
    "path/filepath"
    "time"

    "go.uber.org/zap"
    "go.uber.org/zap/zapcore"
    "gopkg.in/natefinch/lumberjack.v2"
)

func InitLogger(path string) {
    if path == "" {
        path = "logs/app.log"
    }
    _ = os.MkdirAll(filepath.Dir(path), 0755)

    writeSyncer := getLogWriter(path)
    encoder := getEncoder()

    core := zapcore.NewCore(
        encoder,
        writeSyncer,
        zapcore.InfoLevel,
    )

    logger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
    zap.ReplaceGlobals(logger)
}

func getEncoder() zapcore.Encoder {
	return zapcore.NewJSONEncoder(zapcore.EncoderConfig{
		MessageKey: "msg",
		LevelKey:   "level",
		TimeKey:    "time",
		NameKey:    "logger",
		CallerKey:  "caller",

		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     func(t time.Time, enc zapcore.PrimitiveArrayEncoder) { enc.AppendString(t.Format(time.RFC3339)) },
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	})
}

func getLogWriter(path string) zapcore.WriteSyncer {
    lumberJackLogger := &lumberjack.Logger{
        Filename:   path,
        MaxSize:    20,
        MaxBackups: 30,
        MaxAge:     7,
        Compress:   true,
    }
    return zapcore.AddSync(lumberJackLogger)
}
