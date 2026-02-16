package zerolog_test

import (
	"fmt"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/diode"
)

func ExampleDiodeWriter() {
	wr := diode.NewWriter(os.Stdout, 1000, 10*time.Millisecond, func(missed int) {
		fmt.Printf("Logger Dropped %d messages", missed)
	})
	log := zerolog.New(wr)
	log.Print("test")
}

func ExampleBasicSampler() {
	log := zerolog.New(os.Stdout)
	sampled := log.Sample(&zerolog.BasicSampler{N: 10})
	sampled.Info().Msg("will be logged every 10 messages")
}

func ExampleBurstSampler() {
	log := zerolog.New(os.Stdout)
	sampled := log.Sample(zerolog.LevelSampler{
		DebugSampler: &zerolog.BurstSampler{
			Burst:       5,
			Period:      1 * time.Second,
			NextSampler: &zerolog.BasicSampler{N: 100},
		},
	})
	sampled.Debug().Msg("hello world")
}

func ExamplePredefinedSamplers() {
	log := zerolog.New(os.Stdout)

	log.Sample(zerolog.Often).Info().Msg("~1 in 10")
	log.Sample(zerolog.Sometimes).Info().Msg("~1 in 100")
	log.Sample(zerolog.Rarely).Info().Msg("~1 in 1000")
}
