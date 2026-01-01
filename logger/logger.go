package logger

import (
    "os"
    "path/filepath"
    "time"

    "go.uber.org/zap"
    "go.uber.org/zap/zapcore"
    "gopkg.in/natefinch/lumberjack.v2"
)

var Log *zap.Logger

type LogConfig struct {
    Level  string
    Format string
    File   string
    Env    string
}

func Init(cfg LogConfig) error {
    if cfg.File == "" {
        cfg.File = "logs/app.log"
    }
    _ = os.MkdirAll(filepath.Dir(cfg.File), 0755)
    var level zapcore.Level
    if err := level.Set(cfg.Level); err != nil {
        level = zap.InfoLevel
    }
    var encoder zapcore.Encoder
    if cfg.Format == "json" || cfg.Env == "prod" {
        encCfg := zap.NewProductionEncoderConfig()
        encCfg.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) { enc.AppendString(t.Format(time.RFC3339)) }
        encoder = zapcore.NewJSONEncoder(encCfg)
    } else {
        encCfg := zap.NewDevelopmentEncoderConfig()
        encCfg.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) { enc.AppendString(t.Format(time.RFC3339)) }
        encoder = zapcore.NewConsoleEncoder(encCfg)
    }
    fileSyncer := getLogWriter(cfg.File)
    ws := zapcore.NewMultiWriteSyncer(fileSyncer, zapcore.AddSync(os.Stdout))
    core := zapcore.NewCore(encoder, ws, level)
    Log = zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
    zap.ReplaceGlobals(Log)
    return nil
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
