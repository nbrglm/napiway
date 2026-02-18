package utils

import (
	"cmp"
	"fmt"
	"maps"
	"os"
	"os/exec"
	"slices"
	"strings"
)

func SortedMapKeys[V any, K cmp.Ordered](m map[K]V) []K {
	keys := slices.Collect(maps.Keys(m))
	slices.Sort(keys)
	return keys
}

func ExecCommand(command string, workDir string) error {
	// split the command into name and args
	parts := strings.Fields(command)
	name := parts[0]
	args := parts[1:]

	cmd := exec.Command(name, args...)
	if workDir == "" {
		wd, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("failed to get working directory: %w", err)
		}
		workDir = wd
	}
	cmd.Dir = workDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to run command %s: %w", command, err)
	}
	return nil
}

func ClearOutputDir(dir string) error {
	err := os.RemoveAll(dir)
	if err != nil {
		return fmt.Errorf("failed to remove output directory: %w", err)
	}

	// recreate the directory
	err = os.MkdirAll(dir, 0o755)
	if err != nil {
		return fmt.Errorf("failed to recreate output directory: %w", err)
	}
	return nil
}

func WriteFile(path string, content []byte) error {
	return os.WriteFile(path, content, 0o644)
}
