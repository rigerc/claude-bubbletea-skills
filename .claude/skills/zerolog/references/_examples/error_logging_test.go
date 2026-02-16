package zerolog_test

import (
	"errors"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
)

func ExampleErrorLogging() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	err := errors.New("seems we have an error here")
	log.Error().Err(err).Msg("")
}

func ExampleErrorLogging_Stacktrace() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack

	err := outer()
	log.Error().Stack().Err(err).Msg("")
}

func inner() error {
	return errors.New("seems we have an error here")
}

func middle() error {
	err := inner()
	if err != nil {
		return err
	}
	return nil
}

func outer() error {
	err := middle()
	if err != nil {
		return err
	}
	return nil
}

func ExampleFatalLogging() {
	err := errors.New("A repo man spends his life getting into tense situations")
	service := "myservice"

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	log.Fatal().
		Err(err).
		Str("service", service).
		Msgf("Cannot start %s", service)
}
