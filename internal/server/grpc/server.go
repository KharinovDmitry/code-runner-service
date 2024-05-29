package grpc

import (
	"code-runner-service/internal/domain/enum"
	"code-runner-service/internal/domain/model"
	"code-runner-service/internal/domain/service"
	codeRunnerV1 "code-runner-service/proto/gen/go/code-runner"
	"context"
	"google.golang.org/grpc"
)

type serverAPI struct {
	codeRunnerV1.UnimplementedCodeRunnerServer
	testRunner service.TestRunner
}

func Register(gRPCServer *grpc.Server, testRunner service.TestRunner) {
	codeRunnerV1.RegisterCodeRunnerServer(gRPCServer, &serverAPI{testRunner: testRunner})
}

func (s *serverAPI) RunTestsOnCode(ctx context.Context, in *codeRunnerV1.RunTestsOnCodeRequest) (*codeRunnerV1.RunTestsOnCodeResponse, error) {
	testsModel := make([]model.Test, len(in.Tests))
	for i, testProto := range in.Tests {
		testsModel[i] = model.Test{
			Input:          testProto.Input,
			ExpectedResult: testProto.Output,
			Points:         int(testProto.Points),
		}
	}

	res, err := s.testRunner.RunTest(ctx, enum.Language(in.Language), in.Code, int(in.MemoryLimitKb), int(in.TimeLimitMs), testsModel)
	return &codeRunnerV1.RunTestsOnCodeResponse{
		ResultCode:  string(res.ResultCode),
		Points:      int32(res.Points),
		Description: res.Description,
	}, err
}
