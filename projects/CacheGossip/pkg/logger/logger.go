package logger

import (
	"fmt"
	"os"

	"github.com/rs/zerolog"
)

type Logger struct {
	log zerolog.Logger
}

func NewLogger(name string) (*Logger, error) {
	if name == "" {
		return nil, fmt.Errorf("Name is empty")
	}
	l := zerolog.New(os.Stdout).With().Timestamp().Str("NodeName", name).Logger()
	return &Logger{l}, nil
}

func (l *Logger) Info(str string, arg ...any) {
	l.msg(zerolog.InfoLevel, str, arg...)
}

func (l *Logger) Warning(str string, arg ...any) {
	l.msg(zerolog.WarnLevel, str, arg...)
}

func (l *Logger) Error(str string, arg ...any) {
	l.msg(zerolog.ErrorLevel, str, arg...)
}

func (l *Logger) msg(level zerolog.Level, str string, arg ...any) {
	l.log.WithLevel(level).Msg(fmt.Sprintf(str, arg...))
}
