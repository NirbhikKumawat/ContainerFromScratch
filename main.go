package main

import (
	"fmt"
	"os"
	"os/exec"
)

func run() {
	fmt.Printf("Running command %v as PID %d\n", os.Args[1:], os.Getpid())

	cmdArgs := os.Args[2:]

	cmd := exec.Command(cmdArgs[0], cmdArgs[1:]...)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Printf("Error running command: %v\n", err)
		os.Exit(1)
	}
}
func main() {
	if len(os.Args) > 1 && os.Args[1] == "run" {
		run()
	} else {
		fmt.Println("Usage: go run main.go run <command> [arguments]")
	}
}
