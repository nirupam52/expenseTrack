//go:build ignore

package main

import (
	"os"
	"os/exec"
	"runtime"
)

func main() {
	// Pick the right Air config for the current OS
	config := ".air.unix.toml"
	if runtime.GOOS == "windows" {
		config = ".air.toml"
	}

	cmd := exec.Command("air", "-c", config)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Run()
}
