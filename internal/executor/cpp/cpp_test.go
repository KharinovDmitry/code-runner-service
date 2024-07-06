package cpp

import (
	"code-runner-service/lib/byteconv"
	"context"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestBase(t *testing.T) {
	code, err := os.ReadFile("testFiles/test.cpp")
	assert.Nil(t, err)

	executor := NewCPPExecutor(byteconv.String(code), 1024, 1000)
	err = executor.Init()
	assert.Nil(t, err)
	defer executor.Close()

	actual, err := executor.Execute(context.Background(), "2")
	assert.Nil(t, err)

	expected := "4"
	assert.Equal(t, expected, actual)
}
