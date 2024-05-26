package tests

import (
	"code-runner-service/internal/domain/enum"
	mockExecutor2 "code-runner-service/internal/domain/executor/mocks"
	"code-runner-service/internal/domain/model"
	mockService "code-runner-service/internal/domain/service/mocks"

	"code-runner-service/internal/service/implementations"
	"context"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func initMockExecutorFactory(cntrl *gomock.Controller) *mockService.MockExecutorFactory {
	mockFactory := mockService.NewMockExecutorFactory(cntrl)
	mockExecutor := mockExecutor2.NewMockExecutor(cntrl)
	mockExecutor.EXPECT().Execute(gomock.Any(), "1", gomock.Any(), gomock.Any()).Return("1", nil)
	mockExecutor.EXPECT().Execute(gomock.Any(), "2", gomock.Any(), gomock.Any()).Return("4", nil)
	mockExecutor.EXPECT().Execute(gomock.Any(), "3", gomock.Any(), gomock.Any()).Return("9", nil)

	mockExecutor.EXPECT().Close()

	mockFactory.EXPECT().NewExecutor(gomock.Any(), gomock.Any()).Return(mockExecutor, nil)

	return mockFactory
}

func Test(t *testing.T) {
	cntrl := gomock.NewController(t)
	executorFactory := initMockExecutorFactory(cntrl)

	service := implementations.NewTestRunnerService(executorFactory)

	tests := []model.Test{
		{
			ID:             1,
			Input:          "1",
			ExpectedResult: "1",
			Points:         1,
		},
		{
			ID:             2,
			Input:          "2",
			ExpectedResult: "4",
			Points:         1,
		},
		{
			ID:             3,
			Input:          "3",
			ExpectedResult: "9",
			Points:         1,
		},
	}
	actual, err := service.RunTest(context.Background(), enum.CPP, "", 10, 0, tests)
	assert.Nil(t, err)
	expected := model.TestsResult{
		ResultCode:  "SC",
		Description: "",
		Points:      3,
	}
	assert.Equal(t, expected, actual)
}
