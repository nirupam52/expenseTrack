//go:build ignore

package main

import (
	"log"
	"os"
	"os/exec"
	"runtime"
)

func main() {
	vite := exec.Command(npm(), "run", "build")
	vite.Dir = "frontend"
	vite.Stdout = os.Stdout
	vite.Stderr = os.Stderr
	if err := vite.Run(); err != nil {
		log.Fatalf("frontend build failed: %v", err)
	}

	bin := "expensetrack"
	if runtime.GOOS == "windows" {
		bin = "expensetrack.exe"
	}
	build := exec.Command("go", "build", "-o", bin, ".")
	build.Stdout = os.Stdout
	build.Stderr = os.Stderr
	if err := build.Run(); err != nil {
		log.Fatalf("go build failed: %v", err)
	}

	run := exec.Command("./" + bin)
	if runtime.GOOS == "windows" {
		run = exec.Command(bin)
	}
	run.Stdout = os.Stdout
	run.Stderr = os.Stderr
	run.Stdin = os.Stdin
	if err := run.Run(); err != nil {
		log.Fatalf("server exited: %v", err)
	}
}

func npm() string {
	if runtime.GOOS == "windows" {
		return "npm.cmd"
	}
	return "npm"
}
