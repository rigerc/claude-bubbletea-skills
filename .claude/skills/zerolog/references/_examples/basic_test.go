package zerolog_test

import (
	"errors"
	"fmt"
	stdlog "log"
	"net"
	"os"
	"time"

	"github.com/rs/zerolog"
)

func ExampleNew() {
	log := zerolog.New(os.Stdout)

	log.Info().Msg("hello world")
}

func ExampleLogger_With() {
	log := zerolog.New(os.Stdout).
		With().
		Str("foo", "bar").
		Logger()

	log.Info().Msg("hello world")
}

func ExampleLogger_Level() {
	log := zerolog.New(os.Stdout).Level(zerolog.WarnLevel)

	log.Info().Msg("filtered out message")
	log.Error().Msg("kept message")
}

func ExampleLogger_Sample() {
	log := zerolog.New(os.Stdout).Sample(&zerolog.BasicSampler{N: 2})

	log.Info().Msg("message 1")
	log.Info().Msg("message 2")
	log.Info().Msg("message 3")
	log.Info().Msg("message 4")
}

type LevelNameHook struct{}

func (h LevelNameHook) Run(e *zerolog.Event, l zerolog.Level, msg string) {
	if l != zerolog.NoLevel {
		e.Str("level_name", l.String())
	} else {
		e.Str("level_name", "NoLevel")
	}
}

type MessageHook string

func (h MessageHook) Run(e *zerolog.Event, l zerolog.Level, msg string) {
	e.Str("the_message", msg)
}

func ExampleLogger_Hook() {
	var levelNameHook LevelNameHook
	var messageHook MessageHook = "The message"

	log := zerolog.New(os.Stdout).Hook(levelNameHook, messageHook)

	log.Info().Msg("hello world")
}

func ExampleLogger_Trace() {
	log := zerolog.New(os.Stdout)

	log.Trace().
		Str("foo", "bar").
		Int("n", 123).
		Msg("hello world")
}

func ExampleLogger_Debug() {
	log := zerolog.New(os.Stdout)

	log.Debug().
		Str("foo", "bar").
		Int("n", 123).
		Msg("hello world")
}

func ExampleLogger_Info() {
	log := zerolog.New(os.Stdout)

	log.Info().
		Str("foo", "bar").
		Int("n", 123).
		Msg("hello world")
}

func ExampleLogger_Warn() {
	log := zerolog.New(os.Stdout)

	log.Warn().
		Str("foo", "bar").
		Msg("a warning message")
}

func ExampleLogger_Error() {
	log := zerolog.New(os.Stdout)

	log.Error().
		Err(errors.New("some error")).
		Msg("error doing something")
}

func ExampleLogger_WithLevel() {
	log := zerolog.New(os.Stdout)

	log.WithLevel(zerolog.InfoLevel).
		Msg("hello world")
}

func ExampleLogger_Write() {
	log := zerolog.New(os.Stdout).With().
		Str("foo", "bar").
		Logger()

	stdlog.SetFlags(0)
	stdlog.SetOutput(log)

	stdlog.Print("hello world")
}

func ExampleLogger_Log() {
	log := zerolog.New(os.Stdout)

	log.Log().
		Str("foo", "bar").
		Str("bar", "baz").
		Msg("")
}

func ExampleEvent_Dict() {
	log := zerolog.New(os.Stdout)

	e := log.Log().
		Str("foo", "bar")

	e.Dict("dict", e.CreateDict().
		Str("bar", "baz").
		Int("n", 1),
	).
		Msg("hello world")
}

type User struct {
	Name    string
	Age     int
	Created time.Time
}

func (u User) MarshalZerologObject(e *zerolog.Event) {
	e.Str("name", u.Name).
		Int("age", u.Age).
		Time("created", u.Created)
}

type Price struct {
	val  uint64
	prec int
	unit string
}

func (p Price) MarshalZerologObject(e *zerolog.Event) {
	denom := uint64(1)
	for i := 0; i < p.prec; i++ {
		denom *= 10
	}
	result := []byte(p.unit)
	result = append(result, fmt.Sprintf("%d.%d", p.val/denom, p.val%denom)...)
	e.Str("price", string(result))
}

type Users []User

func (uu Users) MarshalZerologArray(a *zerolog.Array) {
	for _, u := range uu {
		a.Object(u)
	}
}

func ExampleEvent_Array() {
	log := zerolog.New(os.Stdout)

	e := log.Log().
		Str("foo", "bar")

	e.Array("array", e.CreateArray().
		Str("baz").
		Int(1).
		Dict(e.CreateDict().
			Str("bar", "baz").
			Int("n", 1),
		),
	).
		Msg("hello world")
}

func ExampleEvent_Array_object() {
	log := zerolog.New(os.Stdout)

	u := Users{
		User{"John", 35, time.Time{}},
		User{"Bob", 55, time.Time{}},
	}

	log.Log().
		Str("foo", "bar").
		Array("users", u).
		Msg("hello world")
}

func ExampleEvent_Object() {
	log := zerolog.New(os.Stdout)

	u := User{"John", 35, time.Time{}}

	log.Log().
		Str("foo", "bar").
		Object("user", u).
		Msg("hello world")
}

func ExampleEvent_EmbedObject() {
	log := zerolog.New(os.Stdout)

	price := Price{val: 6449, prec: 2, unit: "$"}

	log.Log().
		Str("foo", "bar").
		EmbedObject(price).
		Msg("hello world")
}

func ExampleEvent_Interface() {
	log := zerolog.New(os.Stdout)

	obj := struct {
		Name string `json:"name"`
	}{
		Name: "john",
	}

	log.Log().
		Str("foo", "bar").
		Interface("obj", obj).
		Msg("hello world")
}

func ExampleEvent_Dur() {
	d := 10 * time.Second

	log := zerolog.New(os.Stdout)

	log.Log().
		Str("foo", "bar").
		Dur("dur", d).
		Msg("hello world")
}

func ExampleEvent_Durs() {
	d := []time.Duration{
		10 * time.Second,
		20 * time.Second,
	}

	log := zerolog.New(os.Stdout)

	log.Log().
		Str("foo", "bar").
		Durs("durs", d).
		Msg("hello world")
}

func ExampleEvent_Fields_map() {
	fields := map[string]interface{}{
		"bar": "baz",
		"n":   1,
	}

	log := zerolog.New(os.Stdout)

	log.Log().
		Str("foo", "bar").
		Fields(fields).
		Msg("hello world")
}

func ExampleEvent_Fields_slice() {
	fields := []interface{}{
		"bar", "baz",
		"n", 1,
	}

	log := zerolog.New(os.Stdout)

	log.Log().
		Str("foo", "bar").
		Fields(fields).
		Msg("hello world")
}

func ExampleContext_Dict() {
	ctx := zerolog.New(os.Stdout).With().
		Str("foo", "bar")

	logger := ctx.Dict("dict", ctx.CreateDict().
		Str("bar", "baz").
		Int("n", 1),
	).Logger()

	logger.Log().Msg("hello world")
}

func ExampleContext_Array() {
	ctx := zerolog.New(os.Stdout).With().
		Str("foo", "bar")

	logger := ctx.Array("array", ctx.CreateArray().
		Str("baz").
		Int(1),
	).Logger()

	logger.Log().Msg("hello world")
}

func ExampleContext_IPAddr() {
	hostIP := net.IP{192, 168, 0, 100}
	log := zerolog.New(os.Stdout).With().
		IPAddr("HostIP", hostIP).
		Logger()

	log.Log().Msg("hello world")
}

func ExampleContext_IPPrefix() {
	route := net.IPNet{IP: net.IP{192, 168, 0, 0}, Mask: net.CIDRMask(24, 32)}
	log := zerolog.New(os.Stdout).With().
		IPPrefix("Route", route).
		Logger()

	log.Log().Msg("hello world")
}

func ExampleContext_MACAddr() {
	mac := net.HardwareAddr{0x00, 0x14, 0x22, 0x01, 0x23, 0x45}
	log := zerolog.New(os.Stdout).With().
		MACAddr("hostMAC", mac).
		Logger()

	log.Log().Msg("hello world")
}
