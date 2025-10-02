package upd8

import (
	"bytes"
	"context"
	"fmt"
	"sort"
	"strings"
)

type snapManager struct {
	runner CommandRunner
}

func (m *snapManager) Name() string { return "snap" }

func (m *snapManager) Detect(ctx context.Context) bool {
	_ = ctx
	return lookupBinary("snap")
}

func (m *snapManager) CheckUpdates(ctx context.Context) Result {
	r := Result{Manager: m.Name(), UpdateCommand: "snap refresh"}
	runner := safeRunner(m.runner)

	cmdRes := runner.Run(ctx, "snap", "refresh", "--list")
	if cmdRes.Error != nil && cmdRes.ExitCode != 0 {
		r.Err = fmt.Errorf("snap refresh --list failed: %w", cmdRes.Error)
		return r
	}

	payload := bytes.TrimSpace(cmdRes.Stdout)
	if len(payload) == 0 {
		return r
	}

	lines := strings.Split(string(payload), "\n")
	for idx, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		// Skip header line containing column markers.
		if idx == 0 && strings.Contains(strings.ToLower(line), "version") {
			continue
		}

		parts := strings.Fields(line)
		if len(parts) < 2 {
			continue
		}

		name := parts[0]
		latest := parts[1]

		r.Items = append(r.Items, Item{
			Name:   name,
			Latest: latest,
		})
	}

	sort.Slice(r.Items, func(i, j int) bool { return r.Items[i].Name < r.Items[j].Name })
	return r
}
