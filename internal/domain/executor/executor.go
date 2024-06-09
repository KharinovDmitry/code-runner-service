package executor

import (
	"context"
	"errors"
)

var (
	ErrUnknownLanguage = errors.New("unknown language")
	CompileError       = errors.New("compile error")
	MemoryLimitError   = errors.New("memory limit")
	TimeLimitError     = errors.New("time limit")
	RuntimeError       = errors.New("runtime error")
)

type Executor interface {
	Execute(ctx context.Context, input string) (output string, err error)
	Init() error
	Close() error
}
