package service

import (
	"code-runner-service/internal/domain/enum"
	"code-runner-service/internal/domain/model"
	"context"
	"errors"
)

var (
	ErrTaskNotFound  = errors.New("task not found")
	ErrTestsNotFound = errors.New("tests not found")
)

type TestRunner interface {
	RunTest(ctx context.Context, language enum.Language, code string, memoryLimitInKB, timeLimitInMs int, tests []model.Test) (model.TestsResult, error)
}
