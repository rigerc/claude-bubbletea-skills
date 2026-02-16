package zerolog_test

import (
	"net/http"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/hlog"
)

func ExampleHlog_Basic() {
	log := zerolog.New(os.Stdout).With().
		Timestamp().
		Str("role", "my-service").
		Logger()

	handler := hlog.NewHandler(log)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hlog.FromRequest(r).Info().
			Str("path", r.URL.Path).
			Msg("request received")
	}))
}

func ExampleHlog_WithMiddleware() {
	log := zerolog.New(os.Stdout).With().
		Timestamp().
		Str("role", "my-service").
		Logger()

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hlog.FromRequest(r).Info().
			Str("user", "current user").
			Str("status", "ok").
			Msg("Something happened")
	})

	handler = hlog.NewHandler(log)(handler)
	handler = hlog.RemoteAddrHandler("ip")(handler)
	handler = hlog.UserAgentHandler("user_agent")(handler)
	handler = hlog.RefererHandler("referer")(handler)
	handler = hlog.AccessHandler(func(r *http.Request, status, size int, duration time.Duration) {
		hlog.FromRequest(r).Info().
			Str("method", r.Method).
			Stringer("url", r.URL).
			Int("status", status).
			Int("size", size).
			Dur("duration", duration).
			Msg("")
	})(handler)
}
