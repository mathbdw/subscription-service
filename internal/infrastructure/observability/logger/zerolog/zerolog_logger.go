package zerolog

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	zlog "github.com/rs/zerolog"

	"github.com/mathbdw/subscription-service/config"
	"github.com/mathbdw/subscription-service/internal/interfaces/observability"
)

type logger struct {
	level *atomic.Int32
	mu    sync.RWMutex

	baseLogger zlog.Logger
}

// New - Constructor logger
func New(cfg *config.Config) observability.Logger {
	skipFrameCount := 1

	level := &atomic.Int32{}

	var zlevel zlog.Level
	if cfg.Project.Debug {
		zlevel = zlog.DebugLevel
		level.Store(int32(zlog.DebugLevel))
	} else {
		zlevel = zlog.InfoLevel
		level.Store(int32(zlog.InfoLevel))
	}

	zlog.TimeFieldFormat = zlog.TimeFormatUnix

	// Console writer
	consoleWriter := zlog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
	consoleWriter.FormatLevel = func(i any) string {
		return strings.ToUpper(fmt.Sprintf("| %-6s|", i))
	}
	consoleWriter.FormatMessage = func(i any) string {
		return fmt.Sprintf("***%s***", i)
	}
	consoleWriter.FormatFieldName = func(i any) string {
		return fmt.Sprintf("%s:", i)
	}
	consoleWriter.FormatFieldValue = func(i any) string {
		return strings.ToUpper(fmt.Sprintf("%s", i))
	}

	baseLogger := zlog.New(consoleWriter).With().Timestamp().
		CallerWithSkipFrameCount(zlog.CallerSkipFrameCount + skipFrameCount).
		Logger().Level(zlevel)

	return &logger{
		level:      level,
		baseLogger: baseLogger,
	}
}

// SetLevel - set new level
func (l *logger) SetLevel(newLevel int8) {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.level.Store(int32(newLevel))
	l.baseLogger = l.baseLogger.Level(zlog.Level(newLevel))
}

// Debug - implementation of Debug for zerolog
func (l *logger) Debug(msg string, fields observability.Field) {
	l.mu.RLock()
	defer l.mu.RUnlock()

	Msg(l.baseLogger.Debug(), fields).Msg(msg)
}

// Info - implementation of Info for zerolog
func (l *logger) Info(msg string, fields observability.Field) {
	l.mu.RLock()
	defer l.mu.RUnlock()

	Msg(l.baseLogger.Info(), fields).Msg(msg)
}

// Warn - implementation of Warn for zerolog
func (l *logger) Warn(msg string, fields observability.Field) {
	l.mu.RLock()
	defer l.mu.RUnlock()

	Msg(l.baseLogger.Warn(), fields).Msg(msg)
}

// Error - implementation of Error for zerolog
func (l *logger) Error(msg string, fields observability.Field) {
	l.mu.RLock()
	defer l.mu.RUnlock()

	Msg(l.baseLogger.Error(), fields).Msg(msg)
}

// Fatal - implementation of Fatal for zerolog
func (l *logger) Fatal(msg string, fields observability.Field) {
	l.mu.RLock()
	defer l.mu.RUnlock()

	Msg(l.baseLogger.Error(), fields).Msg(msg)
	os.Exit(1)
}

// Msg - adds fields for the zerolog event
func Msg(event *zlog.Event, fields observability.Field) *zlog.Event {
	for key, value := range fields {
		switch v := value.(type) {
		case int:
			event.Int(key, v)
		case int8:
			event.Int8(key, v)
		case int16:
			event.Int16(key, v)
		case int32:
			event.Int32(key, v)
		case int64:
			event.Int64(key, v)
		case []int:
			event.Ints(key, v)
		case []int8:
			event.Ints8(key, v)
		case []int16:
			event.Ints16(key, v)
		case []int32:
			event.Ints32(key, v)
		case []int64:
			event.Ints64(key, v)
		case uint8:
			event.Uint8(key, v)
		case uint16:
			event.Uint16(key, v)
		case uint32:
			event.Uint32(key, v)
		case uint64:
			event.Uint64(key, v)
		case []uint8:
			event.Uints8(key, v)
		case []uint16:
			event.Uints16(key, v)
		case []uint32:
			event.Uints32(key, v)
		case []uint64:
			event.Uints64(key, v)
		case string:
			event.Str(key, v)
		case []string:
			event.Strs(key, v)
		case bool:
			event.Bool(key, v)
		default:
			event.Str(key, fmt.Sprintf("%+v", v))
		}
	}

	return event
}
