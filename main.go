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
		Cloneflags: syscall.CLONE_NEWPID | syscall.CLONE_NEWNS, // creates a new PID namespace,also added a new mount namespace
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

	if err := syscall.Mount("proc", "proc", "proc", 0, ""); err != nil { // mount /proc filesystem
		fmt.Printf("Error mounting /proc: %v\n", err)
		os.Exit(1)
	}

	if err := cmd.Run(); err != nil {
		fmt.Printf("Error running command: %v\n", err)
		os.Exit(1)
	}

	if err := syscall.Unmount("/proc", 0); err != nil { // unmount /proc after the commands exit
		fmt.Printf("Error unmounting /proc: %v\n", err)
	}
}
func runChild() {
	fmt.Printf("Running as child (PID %d) inside new namespace\n", os.Getpid())

	cmdArgs := os.Args[2:]

	const rootfs = "./my-rootfs"

	if err := syscall.Chroot(rootfs); err != nil {
		fmt.Printf("Failed to chroot: %v\n", err)
		os.Exit(1)
	}

	if err := syscall.Chdir("/"); err != nil {
		fmt.Printf("Failed to chdir: %v\n", err)
		os.Exit(1)
	}

	if err := syscall.Mount("proc", "proc", "proc", 0, ""); err != nil {
		fmt.Printf("Error mounting /proc: %v\n", err)
		os.Exit(1)
	}

	env := os.Environ()

	if err := syscall.Exec(cmdArgs[0], cmdArgs[1:], env); err != nil {
		fmt.Printf("Error executing command: %v\n", err)
		os.Exit(1)
	}
}

func runParent() {
	childArgs := append([]string{"child"}, os.Args[2:]...)

	cmd := exec.Command("/proc/self/exe", childArgs...)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWPID | syscall.CLONE_NEWNS,
	}

	fmt.Printf("Running as parent (PID %d), launching child\n", os.Getpid())

	if err := cmd.Run(); err != nil {
		fmt.Printf("Error running parent: %v\n", err)
		os.Exit(1)
	}
}
func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: ./main run <command>")
		os.Exit(1)
	}
	switch os.Args[1] {
	case "run":
		runParent()
	case "child":
		runChild()
	default:
		fmt.Println("Usage: ./main run <command>")
	}
}
