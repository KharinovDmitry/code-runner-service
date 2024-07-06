package cpp

import (
	"code-runner-service/internal/domain/executor"
	"code-runner-service/lib/byteconv"
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

var baseContainerMemoryKb = 6 * 1024

type CPPExecutor struct {
	code            string
	fileName        string
	memoryLimitInKb int
	timeLimitInMs   int
}

func NewCPPExecutor(code string, memoryLimitInKb int, timeLimitInMs int) executor.Executor {
	return &CPPExecutor{
		code:            code,
		memoryLimitInKb: memoryLimitInKb,
		timeLimitInMs:   timeLimitInMs,
	}
}

func (p *CPPExecutor) Init() error {
	id, err := uuid.NewRandom()
	if err != nil {
		return err
	}
	p.fileName = id.String() + ".cpp"

	_, err = os.Stat("tmp")
	if err != nil {
		if os.IsNotExist(err) {
			err = os.Mkdir("tmp", 0777)
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}

	file, err := os.Create("tmp/" + p.fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	if _, err = file.WriteString(p.code); err != nil {
		return err
	}

	cmd := exec.Command("docker", "run",
		"--rm",
		"--name", p.fileName,
		fmt.Sprintf("--memory=%dk", p.memoryLimitInKb+baseContainerMemoryKb),
		"-d",
		"-e", "FILE_NAME="+p.fileName,
		"cpp_executor",
		"sleep", "infinity")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("in CPPExecutor(Init, run): %w, %s", err, string(out))
	}

	cmd = exec.Command("g++", "-static", "-o", "./tmp/"+p.fileName+".out", "tmp/"+p.fileName)
	out, err = cmd.CombinedOutput()
	if err != nil {
		return errors.Wrap(executor.CompileError, string(out))
	}

	cmd = exec.Command("docker", "cp", "./tmp/"+p.fileName+".out", p.fileName+":./home/jail/tmp/")
	out, err = cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("in CPPExecutor(Init, copy): %s", err.Error()+" "+string(out))
	}

	return nil
}

func (p *CPPExecutor) Execute(ctx context.Context, input string) (output string, err error) {
	cmd := exec.Command("docker", "exec",
		"-i", p.fileName,
		"./unprivrun", strconv.Itoa(p.timeLimitInMs), strconv.Itoa(p.memoryLimitInKb),
		"./tmp/"+p.fileName+".out")

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return "", fmt.Errorf("In CPPExecutor(Execute): %w", err)
	}
	defer stdin.Close()

	fmt.Fprintln(stdin, input)

	outputBytes, err := cmd.Output()
	outputString := strings.Replace(byteconv.String(outputBytes), "\n", "", -1)
	if err != nil {
		return outputString, fmt.Errorf("In CPPExecutor(Execute): %w", err)
	}

	if outputString == "TIME LIMIT" {
		return "", executor.TimeLimitError
	}

	if outputString == "MEMORY LIMIT" {
		return "", executor.MemoryLimitError
	}

	if strings.Contains(outputString, "RUNTIME ERROR") {
		parts := strings.SplitN(outputString, ":", 2)
		if len(parts) != 2 {
			return outputString, errors.New("INCORRECT RESULT")
		}

		return parts[1], executor.RuntimeError
	}

	return outputString, nil
}

func (p *CPPExecutor) Close() error {
	cmd := exec.Command("docker", "kill", p.fileName)
	err := cmd.Run()
	if err != nil {
		return err
	}

	_ = os.Remove("tmp/" + p.fileName + ".out")
	return os.Remove("tmp/" + p.fileName)
}
