package enum

// Language код языка
// @Model domain.Language
// @Description CPP - C++,
// @Description Python - Python,
type Language string

const (
	CPP        Language = "CPP"
	Python     Language = "Python"
	CSharp     Language = "C#"
	JavaScript Language = "JS"
)
