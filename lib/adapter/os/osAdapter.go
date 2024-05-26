package os

import "os"

type OSAdapter interface {
	CreateTempFileWithText(text string, extension string) (fileName string, err error)
	AddFileExecutablePermission(file *os.File) error
}
