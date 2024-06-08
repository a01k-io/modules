package logger

import (
	"context"
	"os"
	"time"

	"github.com/rs/zerolog"
)

type Wrapper struct {
	lg zerolog.Logger
}

func New(service string) *Wrapper {
	out := zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: time.RFC3339,
		NoColor:    false,
	}
	fields := map[string]interface{}{"service": service}
	logger := zerolog.New(out).
		Level(zerolog.GlobalLevel()).
		With().
		Stack().
		Fields(fields).
		Timestamp().
		Logger()

	return &Wrapper{
		lg: logger,
	}
}

func (logger *Wrapper) Debug(message string) {
	logger.lg.Debug().Caller(1, 2, 3).Msg(message)
}

func (logger *Wrapper) Info(message string) {
	logger.lg.Info().Caller(1, 2, 3).Msg(message)
}

func (logger *Wrapper) Warn(message string) {
	logger.lg.Warn().Caller(1, 2, 3).Msg(message)
}

func (logger *Wrapper) Error(message string) {
	logger.lg.Error().Caller(1, 2, 3).Msg(message)
}

func (logger *Wrapper) Fatal(message string) {
	logger.lg.Fatal().Caller(1, 2, 3).Msg(message)
}

func (logger *Wrapper) Debugf(format string, args ...interface{}) {
	logger.lg.Debug().Caller(1, 2, 3).Msgf(format, args...)
}

func (logger *Wrapper) Infof(format string, args ...interface{}) {
	logger.lg.Info().Caller(1, 2, 3).Msgf(format, args...)
}

func (logger *Wrapper) Warnf(format string, args ...interface{}) {
	logger.lg.Warn().Caller(1, 2, 3).Msgf(format, args...)
}

func (logger *Wrapper) Errorf(format string, args ...interface{}) {
	logger.lg.Error().Caller(1, 2, 3).Msgf(format, args...)
}

func (logger *Wrapper) Fatalf(format string, args ...interface{}) {
	logger.lg.Fatal().Caller(1, 2, 3).Msgf(format, args...)
}

func (logger *Wrapper) DebugfCtx(ctx context.Context, format string, args ...interface{}) {
	logger.lg.Debug().Caller(1, 2, 3).Ctx(ctx).Msgf(format, args...)
}

func (logger *Wrapper) InfofCtx(ctx context.Context, format string, args ...interface{}) {
	logger.lg.Info().Caller(1, 2, 3).Ctx(ctx).Msgf(format, args...)
}

func (logger *Wrapper) WarnfCtx(ctx context.Context, format string, args ...interface{}) {
	logger.lg.Warn().Caller(1, 2, 3).Ctx(ctx).Msgf(format, args...)
}

func (logger *Wrapper) ErrorfCtx(ctx context.Context, format string, args ...interface{}) {
	logger.lg.Error().Caller(1, 2, 3).Ctx(ctx).Msgf(format, args...)
}

func (logger *Wrapper) FatalfCtx(ctx context.Context, format string, args ...interface{}) {
	logger.lg.Fatal().Caller(1, 2, 3).Ctx(ctx).Msgf(format, args...)
}
