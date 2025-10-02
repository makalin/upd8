package upd8

import (
	"bytes"
	"context"
	"fmt"
	"regexp"
	"sort"
	"strings"
)

type cargoManager struct {
	runner CommandRunner
}

func (m *cargoManager) Name() string { return "cargo" }

func (m *cargoManager) Detect(ctx context.Context) bool {
	_ = ctx
	if !lookupBinary("cargo") {
		return false
	}
	// cargo install-update is provided by the cargo-update crate.
	return lookupBinary("cargo-install-update")
}

var cargoLineRegex = regexp.MustCompile(`(?P<name>[^\s]+)\s+v?(?P<current>[0-9][^\s]*)\s+->\s+v?(?P<latest>[0-9][^\s]*)`)

func (m *cargoManager) CheckUpdates(ctx context.Context) Result {
	r := Result{Manager: m.Name(), UpdateCommand: "cargo install-update -a"}
	runner := safeRunner(m.runner)

	cmdRes := runner.Run(ctx, "cargo", "install-update", "--list")
	if cmdRes.Error != nil && cmdRes.ExitCode != 0 {
		r.Err = fmt.Errorf("cargo install-update --list failed: %w", cmdRes.Error)
		return r
	}

	payload := string(bytes.TrimSpace(cmdRes.Stdout))
	if payload == "" {
		return r
	}

	lines := strings.Split(payload, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		matches := cargoLineRegex.FindStringSubmatch(line)
		if len(matches) == 0 {
			continue
		}

		name := matches[cargoLineRegex.SubexpIndex("name")]
		current := matches[cargoLineRegex.SubexpIndex("current")]
		latest := matches[cargoLineRegex.SubexpIndex("latest")]

		r.Items = append(r.Items, Item{ // description omitted to reduce noise
			Name:    name,
			Current: current,
			Latest:  latest,
		})
	}

	sort.Slice(r.Items, func(i, j int) bool { return r.Items[i].Name < r.Items[j].Name })
	return r
}
