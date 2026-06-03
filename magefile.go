//go:build mage

package main

import (
	"fmt"
	"os"
	"runtime"

	"github.com/magefile/mage/sh"
)

const tailwindVersion = "v3.4.17"

// InstallTailwind downloads the Tailwind standalone CLI for the current platform
func InstallTailwind() error {
	binary := tailwindBinaryPath()
	if _, err := os.Stat(binary); err == nil {
		fmt.Println("Tailwind already installed, skipping.")
		return nil
	}

	url := tailwindDownloadURL()
	fmt.Printf("Downloading Tailwind %s from %s\n", tailwindVersion, url)

	if err := sh.Run("curl", "-sLo", binary, url); err != nil {
		return err
	}
	return sh.Run("chmod", "+x", binary)
}

// BuildCSS compiles Tailwind CSS
func BuildCSS() error {
	return sh.Run(
		tailwindBinaryPath(),
		"-c", "./tailwind/tailwind.config.js",
		"-i", "./tailwind/input.css",
		"-o", "./web/static/css/site.css",
		"--minify",
	)
}

// GenerateTempl runs templ generate
func GenerateTempl() error {
	return sh.Run("templ", "generate")
}

// BuildGo compiles the Go binary
func BuildGo() error {
	if err := GenerateTempl(); err != nil {
		return err
	}
	return sh.Run("go", "build", "-o", "./bin/server", "./cmd/server")
}

// Seed creates an admin user. Usage: mage seed admin@example.com
func Seed(email string) error {
	return sh.Run("go", "run", "./cmd/seed", email)
}

// Build runs a full production build
func Build() error {
	if err := BuildCSS(); err != nil {
		return err
	}
	return BuildGo()
}

// Dev generates templ, builds CSS, and starts the server
func Dev() error {
	if err := GenerateTempl(); err != nil {
		return err
	}
	if err := BuildCSS(); err != nil {
		return err
	}
	return sh.Run("go", "run", "./cmd/server/")
}

// Start builds everything (CSS + templ + Go) and starts the server
func Start() error {
	if err := Build(); err != nil {
		return err
	}
	return sh.Run("./bin/server")
}

// Run starts the server
func Run() error {
	return sh.Run("./bin/server")
}

func tailwindBinaryPath() string {
	if runtime.GOOS == "windows" {
		return "./tailwind/tailwindcss.exe"
	}
	return "./tailwind/tailwindcss"
}

func tailwindDownloadURL() string {
	os := runtime.GOOS
	arch := runtime.GOARCH

	osName := map[string]string{
		"darwin":  "macos",
		"linux":   "linux",
		"windows": "windows",
	}[os]

	archName := map[string]string{
		"amd64": "x64",
		"arm64": "arm64",
	}[arch]

	ext := ""
	if os == "windows" {
		ext = ".exe"
	}

	return fmt.Sprintf(
		"https://github.com/tailwindlabs/tailwindcss/releases/download/%s/tailwindcss-%s-%s%s",
		tailwindVersion, osName, archName, ext,
	)
}