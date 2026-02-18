//go:build mage

package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

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
