package query

import (
	"context"
	"fmt"
	"strings"

	"github.com/alpardfm/go-toolkit/codes"
	"github.com/alpardfm/go-toolkit/errors"
)

// FormatQueryForRows builds a multi-row INSERT query by appending parameterized value placeholders
// for each row of inputs. It returns the formatted query string and a flattened slice of arguments.
func FormatQueryForRows(ctx context.Context, q string, inputs [][]any) (string, []any, error) {
	// Add () based on rows
	// Add ? based on cols
	lRow := len(inputs)
	if lRow < 1 {
		return "", nil, errors.NewWithCode(codes.CodeSQLPrepareStmt, "no inputs rows supplied")
	}
	lCol := len(inputs[0])
	if lCol < 1 {
		return "", nil, errors.NewWithCode(codes.CodeSQLPrepareStmt, "no inputs cols supplied")
	}

	ins := []string{}
	for i := 0; i < lCol; i++ {
		ins = append(ins, "?")
	}
	inputTemplate := fmt.Sprintf("(%s)", strings.Join(ins, ", "))

	insQuery := []string{}
	for i := 0; i < lRow; i++ {
		insQuery = append(insQuery, inputTemplate)
	}

	result := fmt.Sprintf("%s %s", q, strings.Join(insQuery, ", "))

	args := []any{}
	for _, r := range inputs {
		args = append(args, r...)
	}

	return result, args, nil
}
