package suite

import (
	"code-runner-service/config"
	"context"
	"fmt"
	codeRunner "github.com/Reholly/kforge-proto/src/gen/code-runner"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"testing"
)

type Suite struct {
	*testing.T
	Cfg    *config.Config
	Client codeRunner.CodeRunnerClient
}

func NewSuite(t *testing.T) (context.Context, *Suite) {
	t.Helper()
	t.Parallel()

	cfg := config.Config{
		Port: 8000,
		Env:  "local",
	}

	ctx, cancel := context.WithCancel(context.Background())

	t.Cleanup(func() {
		cancel()
	})

	addr := fmt.Sprintf("localhost:%d", cfg.Port)
	cc, err := grpc.DialContext(context.Background(), addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatal(err)
	}

	return ctx, &Suite{
		T:      t,
		Cfg:    &cfg,
		Client: codeRunner.NewCodeRunnerClient(cc),
	}

}
