package upd8

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"sort"
)

type pipManager struct {
	runner CommandRunner
	binary string
}

func (m *pipManager) Name() string {
	if m.binary != "" {
		return m.binary
	}
	return "pip"
}

func (m *pipManager) Detect(ctx context.Context) bool {
	_ = ctx
	bin := m.binary
	if bin == "" {
		bin = "pip"
	}
	return lookupBinary(bin)
}

func (m *pipManager) CheckUpdates(ctx context.Context) Result {
	bin := m.binary
	if bin == "" {
		bin = "pip"
	}

	r := Result{Manager: m.Name(), UpdateCommand: fmt.Sprintf("%s install --upgrade -r requirements.txt", bin)}
	runner := safeRunner(m.runner)

	cmdRes := runner.Run(ctx, bin, "list", "--outdated", "--format=json")
	if cmdRes.Error != nil && cmdRes.ExitCode != 0 {
		r.Err = fmt.Errorf("%s list failed: %w", bin, cmdRes.Error)
		return r
	}

	payload := bytes.TrimSpace(cmdRes.Stdout)
	if len(payload) == 0 || bytes.Equal(payload, []byte("[]")) {
		return r
	}

	type pipEntry struct {
		Name          string `json:"name"`
		Version       string `json:"version"`
		LatestVersion string `json:"latest_version"`
	}

	var entries []pipEntry
	if err := json.Unmarshal(payload, &entries); err != nil {
		r.Err = fmt.Errorf("parse %s outdated output: %w", bin, err)
		return r
	}

	sort.Slice(entries, func(i, j int) bool { return entries[i].Name < entries[j].Name })
	for _, entry := range entries {
		r.Items = append(r.Items, Item{
			Name:    entry.Name,
			Current: entry.Version,
			Latest:  entry.LatestVersion,
		})
	}

	return r
}
