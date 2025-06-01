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
