package format

import (
	"fmt"
	"time"

	"github.com/alpardfm/go-toolkit/codes"
	"github.com/alpardfm/go-toolkit/errors"
)

const (
	// DayMonthYearHourMinSecMilisec is a time format with milliseconds
	DayMonthYearHourMinSecMilisec = "02/01/2006 15:04:05.000000000"
	// DayMonthYearHourMinSec is a time format without milliseconds
	DayMonthYearHourMinSec = "02/01/2006 15:04:05"
	// DayMonthYear is a date format
	DayMonthYear = "02/01/2006"
	// HourMinSec is a time format for hours, minutes, and seconds
	HourMinSec = "15:04:05"
)

// TimeParseWithDefaultFormat parses a time string with the default format "02/01/2006 15:04:05".
// If parsing fails, it returns a zero time and an error with an appropriate error code.
func TimeParseWithDefaultFormat(value string) (time.Time, error) {
	result, err := time.Parse(DayMonthYearHourMinSec, value)
	if err != nil {
		return time.Time{}, errors.NewWithCode(codes.CodeInvalidValue, fmt.Sprintf("invalid time format: %v", err))
	}
	return result, nil
}
