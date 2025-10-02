package upd8

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"sort"
)

type npmManager struct {
	runner CommandRunner
}

func (m *npmManager) Name() string { return "npm" }

func (m *npmManager) Detect(ctx context.Context) bool {
	_ = ctx
	return lookupBinary("npm")
}

func (m *npmManager) CheckUpdates(ctx context.Context) Result {
	r := Result{Manager: m.Name(), UpdateCommand: "npm update -g"}
	runner := safeRunner(m.runner)

	cmdRes := runner.Run(ctx, "npm", "outdated", "-g", "--json")
	if cmdRes.Error != nil && cmdRes.ExitCode != 0 && cmdRes.ExitCode != 1 {
		r.Err = fmt.Errorf("npm outdated failed: %w", cmdRes.Error)
		return r
	}

	payload := bytes.TrimSpace(cmdRes.Stdout)
	if len(payload) == 0 || bytes.Equal(payload, []byte("null")) {
		return r
	}

	type npmEntry struct {
		Current string `json:"current"`
		Wanted  string `json:"wanted"`
		Latest  string `json:"latest"`
	}

	entries := map[string]npmEntry{}
	if err := json.Unmarshal(payload, &entries); err != nil {
		r.Err = fmt.Errorf("parse npm outdated output: %w", err)
		return r
	}

	keys := make([]string, 0, len(entries))
	for name := range entries {
		keys = append(keys, name)
	}
	sort.Strings(keys)

	for _, name := range keys {
		entry := entries[name]
		// npm sometimes returns wanted/Latest empty; fall back to wanted.
		latest := entry.Latest
		if latest == "" {
			latest = entry.Wanted
		}
		r.Items = append(r.Items, Item{
			Name:    name,
			Current: entry.Current,
			Latest:  latest,
		})
	}

	return r
}
