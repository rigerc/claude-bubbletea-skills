package adapter

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"slices"
	"strings"
	"sync"
)

var (
	modelCacheMu sync.Mutex
	modelCache   = map[AgentType][]string{}
)

// FetchModels returns the list of available models for the given agent. Results
// are cached in memory. Returns an empty slice for agents not in
// AgentsSupportingModel. Ported from ralph/src/lib/model-fetcher.ts.
func FetchModels(ctx context.Context, agent AgentType) ([]string, error) {
	if !supportsModelSelection(agent) {
		return []string{}, nil
	}

	modelCacheMu.Lock()
	if cached, ok := modelCache[agent]; ok {
		modelCacheMu.Unlock()
		return cached, nil
	}
	modelCacheMu.Unlock()

	models, err := fetchModelsUncached(ctx, agent)
	if err != nil {
		return nil, err
	}

	modelCacheMu.Lock()
	modelCache[agent] = models
	modelCacheMu.Unlock()

	return models, nil
}

// ClearModelCache removes all cached model lists, forcing the next call to
// FetchModels to re-query the agent binaries.
func ClearModelCache() {
	modelCacheMu.Lock()
	defer modelCacheMu.Unlock()
	modelCache = map[AgentType][]string{}
}

func supportsModelSelection(agent AgentType) bool {
	return slices.Contains(AgentsSupportingModel, agent)
}

func fetchModelsUncached(ctx context.Context, agent AgentType) ([]string, error) {
	switch agent {
	case AgentOpencode, AgentKilo:
		return fetchOpencodeStyleModels(ctx, string(agent))
	case AgentPi:
		return fetchPiModels(ctx)
	default:
		return nil, fmt.Errorf("model listing not supported for agent %q", agent)
	}
}

// fetchOpencodeStyleModels runs `<binary> models` and returns one model per
// non-empty line.
func fetchOpencodeStyleModels(ctx context.Context, binary string) ([]string, error) {
	out, err := runCommand(ctx, binary, "models")
	if err != nil {
		return nil, fmt.Errorf("fetching models for %s: %w", binary, err)
	}

	var models []string
	sc := bufio.NewScanner(bytes.NewReader(out))
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line != "" {
			models = append(models, line)
		}
	}
	return models, nil
}

// fetchPiModels runs `pi --list-models` and parses lines of the form
// "provider model" into "provider/model" strings, skipping header/warning
// lines.
func fetchPiModels(ctx context.Context) ([]string, error) {
	out, err := runCommand(ctx, "pi", "--list-models")
	if err != nil {
		return nil, fmt.Errorf("fetching pi models: %w", err)
	}

	var models []string
	sc := bufio.NewScanner(bytes.NewReader(out))
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" {
			continue
		}
		// Skip header and status lines.
		lower := strings.ToLower(line)
		if strings.HasPrefix(lower, "provider") ||
			strings.HasPrefix(lower, "warning") ||
			strings.HasPrefix(lower, "error") {
			continue
		}
		parts := strings.Fields(line)
		if len(parts) >= 2 {
			models = append(models, parts[0]+"/"+parts[1])
		}
	}
	return models, nil
}

// runCommand executes a binary with the given arguments and returns combined
// stdout output.
func runCommand(ctx context.Context, binary string, args ...string) ([]byte, error) {
	cmd := exec.CommandContext(ctx, binary, args...)
	cmd.Env = os.Environ()
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	return out, nil
}
