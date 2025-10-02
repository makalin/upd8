package upd8

import "context"

// DefaultManagers returns the built-in package manager implementations.
func DefaultManagers(runner CommandRunner) []Manager {
	return []Manager{
		&npmManager{runner: runner},
		&pipManager{runner: runner, binary: "pip"},
		&pipManager{runner: runner, binary: "pip3"},
		&brewManager{runner: runner},
		&cargoManager{runner: runner},
		&flatpakManager{runner: runner},
		&snapManager{runner: runner},
	}
}

// safeRunner ensures managers always have a runner to use.
func safeRunner(runner CommandRunner) CommandRunner {
	if runner != nil {
		return runner
	}
	return ExecRunner{}
}

// detectBinary wraps exec.LookPath for tests.
func detectBinary(_ context.Context, binary string) bool {
	return lookupBinary(binary)
}
