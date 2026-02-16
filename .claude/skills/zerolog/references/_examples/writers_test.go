package zerolog_test

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog"
)

func ExampleConsoleWriter_Basic() {
	log := zerolog.New(os.Stderr).Output(zerolog.ConsoleWriter{Out: os.Stderr})

	log.Info().Str("foo", "bar").Msg("Hello World")
}

func ExampleConsoleWriter_CustomFormat() {
	output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
	output.FormatLevel = func(i interface{}) string {
		return strings.ToUpper(fmt.Sprintf("| %-6s|", i))
	}
	output.FormatMessage = func(i interface{}) string {
		return fmt.Sprintf("***%s****", i)
	}
	output.FormatFieldName = func(i interface{}) string {
		return fmt.Sprintf("%s:", i)
	}
	output.FormatFieldValue = func(i interface{}) string {
		return strings.ToUpper(fmt.Sprintf("%s", i))
	}

	log := zerolog.New(output).With().Timestamp().Logger()
	log.Info().Str("foo", "bar").Msg("Hello World")
}

func ExampleMultiLevelWriter() {
	consoleWriter := zerolog.ConsoleWriter{Out: os.Stdout}
	multi := zerolog.MultiLevelWriter(consoleWriter, os.Stdout)
	logger := zerolog.New(multi).With().Timestamp().Logger()
	logger.Info().Msg("Hello World!")
}

func ExampleFileOutput() {
	file, _ := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	logger := zerolog.New(file).With().Timestamp().Logger()
	logger.Info().Msg("logging to file")
}
