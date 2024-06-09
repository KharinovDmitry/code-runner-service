package service

import (
	"code-runner-service/internal/domain/enum"
	"code-runner-service/internal/domain/executor"
)

type ExecutorFactory interface {
	NewExecutor(code string, memoryLimitKb int, timeLimitInMs int, language enum.Language) (executor.Executor, error)
}
