package logging

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/ykhdr/mss-bot/internal/config"
)

// Setup initializes the global logger with the given configuration
func Setup(cfg config.LoggingConfig) {
	lvl := parseLevel(cfg.Level)
	zerolog.SetGlobalLevel(lvl)

	zerolog.CallerMarshalFunc = func(pc uintptr, file string, line int) string {
		// Get the short path relative to the project root
		if idx := strings.Index(file, "mss-bot/"); idx != -1 {
			file = file[idx+len("mss-bot/"):]
		} else {
			// Fallback to filename only
			file = filepath.Base(file)
		}
		return file + ":" + itoa(line)
	}

	writer := zerolog.NewConsoleWriter(func(w *zerolog.ConsoleWriter) {
		w.Out = os.Stdout
		w.TimeFormat = time.RFC3339
	})

	log.Logger = zerolog.
		New(writer).
		With().
		Timestamp().
		Caller().
		Logger()
}

// parseLevel converts a string log level to zerolog.Level
func parseLevel(lvl string) zerolog.Level {
	parsedLevel, err := zerolog.ParseLevel(strings.ToLower(lvl))
	if err != nil || parsedLevel == zerolog.NoLevel {
		return zerolog.InfoLevel
	}
	return parsedLevel
}

// itoa converts an integer to a string without using strconv
func itoa(i int) string {
	if i == 0 {
		return "0"
	}

	var buf [20]byte
	pos := len(buf)
	negative := i < 0
	if negative {
		i = -i
	}

	for i > 0 {
		pos--
		buf[pos] = byte('0' + i%10)
		i /= 10
	}

	if negative {
		pos--
		buf[pos] = '-'
	}

	return string(buf[pos:])
}

// init sets up default logger in case Setup is not called
func init() {
	// Prevent unused import warning
	_ = runtime.Caller
}
