package implementations

import (
	"code-runner-service/internal/domain/enum"
	executor2 "code-runner-service/internal/domain/executor"
	"code-runner-service/internal/domain/model"
	"code-runner-service/internal/domain/service"
	"context"
	"errors"
	"fmt"
	"strings"
)

type TestRunnerService struct {
	codeRunnerFactory service.ExecutorFactory
}

func NewTestRunnerService(codeRunnerFactory service.ExecutorFactory) *TestRunnerService {
	return &TestRunnerService{
		codeRunnerFactory: codeRunnerFactory,
	}
}

func (s *TestRunnerService) RunTest(ctx context.Context, language enum.Language, code string, memoryLimitInKB, timeLimitInMs int, tests []model.Test) (model.TestsResult, error) {
	program, err := s.codeRunnerFactory.NewExecutor(code, language)
	if err != nil {
		if errors.Is(err, executor2.CompileError) {
			res := model.TestsResult{
				ResultCode:  enum.CompileErrorCode,
				Description: err.Error(),
				Points:      0,
			}
			return res, nil
		}
		return model.TestsResult{}, fmt.Errorf("In TestRunnerService(RunTest): %w", err)
	}
	defer program.Close()

	res := model.TestsResult{
		ResultCode:  enum.SuccessCode,
		Description: "",
		Points:      0,
	}
	for i, test := range tests {
		output, err := program.Execute(ctx, test.Input, memoryLimitInKB, timeLimitInMs)
		if err != nil {
			if errors.Is(err, executor2.TimeLimitError) {
				res = model.TestsResult{
					ResultCode:  enum.TimeLimitCode,
					Description: fmt.Sprintf("Test failed: %d. Time limit error", i),
					Points:      res.Points,
				}
				break
			}
			if errors.Is(err, executor2.RuntimeError) {
				res = model.TestsResult{
					ResultCode:  enum.RuntimeErrorCode,
					Description: fmt.Sprintf("Test failed: %d. Description: %s Output: %s", i, err.Error(), output),
					Points:      res.Points,
				}
				break
			}
			return model.TestsResult{}, fmt.Errorf("In TestService(RunTests): %w", err)
		}
		if strings.Replace(output, "\n", "", 1) != test.ExpectedResult {
			res = model.TestsResult{
				ResultCode:  enum.IncorrectAnswerCode,
				Description: fmt.Sprintf("Test failed: %d Excepted: %s Received: %s", i, test.ExpectedResult, output),
				Points:      res.Points,
			}
			break
		}
		res.Points += test.Points
	}

	return res, nil
}
