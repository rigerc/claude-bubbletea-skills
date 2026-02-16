package zerolog_test

import (
	"context"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func ExampleContext_Basic() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	ctx := log.With().Str("component", "module").Logger().WithContext(context.Background())

	log.Ctx(ctx).Info().Msg("hello world")
}

func ExampleContext_AttachRetrieve() {
	logger := zerolog.New(os.Stdout)
	ctx := context.Background()

	ctx = logger.WithContext(ctx)

	retrieved := zerolog.Ctx(ctx)
	retrieved.Info().Msg("Hello")
}

func ExampleContext_WithHook() {
	type TracingHook struct{}

	getSpanIdFromContext := func(ctx context.Context) string { return "span-123" }
	_ = getSpanIdFromContext

	var h TracingHook
	logger := zerolog.New(os.Stdout).Hook(h)

	ctx := context.Background()
	logger.Info().Ctx(ctx).Msg("Hello")
}
