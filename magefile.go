//go:build mage

package main

import (
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

// Build runs a full CI build.
func Build() {
	mg.Deps(Generate, Lint, Test)
}

// Clean removes all build artifacts.
func Clean() error {
	return sh.Run("rm", "-rf", "build")
}

// Lint runs the Go linter.
func Lint() error {
	return sh.RunV("go", "tool", "golangci-lint", "run")
}

// Test runs the Go tests.
func Test() error {
	mg.Deps(Generate)
	return sh.RunV("go", "test", "-v", "./...")
}

// Generate runs all code generators.
func Generate() error {
	return sh.RunV("go", "generate", "./...")
}

// DockerPush pushes the demo server Docker image to the registry.
func DockerPush() error {
	return sh.RunWith(
		map[string]string{
			"KO_DOCKER_REPO": "ghcr.io/way-platform/ileap-go",
		},
		"go", "tool", "ko", "build",
		"--base-import-paths",
		"--platform", "linux/amd64",
		"./cmd/demo-server",
	)
}

// DockerBuildLocal builds the demo server Docker image locally.
func DockerBuild() error {
	return sh.RunWith(
		map[string]string{
			"KO_DOCKER_REPO": "ko.local",
		},
		"go", "tool", "ko", "build",
		"--base-import-paths",
		"--platform", "linux/amd64",
		"./cmd/demo-server",
	)
}
