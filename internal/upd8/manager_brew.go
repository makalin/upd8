package upd8

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"sort"
)

type brewManager struct {
	runner CommandRunner
}

func (m *brewManager) Name() string { return "brew" }

func (m *brewManager) Detect(ctx context.Context) bool {
	_ = ctx
	return lookupBinary("brew")
}

func (m *brewManager) CheckUpdates(ctx context.Context) Result {
	r := Result{Manager: m.Name(), UpdateCommand: "brew upgrade"}
	runner := safeRunner(m.runner)

	cmdRes := runner.Run(ctx, "brew", "outdated", "--json=v2")
	if cmdRes.Error != nil && cmdRes.ExitCode != 0 {
		r.Err = fmt.Errorf("brew outdated failed: %w", cmdRes.Error)
		return r
	}

	payload := bytes.TrimSpace(cmdRes.Stdout)
	if len(payload) == 0 {
		return r
	}

	var parsed struct {
		Formulae []struct {
			Name              string   `json:"name"`
			InstalledVersions []string `json:"installed_versions"`
			CurrentVersion    string   `json:"current_version"`
		} `json:"formulae"`
		Casks []struct {
			Name             string `json:"name"`
			InstalledVersion string `json:"installed_version"`
			CurrentVersion   string `json:"current_version"`
		} `json:"casks"`
	}

	if err := json.Unmarshal(payload, &parsed); err != nil {
		r.Err = fmt.Errorf("parse brew outdated output: %w", err)
		return r
	}

	type row struct {
		name    string
		current string
		latest  string
	}
	var rows []row

	for _, formula := range parsed.Formulae {
		current := ""
		if len(formula.InstalledVersions) > 0 {
			current = formula.InstalledVersions[len(formula.InstalledVersions)-1]
		}
		rows = append(rows, row{
			name:    formula.Name,
			current: current,
			latest:  formula.CurrentVersion,
		})
	}

	for _, cask := range parsed.Casks {
		rows = append(rows, row{
			name:    cask.Name,
			current: cask.InstalledVersion,
			latest:  cask.CurrentVersion,
		})
	}

	sort.Slice(rows, func(i, j int) bool { return rows[i].name < rows[j].name })

	for _, row := range rows {
		if row.name == "" {
			continue
		}
		r.Items = append(r.Items, Item{
			Name:    row.name,
			Current: row.current,
			Latest:  row.latest,
		})
	}

	return r
}
