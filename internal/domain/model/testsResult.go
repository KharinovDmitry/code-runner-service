package model

import "code-runner-service/internal/domain/enum"

type TestsResult struct {
	ResultCode  enum.TestResultCode
	Description string
	Points      int
}
