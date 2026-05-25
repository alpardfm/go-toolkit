package files

import (
	"fmt"
	"os"
	"runtime"
	"strings"
)

// GetExtension extracts the file extension (e.g. "txt", "csv", "docx") from a filename.
// Returns an empty string if the file has no extension.
func GetExtension(filename string) string {
	filenamewithext := strings.Split(filename, ".")
	if len(filenamewithext) < 1 {
		return ""
	}

	fileextension := filenamewithext[len(filenamewithext)-1]
	if fileextension == filename {
		return ""
	}
	return fileextension
}

// IsExist checks if a file exists at the given path and is not a directory.
func IsExist(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}

	return !info.IsDir()
}

// GetCurrentFileLocation returns the absolute path of the source file where this function is called.
func GetCurrentFileLocation() string {
	_, file, _, isOk := runtime.Caller(1)
	if isOk {
		return file
	}
	return ""
}

// GetCurrentMethodName returns the fully qualified name of the function where this is called.
func GetCurrentMethodName() string {
	pc, _, _, isOk := runtime.Caller(1)
	if !isOk {
		return ""
	}

	f := runtime.FuncForPC(pc)
	return fmt.Sprintf("%v()", f.Name())
}
