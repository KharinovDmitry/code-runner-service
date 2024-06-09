package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
	"time"
)

func main() {
	ctx := context.Background()
	timeout, err := strconv.Atoi(os.Args[1])
	if err != nil {
		fmt.Println("INCORRECT TIME LIMIT: " + os.Args[1])
		return
	}

	memLimit, err := strconv.Atoi(os.Args[2])
	if err != nil {
		fmt.Println("INCORRECT MEMORY LIMIT" + os.Args[2])
		return
	}

	start := time.Now()
	ctx, closeCtx := context.WithTimeout(ctx, time.Duration(timeout)*time.Millisecond)
	defer closeCtx()

	runCommand := strings.Join(os.Args[3:], " ")
	args := []string{"chroot", "jail", "su", "unprivuser", "-c", runCommand}

	cmd := exec.CommandContext(ctx, args[0], args[1:]...)

	_ = int64(memLimit)

	stdin, err := cmd.StdinPipe()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer stdin.Close()

	var input string
	fmt.Fscan(os.Stdin, &input)
	fmt.Fprintln(stdin, input)

	outputBytes, err := cmd.CombinedOutput()
	outputString := string(outputBytes)
	end := time.Now()

	if exitErr, ok := err.(*exec.ExitError); ok {
		if status, ok := exitErr.Sys().(syscall.WaitStatus); ok {
			if end.Sub(start) >= time.Duration(timeout)*time.Millisecond {
				fmt.Println("TIME LIMIT")
				return
			}
			if status.Signaled() && status.Signal() == syscall.SIGKILL {
				fmt.Println("MEMORY LIMIT")
				return
			}
		}

		fmt.Printf("RUNTIME ERROR: %s, %s\n", outputString, err.Error())
		return
	}

	if err != nil {
		fmt.Println(err.Error())
	}

	fmt.Println(outputString)
}
