package main

import (
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	LogMaxSizeMB    = 512
	LogMaxBackupNum = 3
	LogMaxAgeDays   = 30
)

var (
	core zapcore.Core

	highLogPath = strings.Join([]string{"logs", "high_log.log"}, string(filepath.Separator))
	lowLogPath  = strings.Join([]string{"logs", "low_log.log"}, string(filepath.Separator))
)

func init() {
	highPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl == zapcore.ErrorLevel
	})
	lowPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl != zapcore.ErrorLevel
	})

	consoleDebugging := zapcore.Lock(os.Stdout)
	consoleErrors := zapcore.Lock(os.Stderr)
	fileErrorDebugging := zapcore.AddSync(&lumberjack.Logger{
		Filename:   GetHighLogPath(),
		MaxSize:    LogMaxSizeMB,
		MaxBackups: LogMaxBackupNum,
		MaxAge:     LogMaxAgeDays,
		Compress:   true,
	})
	fileInfoDebugging := zapcore.AddSync(&lumberjack.Logger{
		Filename:   GetLowLogPath(),
		MaxSize:    LogMaxSizeMB,
		MaxBackups: LogMaxBackupNum,
		MaxAge:     LogMaxAgeDays,
		Compress:   true,
	})

	consoleEncoder := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
	fileEncoder := zapcore.NewConsoleEncoder(zapcore.EncoderConfig{
		// Keys can be anything except the empty string.
		TimeKey:        "T",
		LevelKey:       "L",
		NameKey:        "N",
		CallerKey:      "C",
		MessageKey:     "M",
		StacktraceKey:  "S",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     timeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	})

	core = zapcore.NewTee(
		zapcore.NewCore(fileEncoder, fileErrorDebugging, highPriority),
		zapcore.NewCore(fileEncoder, fileInfoDebugging, lowPriority),

		zapcore.NewCore(consoleEncoder, consoleErrors, highPriority),
		zapcore.NewCore(consoleEncoder, consoleDebugging, lowPriority),
	)
}

// NewLog 获取日志输出对象
func NewLog() *zap.Logger {
	return zap.New(core)
}

// NewSugaredLogger 获取包装日志输出对象
func NewSugaredLogger() *zap.SugaredLogger {
	return zap.New(core).Sugar()
}

// GetHighLogPath 获得高重要等级的日志文件路径
func GetHighLogPath() string {
	return TmpDir + highLogPath
}

// GetLowLogPath 获得低重要等级的日志文件路径
func GetLowLogPath() string {
	return TmpDir + lowLogPath
}

func timeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05.000"))
}
