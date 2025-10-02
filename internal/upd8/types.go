package upd8

import (
	"context"
)

// Item describes a single outdated package.
type Item struct {
	Name        string
	Current     string
	Latest      string
	Description string
}

// Result captures the outcome of running an update check for a package manager.
type Result struct {
	Manager       string
	Items         []Item
	UpdateCommand string
	Err           error
	DurationMs    int64
}

// Manager defines the capabilities of a package manager implementation.
type Manager interface {
	Name() string
	Detect(ctx context.Context) bool
	CheckUpdates(ctx context.Context) Result
}
