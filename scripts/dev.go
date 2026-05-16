//go:build ignore

package main

import (
	"log"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
)

func main() {
	vite := exec.Command(npm(), "run", "dev")
	vite.Dir = "frontend"
	vite.Stdout = os.Stdout
	vite.Stderr = os.Stderr
	if err := vite.Start(); err != nil {
		log.Fatalf("failed to start vite: %v", err)
	}

	config := ".air.unix.toml"
	if runtime.GOOS == "windows" {
		config = ".air.toml"
	}
	air := exec.Command("air", "-c", config)
	air.Stdout = os.Stdout
	air.Stderr = os.Stderr
	air.Stdin = os.Stdin
	if err := air.Start(); err != nil {
		vite.Process.Kill()
		log.Fatalf("failed to start air (is it installed? go install github.com/air-verse/air@latest): %v", err)
	}

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)
	<-sig

	vite.Process.Kill()
	air.Process.Kill()
}

func npm() string {
	if runtime.GOOS == "windows" {
		return "npm.cmd"
	}
	return "npm"
}
