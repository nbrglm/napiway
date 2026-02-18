package main

import (
	"bytes"
	"context"
	"embed"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path"
	"strings"
	"testing"
	"time"
)

//go:embed testdata/spec.yaml
var testdata embed.FS

var serverInterface = "localhost:5000"
var httpServerURL = "http://" + serverInterface
var httpServerHealthURL = httpServerURL + "/health"

type ClientResult map[string]bool

func TestSpec(t *testing.T) {
	// TODO: FIX THE GENERATED SERVER NOT EXITING PROPERLY ON TEST FAILURE, CAUSING HANGS
	// AND SERVER PORT BEING OCCUPIED FOR SUBSEQUENT TESTS
	dir, err := os.Getwd()
	if err != nil {
		t.Logf("Failed to get current working directory: %v", err)
		os.Exit(1)
	}
	specPath := path.Join(dir, "testdata/spec.yaml")
	serverPath := path.Join(dir, "testdata/out/server")
	clientPath := path.Join(dir, "testdata/out/client")

	runTestForSpec(t, dir, specPath, serverPath, clientPath)
}

// func TestSpecLow(t *testing.T) {
// 	dir, err := os.Getwd()
// 	if err != nil {
// 		t.Fatalf("Failed to get current working directory: %v", err)
// 	}
// 	lowSpecPath := path.Join(dir, "testdata/spec-low.yaml")
// 	lowServerPath := path.Join(dir, "testdata/out/low/server")
// 	lowClientPath := path.Join(dir, "testdata/out/low/client")

// 	runTestForSpec(t, dir, lowSpecPath, lowServerPath, lowClientPath)
// }

func runTestForSpec(t *testing.T, workDir, specPath, serverPath, clientPath string) {
	generateTestServerHelpers(t, workDir, specPath)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Run the generated server in a separate process
	serverCmd := runGeneratedServer(t, ctx, serverPath)
	defer func() {
		// Stop the server process
		t.Logf("Killing server process PID: %v...", serverCmd.Process.Pid)
		err := serverCmd.Process.Kill()
		if err != nil {
			t.Logf("Failed to kill server process: %v", err)
		} else {
			t.Log("Server process killed successfully")
		}
		t.Log("Waiting for server process to exit...")
		err = serverCmd.Wait()
		if err != nil && !strings.Contains(err.Error(), "signal: killed") {
			t.Logf("Server process exited with error: %v", err)
		} else {
			t.Log("Server process exited successfully")
		}
	}()

	output, err := runCommandForOutput(t, clientPath, "go", "run", ".", httpServerURL)
	if err != nil {
		t.Logf("%s TEST FAILED: Failed to run generated client: %v", t.Name(), err)
		os.Exit(1)
	}

	clientResult := ClientResult{}

	if err := json.Unmarshal([]byte(output), &clientResult); err != nil {
		t.Logf("%s TEST FAILED: Failed to unmarshal client output:\n\tError: %v\n\tOutput: %s", t.Name(), err, output)
		os.Exit(1)
	}

	allPassed := true

	for test, pass := range clientResult {
		if !pass {
			allPassed = false
			t.Logf("%s TEST FAILED: Test %s did not pass", t.Name(), test)
		} else {
			t.Logf("%s TEST PASSED: Test %s passed", t.Name(), test)
		}
	}

	if allPassed {
		t.Logf("%s TEST PASSED", t.Name())
	} else {
		t.Logf("%s TEST FAILED: Some tests did not pass", t.Name())
		t.FailNow()
	}
}

func generateTestServerHelpers(t *testing.T, dir, specPath string) {
	t.Logf("Generating server helpers from spec: %s", specPath)

	if err := runCommand(t, dir, "go", "run", ".", "--config", specPath); err != nil {
		t.Logf("Failed to generate server helpers: %v", err)
		os.Exit(1)
	}
}

func runCommand(t *testing.T, workDir, name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdout = t.Output()
	cmd.Stderr = t.Output()
	cmd.Dir = workDir

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("Failed to run command %s %v: %v", name, args, err)
	}
	return nil
}

func runCommandForOutput(t *testing.T, workDir, name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = t.Output()
	cmd.Dir = workDir

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("Failed to run command %s %v: %v", name, args, err)
	}
	return out.String(), nil
}

func runGeneratedServer(t *testing.T, ctx context.Context, serverDir string) *exec.Cmd {
	buildCmd := exec.CommandContext(ctx, "go", "build", "-o", "out/server_binary", ".")
	buildCmd.Dir = serverDir
	buildCmd.Stdout = t.Output()
	buildCmd.Stderr = t.Output()

	t.Logf("Building generated server binary in directory: %s", serverDir)
	if err := buildCmd.Run(); err != nil {
		t.Logf("Failed to build generated server binary: %v", err)
		os.Exit(1)
	}

	cmd := exec.CommandContext(ctx, "./out/server_binary", serverInterface)
	cmd.Dir = serverDir
	cmd.Stdout = t.Output()
	cmd.Stderr = t.Output()

	t.Logf("Running generated server in directory: %s", serverDir)
	if err := cmd.Start(); err != nil {
		t.Logf("Failed to start generated server: %v", err)
		t.Fail()
		os.Exit(1)
	}

	// Give the server a moment to start

	waitForServerStartup(t, httpServerHealthURL)

	return cmd
}

func waitForServerStartup(t *testing.T, url string) {
	deadline := time.Now().Add(10 * time.Second)

	for time.Now().Before(deadline) {
		resp, err := http.Get(url)
		if err == nil && resp.StatusCode == http.StatusOK {
			t.Logf("Server is up and running at %s", url)
			return
		}
		time.Sleep(50 * time.Millisecond)
	}

	t.Logf("Server did not start within the expected time")
	os.Exit(1)
}
