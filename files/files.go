package files

import (
	"fmt"
	"os"
	"runtime"
	"strings"
)

// Extract file extension e.g. "txt, csv, docx" from full filename.
// If file has no extension then return empty string
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

// Checks if a file exists and is not a directory
func IsExist(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}

	return !info.IsDir()
}

// Will return string of current file location where this function is called.
func GetCurrentFileLocation() string {
	_, file, _, isOk := runtime.Caller(1)
	if isOk {
		return file
	}
	return ""
}

// Will return string of current method location where this function is called.
func GetCurrentMethodName() string {
	pc, _, _, isOk := runtime.Caller(1)
	if !isOk {
		return ""
	}

	f := runtime.FuncForPC(pc)
	return fmt.Sprintf("%v()", f.Name())
}
