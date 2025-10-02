package upd8

import (
	"bytes"
	"context"
	"fmt"
	"sort"
	"strings"
)

type flatpakManager struct {
	runner CommandRunner
}

func (m *flatpakManager) Name() string { return "flatpak" }

func (m *flatpakManager) Detect(ctx context.Context) bool {
	_ = ctx
	return lookupBinary("flatpak")
}

func (m *flatpakManager) CheckUpdates(ctx context.Context) Result {
	r := Result{Manager: m.Name(), UpdateCommand: "flatpak update"}
	runner := safeRunner(m.runner)

	cmdRes := runner.Run(ctx, "flatpak", "remote-ls", "--updates", "--columns=ref,version")
	if cmdRes.Error != nil && cmdRes.ExitCode != 0 {
		r.Err = fmt.Errorf("flatpak remote-ls failed: %w", cmdRes.Error)
		return r
	}

	payload := bytes.TrimSpace(cmdRes.Stdout)
	if len(payload) == 0 {
		return r
	}

	lines := strings.Split(string(payload), "\n")
	for _, raw := range lines {
		line := strings.TrimSpace(raw)
		if line == "" {
			continue
		}
		lower := strings.ToLower(line)
		if strings.Contains(lower, "ref") && strings.Contains(lower, "version") {
			continue
		}

		parts := strings.Fields(line)
		if len(parts) == 0 {
			continue
		}

		name := parts[0]
		latest := ""
		if len(parts) > 1 {
			latest = parts[len(parts)-1]
		}

		r.Items = append(r.Items, Item{
			Name:   name,
			Latest: latest,
		})
	}

	sort.Slice(r.Items, func(i, j int) bool { return r.Items[i].Name < r.Items[j].Name })
	return r
}
