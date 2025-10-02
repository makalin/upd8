package upd8

import (
	"context"
	"os/exec"
)

func lookupBinary(name string) bool {
	_, err := exec.LookPath(name)
	return err == nil
}

// detectViaCommands checks that all commands/binaries exist.
func detectViaCommands(ctx context.Context, binaries ...string) bool {
	for _, bin := range binaries {
		if !lookupBinary(bin) {
			return false
		}
	}
	return true
}
