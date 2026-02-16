package zerolog_test

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type SeverityHook struct{}

func (h SeverityHook) Run(e *zerolog.Event, level zerolog.Level, msg string) {
	if level != zerolog.NoLevel {
		e.Str("severity", level.String())
	}
}

func ExampleHook_Basic() {
	hooked := log.Hook(SeverityHook{})
	hooked.Warn().Msg("")
}

func ExampleHook_Multiple() {
	type FieldHook struct{}

	h := FieldHook{}
	_ = h

	logger := zerolog.New(os.Stdout).Hook(SeverityHook{})
	logger.Info().Msg("message with hooks")
}

func ExampleHookFunc() {
	logger := zerolog.New(os.Stdout).Hook(zerolog.HookFunc(func(e *zerolog.Event, level zerolog.Level, msg string) {
		e.Str("hooked", "true")
	}))
	logger.Info().Msg("message")
}

func ExampleLevelHook() {
	logger := zerolog.New(os.Stdout).Hook(zerolog.NewLevelHook(
		zerolog.HookFunc(func(e *zerolog.Event, level zerolog.Level, msg string) {
			e.Bool("is_error", level == zerolog.ErrorLevel)
		}),
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
	))
	logger.Error().Msg("error occurred")
}
