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

	dlv := exec.Command("dlv", "debug", "--headless", "--listen=:2345",
		"--api-version=2", "--accept-multiclient", ".")
	dlv.Stdout = os.Stdout
	dlv.Stderr = os.Stderr
	dlv.Stdin = os.Stdin
	if err := dlv.Start(); err != nil {
		vite.Process.Kill()
		log.Fatalf("failed to start dlv: %v", err)
	}

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)
	<-sig

	vite.Process.Kill()
	dlv.Process.Kill()
}

func npm() string {
	if runtime.GOOS == "windows" {
		return "npm.cmd"
	}
	return "npm"
}
