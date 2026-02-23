// Package validator runs shell validation commands and reports results.
package validator

import (
	"context"
	"os/exec"
	"strings"
	"time"
)

// Result holds the outcome of a single validation command execution.
type Result struct {
	Command  string
	Passed   bool
	Output   string
	Duration time.Duration
}

// Run executes command (split by whitespace) and returns the result.
// The first word is the executable; remaining words are arguments.
// Combined stdout and stderr are captured in Result.Output.
// Passed is true when the process exits with code 0.
func Run(ctx context.Context, command string) Result {
	fields := strings.Fields(command)
	r := Result{Command: command}

	if len(fields) == 0 {
		r.Output = "empty command"
		return r
	}

	start := time.Now()
	cmd := exec.CommandContext(ctx, fields[0], fields[1:]...)
	out, err := cmd.CombinedOutput()
	r.Duration = time.Since(start)
	r.Output = string(out)
	r.Passed = err == nil
	return r
}
