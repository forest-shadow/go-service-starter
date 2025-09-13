package logger

import (
	"os"
	"path/filepath"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"repo-starter/internal/config"
)

// Added an alias to the logger from the zap package
// so it can be easily replaced in the future if necessary.
type Logger = *zap.SugaredLogger

func NewLogger(cfg *config.Config) Logger {
	consoleOutput, fileOutput, errorOutput := getOutputs(cfg)
	consoleEncoder := getConsoleEncoder(cfg)
	fileEncoder := getFileEncoder(cfg)
	level := getLevel(cfg)

	var cores []zapcore.Core

	// Console core - with colors
	consoleCore := zapcore.NewCore(consoleEncoder, consoleOutput, level)
	cores = append(cores, consoleCore)

	// File core - without colors
	if fileOutput != nil {
		fileCore := zapcore.NewCore(fileEncoder, fileOutput, level)
		cores = append(cores, fileCore)
	}

	// Error file core - only errors, without colors
	if errorOutput != nil {
		errorCore := zapcore.NewCore(
			fileEncoder,
			errorOutput,
			zapcore.ErrorLevel, // Only errors and above
		)
		cores = append(cores, errorCore)
	}

	core := zapcore.NewTee(cores...)

	logger := zap.New(core, zap.AddCaller())
	return logger.Sugar()
}

// ANSI color codes
const (
	ColorReset  = "\033[0m"
	ColorPurple = "\033[35m"
	ColorCyan   = "\033[36m"
	ColorGray   = "\033[37m"
)

// Custom colored encoders
func coloredTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(ColorCyan + t.Format("2006-01-02T15:04:05.000Z0700") + ColorReset)
}

func coloredCallerEncoder(caller zapcore.EntryCaller, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(ColorGray + caller.TrimmedPath() + ColorReset)
}

func coloredDurationEncoder(d time.Duration, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(ColorPurple + d.String() + ColorReset)
}

var encoderConfig = zapcore.EncoderConfig{
    TimeKey:        "ts",
	LevelKey:       "level",
	NameKey:        "logger",
	MessageKey:     "message",
	StacktraceKey:  "stacktrace",
	CallerKey:      "caller",
	EncodeLevel:    zapcore.LowercaseLevelEncoder,
	EncodeTime:     zapcore.ISO8601TimeEncoder,
	EncodeDuration: zapcore.SecondsDurationEncoder,
	EncodeCaller:   zapcore.ShortCallerEncoder,
}

func getLevel(cfg *config.Config) zapcore.Level {
	level, err := zapcore.ParseLevel(cfg.Logging.Level.String())
	if err != nil {
		if cfg.Env.IsProduction() {
			return zapcore.InfoLevel
		}
		return zapcore.DebugLevel
	}
	return level
}

func getConsoleEncoder(cfg *config.Config) zapcore.Encoder {
	config := encoderConfig

	// Only use colors in dev and when output is a terminal
	if !cfg.Env.IsProduction() && isTerminal() {
		config.EncodeLevel = zapcore.LowercaseColorLevelEncoder
		config.EncodeTime = coloredTimeEncoder
		config.EncodeCaller = coloredCallerEncoder
		config.EncodeDuration = coloredDurationEncoder
	} else {
		config.EncodeLevel = zapcore.LowercaseLevelEncoder
		config.EncodeTime = zapcore.ISO8601TimeEncoder
		config.EncodeCaller = zapcore.ShortCallerEncoder
		config.EncodeDuration = zapcore.SecondsDurationEncoder
	}

	if cfg.Logging.Format == "json" {
		return zapcore.NewJSONEncoder(config)
	}
	return zapcore.NewConsoleEncoder(config)
}

func getFileEncoder(cfg *config.Config) zapcore.Encoder {
	config := encoderConfig

	// Files always without colors
	config.EncodeLevel = zapcore.LowercaseLevelEncoder
	config.EncodeTime = zapcore.ISO8601TimeEncoder
	config.EncodeCaller = zapcore.ShortCallerEncoder
	config.EncodeDuration = zapcore.SecondsDurationEncoder

	return zapcore.NewJSONEncoder(config)
}

func isTerminal() bool {
	fileInfo, _ := os.Stdout.Stat()
	return (fileInfo.Mode() & os.ModeCharDevice) != 0
}

func getOutputs(cfg *config.Config) (zapcore.WriteSyncer, zapcore.WriteSyncer, zapcore.WriteSyncer) {
	// Console output
	consoleOutput := zapcore.AddSync(os.Stdout)

	// Main log file
	var fileOutput zapcore.WriteSyncer
	if cfg.Logging.File != "" {
		if err := ensureDir(cfg.Logging.File); err == nil {
			if file, err := os.OpenFile(cfg.Logging.File, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644); err == nil {
				fileOutput = zapcore.AddSync(file)
			}
		}
	}

	// Error log file
	var errorOutput zapcore.WriteSyncer
	if cfg.Logging.ErrorFile != "" {
		if err := ensureDir(cfg.Logging.ErrorFile); err == nil {
			if file, err := os.OpenFile(cfg.Logging.ErrorFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644); err == nil {
				errorOutput = zapcore.AddSync(file)
			}
		}
	}

	return consoleOutput, fileOutput, errorOutput
}

func ensureDir(filename string) error {
	dir := filepath.Dir(filename)
	return os.MkdirAll(dir, 0755)
}
