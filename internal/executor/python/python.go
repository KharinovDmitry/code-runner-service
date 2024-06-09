package python

import (
	"code-runner-service/internal/domain/executor"
	executor2 "code-runner-service/internal/domain/executor"
	"code-runner-service/lib/byteconv"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

var baseContainerMemoryKb = 5 * 1024

type PythonExecutor struct {
	code            string
	fileName        string
	memoryLimitInKb int
	timeLimitInMs   int
}

func NewPythonExecutor(code string, memoryLimitInKb int, timeLimitInMs int) executor.Executor {
	return &PythonExecutor{
		code:            code,
		memoryLimitInKb: memoryLimitInKb,
		timeLimitInMs:   timeLimitInMs,
	}
}

func (p *PythonExecutor) Init() error {
	fileName := strconv.FormatInt(time.Now().Unix(), 10) + ".py"
	file, err := os.Create("tmp/" + fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	if _, err = file.WriteString(p.code); err != nil {
		return err
	}

	p.fileName = fileName
	cmd := exec.Command("docker", "run",
		"--mount", "type=bind,source=./tmp,target=/home/jail/tmp",
		"--rm",
		"--name", fileName,
		fmt.Sprintf("--memory=%dk", p.memoryLimitInKb+baseContainerMemoryKb),
		"-d",
		"-e", "FILE_NAME="+fileName,
		"python_executor",
		"sleep", "infinity")

	err = cmd.Run()
	if err != nil {
		return err
	}

	return nil
}

func (p *PythonExecutor) Execute(ctx context.Context, input string) (output string, err error) {
	cmd := exec.Command("docker", "exec",
		"-i", p.fileName,
		"./unprivrun", strconv.Itoa(p.timeLimitInMs), strconv.Itoa(p.memoryLimitInKb),
		"python3", "tmp/"+p.fileName)

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return "", fmt.Errorf("In PythonExecutor(Execute): %w", err)
	}
	defer stdin.Close()

	fmt.Fprintln(stdin, input)

	outputBytes, err := cmd.Output()
	outputString := strings.Replace(byteconv.String(outputBytes), "\n", "", -1)
	if err != nil {
		return outputString, fmt.Errorf("In PythonExecutor(Execute): %w", err)
	}

	if outputString == "TIME LIMIT" {
		return "", executor2.TimeLimitError
	}

	if outputString == "MEMORY LIMIT" {
		return "", executor2.MemoryLimitError
	}

	if strings.Contains(outputString, "RUNTIME ERROR") {
		parts := strings.SplitN(outputString, ":", 2)
		if len(parts) != 2 {
			return outputString, errors.New("INCORRECT RESULT")
		}

		return parts[1], executor2.RuntimeError
	}

	return outputString, nil
}

func (p *PythonExecutor) Close() error {
	cmd := exec.Command("docker", "kill", p.fileName)
	err := cmd.Run()
	if err != nil {
		return err
	}

	return os.Remove("tmp/" + p.fileName)
}
