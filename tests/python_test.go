package tests

import (
	"code-runner-service/internal/domain/enum"
	"code-runner-service/lib/byteconv"
	"code-runner-service/tests/suite"
	"fmt"
	codeRunner "github.com/Reholly/kforge-proto/src/gen/code-runner"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestHappyPath(t *testing.T) {
	ctx, st := suite.NewSuite(t)

	code, err := os.ReadFile("testFiles/test.py")
	assert.Nil(t, err)

	tests := []*codeRunner.Test{
		{
			Input:  "1",
			Output: "1",
			Points: 10,
		},
		{
			Input:  "2",
			Output: "4",
			Points: 20,
		},
		{
			Input:  "-1",
			Output: "1",
			Points: 70,
		},
	}

	res, err := st.Client.RunTestsOnCode(ctx, &codeRunner.RunTestsOnCodeRequest{
		Code:          byteconv.String(code),
		Language:      string(enum.Python),
		MemoryLimitKb: 1024,
		TimeLimitMs:   1024,
		Tests:         tests,
	})
	assert.Nil(t, err)

	fmt.Println(res)
}
