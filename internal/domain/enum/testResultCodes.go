package enum

type TestResultCode string

const (
	TimeLimitCode       TestResultCode = "TL"
	MemoryLimitCode     TestResultCode = "ML"
	CompileErrorCode    TestResultCode = "CE"
	RuntimeErrorCode    TestResultCode = "RE"
	SuccessCode         TestResultCode = "SC"
	IncorrectAnswerCode TestResultCode = "IA"
)
