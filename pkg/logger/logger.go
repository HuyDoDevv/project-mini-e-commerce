package logger

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
	"project-mini-e-commerce/internal/common"
	"time"

	"github.com/natefinch/lumberjack"
	"github.com/rs/zerolog"
)

type Config struct {
	Level       string
	Filename    string
	MaxSize     int
	MaxAge      int
	MaxBackups  int
	Compress    bool
	Environment common.Environment
}

func NewLogger(config Config) *zerolog.Logger {
	zerolog.TimeFieldFormat = time.RFC3339

	lvl, err := zerolog.ParseLevel(config.Level)
	if err != nil {
		lvl = zerolog.InfoLevel
	}
	zerolog.SetGlobalLevel(lvl)
	write := getWriter(config)
	logger := zerolog.New(write).With().Timestamp().Logger()
	return &logger
}

type PrettyJSONWrite struct {
	Writer io.Writer
}

func (w *PrettyJSONWrite) Write(p []byte) (n int, err error) {
	var prettyJSON bytes.Buffer
	err = json.Indent(&prettyJSON, p, "", "  ")
	if err != nil {
		return w.Writer.Write(p)
	}
	return w.Writer.Write(prettyJSON.Bytes())
}

func getWriter(config Config) io.Writer {
	switch config.Environment {
	case common.Development, common.Staging:
		return &PrettyJSONWrite{Writer: os.Stdout}
	case common.Production:
		return &lumberjack.Logger{
			Filename:   config.Filename,
			MaxSize:    config.MaxSize,
			MaxAge:     config.MaxAge,
			MaxBackups: config.MaxBackups,
			Compress:   config.Compress,
		}
	default:
		return &PrettyJSONWrite{Writer: os.Stdout}
	}
}
