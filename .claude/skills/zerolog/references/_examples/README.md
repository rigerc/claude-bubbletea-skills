# Zerolog Examples

This directory contains example code demonstrating zerolog features.

## Files

| File | Description |
|------|-------------|
| `basic_test.go` | Basic logging, log levels, structured fields, objects, arrays, dicts |
| `error_logging_test.go` | Error logging with stack traces, fatal logging |
| `writers_test.go` | ConsoleWriter, MultiLevelWriter, file output |
| `context_test.go` | context.Context integration, context-based logging |
| `hlog_test.go` | HTTP handler integration with hlog package |
| `sampling_test.go` | Log sampling, diode writers for high-throughput |
| `hooks_test.go` | Custom hooks for adding fields to log events |

## Running Examples

These are Go test files. To run them in a project that uses zerolog:

```bash
go get github.com/rs/zerolog
go test -v ./...
```

## Key Patterns

### Always terminate the chain
```go
log.Info().Str("key", "val").Msg("done")
log.Info().Str("key", "val").Send()
```

### Conditional expensive operations
```go
if e := log.Debug(); e.Enabled() {
    e.Str("computed", expensiveValue()).Msg("debug")
}
```

### Contextual loggers
```go
sublogger := log.With().Str("component", "api").Logger()
sublogger.Info().Msg("request handled")
```

### Concurrency safety
```go
logger := log.Logger.With().Logger()
logger.UpdateContext(func(c zerolog.Context) zerolog.Context {
    return c.Str("request_id", id)
})
```
