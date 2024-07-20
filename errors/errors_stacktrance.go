package errors

import (
	"fmt"

	"github.com/alpardfm/go-toolkit/codes"
)

type stacktrace struct {
	message  string
	cause    error
	code     codes.Code
	file     string
	function string
	line     int
}

// Error method returns the message of the stacktrace
func (st *stacktrace) Error() string {
	return fmt.Sprintf("Error: %s", st.message)
}

// ExitCode method returns an appropriate exit code based on the code value
func (st *stacktrace) ExitCode() int {
	if st.code == codes.NoCode {
		return 1
	}
	return int(st.code)
}
