package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Define ANSI color codes
var levelColors = map[zapcore.Level]string{
	zapcore.DebugLevel:  "\033[36m",   // Cyan
	zapcore.InfoLevel:   "\033[32m",   // Green
	zapcore.WarnLevel:   "\033[33m",   // Yellow
	zapcore.ErrorLevel:  "\033[31m",   // Red
	zapcore.DPanicLevel: "\033[35m",   // Magenta
	zapcore.PanicLevel:  "\033[1;31m", // Bright Red
	zapcore.FatalLevel:  "\033[1;31m", // Bright Red
}

const resetColor = "\033[0m"

var levelsMap = map[string]zapcore.Level{
	"debug":  zapcore.DebugLevel,
	"info":   zapcore.InfoLevel,
	"warn":   zapcore.WarnLevel,
	"error":  zapcore.ErrorLevel,
	"dpanic": zapcore.DPanicLevel,
}

// NewLogger creates a new logger with colored output
func NewLogger(level string) *zap.Logger {
	cfg := zap.NewProductionEncoderConfig()
	cfg.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05")

	// Use the custom level encoder
	cfg.EncodeLevel = func(level zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
		color, ok := levelColors[level]
		if !ok {
			color = resetColor
		}
		enc.AppendString(color + level.CapitalString() + resetColor)
	}

	core := zapcore.NewCore(zapcore.NewConsoleEncoder(cfg), zapcore.AddSync(os.Stdout), levelsMap[level])
	return zap.New(core, zap.AddCaller())
}
