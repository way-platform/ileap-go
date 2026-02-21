//go:build mage

package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"

	"github.com/magefile/mage/mg"
)

var Default = Build

// Build runs a full CI build.
func Build() {
	mg.SerialDeps(
		Download,
		Lint,
		Generate,
		Test,
		Tidy,
		CLI,
		Diff,
	)
}

// Download downloads the Go dependencies.
func Download() error {
	log.Println("downloading dependencies")
	return forEachGoMod(func(dir string) error {
		return cmd(dir, "go", "mod", "download").Run()
	})
}

// Lint runs the Go linter and fixes code style issues.
func Lint() error {
	log.Println("linting and fixing code")
	return forEachGoMod(func(dir string) error {
		return tool(
			dir,
			"golangci-lint",
			"run",
			"--fix",
			"--path-prefix",
			dir,
			"--build-tags",
			"mage",
		).Run()
	})
}

// Generate runs all code generators.
func Generate() error {
	log.Println("generating code")
	return forEachGoMod(func(dir string) error {
		return cmd(dir, "go", "generate", "-v", "./...").Run()
	})
}

// Test runs the Go tests.
func Test() error {
	log.Println("running tests")
	return cmd(root(), "go", "test", "-v", "-cover", "./...").Run()
}

// Tidy tidies the Go mod files.
func Tidy() error {
	log.Println("tidying Go mod files")
	return forEachGoMod(func(dir string) error {
		return cmd(dir, "go", "mod", "tidy", "-v").Run()
	})
}

// CLI builds the CLI.
func CLI() error {
	log.Println("building CLI")
	return cmd(root("cmd/ileap"), "go", "install", ".").Run()
}

// VHS records the CLI GIF using VHS.
func VHS() error {
	log.Println("recording CLI GIF")
	mg.Deps(CLI)
	return tool(root("docs"), "vhs", "cli.tape").Run()
}

// Diff checks for git diffs.
func Diff() error {
	log.Println("checking for git diffs")
	return cmd(root(), "git", "diff", "--exit-code").Run()
}

// ACT runs the ACT conformance test suite against the demo server.
// Pass a base URL to test against a remote server, or omit to start
// a local server.
func ACT(baseURL string) error {
	log.Println("running ACT conformance tests")
	// Install ACT binary.
	actBin, err := installACT()
	if err != nil {
		return err
	}
	if baseURL == "" {
		// Start a local server.
		var cleanup func()
		baseURL, cleanup, err = startLocalServer()
		if err != nil {
			return err
		}
		defer cleanup()
	}
	// Run ACT.
	return cmd(
		root(), actBin,
		"test", "-b", baseURL, "-u", "hello", "-p", "pathfinder",
	).Run()
}

func startLocalServer() (baseURL string, cleanup func(), err error) {
	// Find a free port.
	lis, err := net.Listen("tcp", ":0")
	if err != nil {
		return "", nil, fmt.Errorf("find free port: %w", err)
	}
	port := lis.Addr().(*net.TCPAddr).Port
	lis.Close()
	baseURL = fmt.Sprintf("http://localhost:%d", port)
	// Build the ileap binary.
	binary := filepath.Join(root(), "ileap-tmp")
	if err := cmd(root("cmd/ileap"), "go", "build", "-o", binary, ".").Run(); err != nil {
		return "", nil, fmt.Errorf("build ileap: %w", err)
	}
	// Start: ileap-tmp demo-server --port <port> --base-url <url>
	server := cmdWith(nil, root(), binary, "demo-server",
		"--port", fmt.Sprintf("%d", port),
		"--base-url", baseURL,
	)
	if err := server.Start(); err != nil {
		os.Remove(binary)
		return "", nil, fmt.Errorf("start ileap demo-server: %w", err)
	}
	// Wait for the server to be ready.
	if err := waitForServer(baseURL, 10*time.Second); err != nil {
		_ = server.Process.Kill()
		_ = server.Wait()
		os.Remove(binary)
		return "", nil, fmt.Errorf("wait for ileap demo-server: %w", err)
	}
	cleanup = func() {
		_ = server.Process.Kill()
		_ = server.Wait()
		os.Remove(binary)
	}
	return baseURL, cleanup, nil
}

// DockerPush pushes the ileap Docker image to GHCR.
func DockerPush() error {
	log.Println("pushing ileap Docker image to GHCR")
	c := tool(root("cmd", "ileap"), "ko", "build",
		"--base-import-paths",
		"--platform", "linux/amd64",
		".",
	)
	c.Env = append(os.Environ(), "KO_DOCKER_REPO=ghcr.io/way-platform/ileap-go")
	return c.Run()
}

// DockerBuild builds the ileap Docker image locally.
func DockerBuild() error {
	log.Println("building ileap Docker image locally")
	c := tool(root("cmd", "ileap"), "ko", "build",
		"--base-import-paths",
		"--platform", "linux/amd64",
		".",
	)
	c.Env = append(os.Environ(), "KO_DOCKER_REPO=ko.local")
	return c.Run()
}

// installACT downloads the ACT conformance binary and caches it.
func installACT() (string, error) {
	arch := runtime.GOARCH
	switch arch {
	case "amd64":
		arch = "x86_64"
	case "arm64":
		// already correct
	default:
		return "", fmt.Errorf("unsupported architecture: %s", arch)
	}
	binPath := root("tools", "build", "act", "conformance_"+arch)
	if _, err := os.Stat(binPath); err == nil {
		return binPath, nil
	}
	url := fmt.Sprintf(
		"https://actbin.blob.core.windows.net/act-bin/conformance_%s",
		arch,
	)
	log.Printf("downloading ACT binary from %s", url)
	if err := downloadBinary(url, binPath); err != nil {
		return "", fmt.Errorf("download ACT: %w", err)
	}
	return binPath, nil
}

func downloadBinary(url, dst string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP %s", resp.Status)
	}
	if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
		return err
	}
	f, err := os.OpenFile(dst, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o755)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = io.Copy(f, resp.Body)
	return err
}

func waitForServer(baseURL string, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		resp, err := http.Get(baseURL + "/.well-known/openid-configuration")
		if err == nil {
			resp.Body.Close()
			return nil
		}
		time.Sleep(100 * time.Millisecond)
	}
	return fmt.Errorf("server not ready after %v", timeout)
}

// Helpers

func forEachGoMod(f func(dir string) error) error {
	return filepath.WalkDir(root(), func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || d.Name() != "go.mod" {
			return nil
		}
		return f(filepath.Dir(path))
	})
}

// root returns the absolute path to the project root.
func root(subdirs ...string) string {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		panic("unable to get current file path")
	}
	rootDir := filepath.Dir(filepath.Dir(filename))
	return filepath.Join(append([]string{rootDir}, subdirs...)...)
}

// cmd runs a command in a specific directory.
func cmd(dir string, command string, args ...string) *exec.Cmd {
	return cmdWith(nil, dir, command, args...)
}

// cmdWith runs a command with environment variables.
func cmdWith(env map[string]string, dir string, command string, args ...string) *exec.Cmd {
	c := exec.Command(command, args...)
	c.Env = os.Environ()
	for key, value := range env {
		c.Env = append(c.Env, fmt.Sprintf("%s=%s", key, value))
	}
	c.Dir = dir
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	return c
}

// tool runs a go tool command using tools/go.mod.
func tool(dir string, toolName string, args ...string) *exec.Cmd {
	return toolWith(nil, dir, toolName, args...)
}

// toolWith runs a go tool command with environment variables.
func toolWith(env map[string]string, dir string, toolName string, args ...string) *exec.Cmd {
	cmdArgs := []string{"tool", "-modfile", filepath.Join(root(), "tools", "go.mod"), toolName}
	return cmdWith(env, dir, "go", append(cmdArgs, args...)...)
}
