package main

import (
	"context"
	"errors"
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

	_, err = strconv.Atoi(os.Args[2])
	if err != nil {
		fmt.Println("INCORRECT MEMORY LIMIT" + os.Args[2])
		return
	}

	ctx, closeCtx := context.WithTimeout(ctx, time.Duration(timeout)*time.Millisecond)
	defer closeCtx()
	cmd := exec.CommandContext(ctx, "chroot", "jail", "su", "-", "unprivuser")

	stdin, err := cmd.StdinPipe()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer stdin.Close()

	fmt.Fprintln(stdin, strings.Join(os.Args[3:], " "))

	var input string
	fmt.Fscan(os.Stdin, &input)
	fmt.Fprintln(stdin, input)

	outputBytes, err := cmd.CombinedOutput()
	outputString := string(outputBytes)

	/*
		ctx, closeCtx := context.WithTimeout(ctx, time.Duration(timeout)*time.Millisecond)
		defer closeCtx()

		cmd := exec.CommandContext(ctx, os.Args[3], os.Args[4:]...)
		cmd.SysProcAttr = &syscall.SysProcAttr{}

	*/

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
