package service

import (
	"code-runner-service/internal/domain/enum"
	"code-runner-service/internal/domain/executor"
)

type ExecutorFactory interface {
	NewExecutor(code string, language enum.Language) (executor.Executor, error)
}
