package suite

import (
	"code-runner-service/config"
	"context"
	"testing"
)

type Suite struct {
	*testing.T
	cfg *config.Config
}

func NewSuite(t *testing.T) (context.Context, *Suite) {
	t.Helper()
	t.Parallel()

	cfg := config.Config{
		Port: 8080,
		Env:  "local",
	}

	ctx, cancel := context.WithCancel(context.Background())

	t.Cleanup(func() {
		cancel()
	})

	return ctx, &Suite{
		T:   t,
		cfg: &cfg,
	}
}
