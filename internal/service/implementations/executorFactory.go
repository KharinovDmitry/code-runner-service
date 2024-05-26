package implementations

import (
	"code-runner-service/internal/domain/enum"
	"code-runner-service/internal/domain/executor"
	executor2 "code-runner-service/internal/domain/executor"
	"code-runner-service/internal/executor/python"
)

type ExecutorFactory struct {
	languageExecutorMap map[enum.Language]func(fileName string) executor.Executor
}

func NewExecutorFactory() *ExecutorFactory {
	return &ExecutorFactory{
		languageExecutorMap: map[enum.Language]func(code string) executor.Executor{
			enum.Python: python.NewPythonExecutor,
		},
	}
}

func (c *ExecutorFactory) NewExecutor(code string, language enum.Language) (executor.Executor, error) {
	if constructor, exist := c.languageExecutorMap[language]; exist {
		return constructor(code), nil
	}

	return nil, executor2.ErrUnknownLanguage
}
