package grpcServer

import (
	"code-runner-service/internal/domain/enum"
	"code-runner-service/internal/domain/model"
	"code-runner-service/internal/domain/service"
	"context"
	codeRunner "github.com/Reholly/kforge-proto/src/gen/code-runner"
	"google.golang.org/grpc"
)

type serverAPI struct {
	codeRunner.UnimplementedCodeRunnerServer
	testRunner service.TestRunner
}

func Register(gRPCServer *grpc.Server, testRunner service.TestRunner) {
	codeRunner.RegisterCodeRunnerServer(gRPCServer, &serverAPI{testRunner: testRunner})
}

func (s *serverAPI) RunTestsOnCode(ctx context.Context, in *codeRunner.RunTestsOnCodeRequest) (*codeRunner.RunTestsOnCodeResponse, error) {
	testsModel := make([]model.Test, len(in.Tests))
	for i, testProto := range in.Tests {
		testsModel[i] = model.Test{
			Input:          testProto.Input,
			ExpectedResult: testProto.Output,
			Points:         int(testProto.Points),
		}
	}

	res, err := s.testRunner.RunTest(ctx, enum.Language(in.Language), in.Code, int(in.MemoryLimitKb), int(in.TimeLimitMs), testsModel)
	return &codeRunner.RunTestsOnCodeResponse{
		ResultCode:  string(res.ResultCode),
		Points:      int32(res.Points),
		Description: res.Description,
	}, err
}
