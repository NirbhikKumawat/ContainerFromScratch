package main

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

func run() {
	fmt.Printf("Running command %v as PID %d\n", os.Args[1:], os.Getpid())

	cmdArgs := os.Args[2:]
	cmd := exec.Command(cmdArgs[0], cmdArgs[1:]...)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWPID, // crates a new PID namespace
	}
	const rootfs = "./my-rootfs"

	if err := syscall.Chroot(rootfs); err != nil { // change the root filesystem
		fmt.Printf("Failed to chroot: %v\n", err)
		os.Exit(1)
	}

	if err := syscall.Chdir("/"); err != nil { //change working dir to "/" inside the root
		fmt.Printf("Failed to chdir: %v\n", err)
		os.Exit(1)
	}

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
