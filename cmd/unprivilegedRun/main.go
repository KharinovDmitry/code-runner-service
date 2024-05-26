package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strconv"
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

	_, err = strconv.Atoi(os.Args[2])
	if err != nil {
		fmt.Println("INCORRECT MEMORY LIMIT" + os.Args[2])
		return
	}

	ctx, closeCtx := context.WithTimeout(ctx, time.Duration(timeout)*time.Millisecond)
	defer closeCtx()

	args := []string{"jail"}
	for _, arg := range os.Args[3:] {
		args = append(args, arg)
	}

	cmd := exec.CommandContext(ctx, "chroot", args...)
	cmd.SysProcAttr = &syscall.SysProcAttr{}

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

	if err != nil && errors.Is(err, context.DeadlineExceeded) {
		fmt.Println("TIME LIMIT")
		return
	}

	if exitErr, ok := err.(*exec.ExitError); ok {
		if status, ok := exitErr.Sys().(syscall.WaitStatus); ok && status.Signaled() && status.Signal() == syscall.SIGKILL {
			fmt.Println("TIME LIMIT")
			return
		}

		fmt.Printf("RUNTIME ERROR: %s, %s\n", outputString, err.Error())
		return
	}

	if err != nil {
		fmt.Println(err.Error())
	}

	fmt.Println(outputString)
}
