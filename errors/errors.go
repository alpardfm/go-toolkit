package errors

import (
	"fmt"
	"net/http"
	"runtime"
	"strings"

	"github.com/alpardfm/go-toolkit/codes"
	"github.com/alpardfm/go-toolkit/language"
	"github.com/alpardfm/go-toolkit/operator"
)

// App represents an application-level error with a code, user-facing title and body, and an underlying system error.
type App struct {
	Code  codes.Code `json:"code"`
	Title string     `json:"title"`
	Body  string     `json:"body"`
	sys   error
}

// Error returns the string representation of the underlying system error.
func (e *App) Error() string {
	return e.sys.Error()
}

// Compile returns an error and creates new App errors
func Compile(err error, lang string) (int, App) {
	code := GetCode(err)

	if appErr, ok := codes.ErrorMessages[code]; ok {
		return appErr.StatusCode, App{
			Code:  code,
			Title: operator.TernaryString(lang == language.Indonesian, appErr.TitleID, appErr.TitleEN),
			Body:  operator.TernaryString(lang == language.Indonesian, appErr.BodyID, appErr.BodyEN),
			sys:   err,
		}
	}

	// Default Error
	return http.StatusInternalServerError, App{
		Code:  code,
		Title: "Service Error Not Defined",
		Body:  "Unknown error. Please contact admin",
		sys:   err,
	}
}

// NewWithCode creates a new error with a stacktrace, error code, and formatted message.
func NewWithCode(code codes.Code, msg string, val ...any) error {
	return create(nil, code, msg, val...)
}

// GetCaller returns the file, line number, message, and any error from the stacktrace.
// If the error is not a stacktrace type, it returns an error indicating the cast failure.
func GetCaller(err error) (string, int, string, error) {
	if st, isOk := err.(*stacktrace); isOk {
		return st.file, st.line, st.message, nil
	} else {
		return "", 0, "", create(nil, codes.NoCode, operator.Ternary(err == nil, "failed to cast error to stacktrace", err.Error()))
	}
}

func create(cause error, code codes.Code, msg string, val ...any) error {
	if code == codes.NoCode {
		code = GetCode(cause)
	}

	err := &stacktrace{
		message: fmt.Sprintf(msg, val...),
		cause:   cause,
		code:    code,
	}

	pc, file, line, ok := runtime.Caller(2)
	if !ok {
		return err
	}
	err.file, err.line = file, line

	f := runtime.FuncForPC(pc)
	if f == nil {
		return err
	}
	err.function = shortFuncName(f)

	return err
}

func shortFuncName(f *runtime.Func) string {
	// f.Name() is like one of these:
	// - "github.com/alpardfm/go-toolkit/<package>.<FuncName>"
	// - "github.com/alpardfm/go-toolkit/<package>.<Receiver>.<MethodName>"
	// - "github.com/alpardfm/go-toolkit/<package>.<*PtrReceiver>.<MethodName>"
	longName := f.Name()

	withoutPath := longName[strings.LastIndex(longName, "/")+1:]
	withoutPackage := withoutPath[strings.Index(withoutPath, ".")+1:]

	shortName := withoutPackage
	shortName = strings.Replace(shortName, "(", "", 1)
	shortName = strings.Replace(shortName, "*", "", 1)
	shortName = strings.Replace(shortName, ")", "", 1)

	return shortName
}

// GetCode extracts the error code from a stacktrace error.
// If the error is not a stacktrace type, it returns NoCode.
func GetCode(err error) codes.Code {
	if err, ok := err.(*stacktrace); ok {
		return err.code
	}
	return codes.NoCode
}
