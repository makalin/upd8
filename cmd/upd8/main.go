package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/makalin/upd8/internal/upd8"
)

func main() {
	os.Exit(run(os.Args[1:]))
}

func run(args []string) int {
	fs := flag.NewFlagSet("upd8", flag.ContinueOnError)
	watch := fs.Bool("watch", false, "Run in daemon mode, printing summaries at each interval")
	interval := fs.Duration("interval", 24*time.Hour, "Scan interval when running with --watch")
	showPackages := fs.Bool("packages", false, "Render a short list of outdated packages per manager")
	noColor := fs.Bool("no-color", false, "Disable ANSI colors in the output")
	verbose := fs.Bool("verbose", false, "Include managers even when no updates are found")

	if err := fs.Parse(args); err != nil {
		return 2
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)
	go func() {
		<-sigCh
		fmt.Fprintln(os.Stderr, "\nInterrupted, exiting...")
		cancel()
	}()

	runner := upd8.ExecRunner{Timeout: 60 * time.Second}
	scanner := upd8.Scanner{Runner: runner, Managers: upd8.DefaultManagers(runner)}

	renderer := upd8.Renderer{
		Writer:       os.Stdout,
		EnableColor:  !*noColor,
		ShowPackages: *showPackages,
	}

	if *verbose {
		renderer.EmptyMessage = "No supported package managers detected."
	} else {
		renderer.EmptyMessage = "No updates found. (Use --verbose to show all managers.)"
	}

	if *watch {
		if *interval <= 0 {
			fmt.Fprintln(os.Stderr, "interval must be positive when using --watch")
			return 2
		}

		renderer.Timestamp = true
		fmt.Fprintf(os.Stdout, "Watching for updates every %s. Press Ctrl+C to stop.\n", (*interval).Truncate(time.Second))

		scanner.Watch(ctx, *interval, func(results []upd8.Result) {
			batch := results
			if !*verbose {
				batch = filterEmpty(batch)
			}
			renderer.Render(batch)
		})
		return 0
	}

	results := scanner.Scan(ctx)
	if !*verbose {
		results = filterEmpty(results)
	}

	renderer.Render(results)
	return computeExitCode(results)
}

func filterEmpty(results []upd8.Result) []upd8.Result {
	filtered := make([]upd8.Result, 0, len(results))
	for _, r := range results {
		if r.Err != nil || len(r.Items) > 0 {
			filtered = append(filtered, r)
		}
	}
	return filtered
}

func computeExitCode(results []upd8.Result) int {
	for _, r := range results {
		if r.Err != nil {
			return 1
		}
	}
	return 0
}
