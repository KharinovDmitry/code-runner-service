package python

import (
	executor2 "code-runner-service/internal/domain/executor"
	"code-runner-service/lib/byteconv"
	"context"
	"github.com/stretchr/testify/assert"
	"os"
	"strconv"
	"testing"
)

func TestBase(t *testing.T) {
	code, err := os.ReadFile("testFiles/test.py")
	assert.Nil(t, err)

	executor := NewPythonExecutor(byteconv.String(code), 1024, 1000)
	assert.Nil(t, executor.Init())
	defer executor.Close()

	actual, err := executor.Execute(context.Background(), "2")
	assert.Nil(t, err)

	expected := "4"
	assert.Equal(t, expected, actual)
}

func TestShutdown(t *testing.T) {
	code, err := os.ReadFile("testFiles/shutdownTest.py")
	assert.Nil(t, err)

	executor := NewPythonExecutor(byteconv.String(code), 1024, 1000)
	assert.Nil(t, executor.Init())
	defer executor.Close()

	_, err = executor.Execute(context.Background(), "2")
	assert.Nil(t, err)
}

func TestTimeLimit(t *testing.T) {
	code, err := os.ReadFile("testFiles/timeLimitTest.py")
	assert.Nil(t, err)

	executor := NewPythonExecutor(byteconv.String(code), 1024, 1000)
	assert.Nil(t, executor.Init())
	defer executor.Close()

	_, err = executor.Execute(context.Background(), "2")
	assert.ErrorIs(t, err, executor2.TimeLimitError)
}

func TestRuntimeError(t *testing.T) {
	code, err := os.ReadFile("testFiles/runtimeErrorTest.py")
	assert.Nil(t, err)

	executor := NewPythonExecutor(byteconv.String(code), 1024, 1000)
	assert.Nil(t, executor.Init())
	defer executor.Close()

	_, err = executor.Execute(context.Background(), "2")
	assert.ErrorIs(t, err, executor2.RuntimeError)
}

func TestCreateFile(t *testing.T) {
	code, err := os.ReadFile("testFiles/createFileTest.py")
	assert.Nil(t, err)

	executor := NewPythonExecutor(byteconv.String(code), 1024, 1000)
	assert.Nil(t, executor.Init())
	defer executor.Close()

	output, err := executor.Execute(context.Background(), "2")
	assert.Equal(t, false, output == "FILE CREATED")
}

func TestMemoryLimit(t *testing.T) {
	code, err := os.ReadFile("testFiles/memoryLimitError.py")
	assert.Nil(t, err)

	executor := NewPythonExecutor(byteconv.String(code), 1024, 1000)
	assert.Nil(t, executor.Init())
	defer executor.Close()

	output, err := executor.Execute(context.Background(), "2")
	assert.Equal(t, false, output == "SUCCESS")
	assert.ErrorIs(t, err, executor2.MemoryLimitError)
}

func BenchmarkSample(b *testing.B) {
	code, err := os.ReadFile("testFiles/test.py")
	assert.Nil(b, err)

	executor := NewPythonExecutor(byteconv.String(code), 1024, 1000)
	assert.Nil(b, executor.Init())
	defer executor.Close()

	for i := 0; i < b.N; i++ {
		iStr := strconv.Itoa(i)

		actual, err := executor.Execute(context.Background(), iStr)
		assert.Nil(b, err)

		expected := strconv.Itoa(i * i)
		assert.Equal(b, expected, actual)
	}
}
