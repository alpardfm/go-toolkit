package query_test

import (
	"context"
	"fmt"

	"github.com/alpardfm/go-toolkit/heavy/query"
)

func ExampleFormatQueryForRows() {
	// Build a multi-row INSERT query
	baseQuery := "INSERT INTO users (name, age) VALUES"
	inputs := [][]any{
		{"Alice", 30},
		{"Bob", 25},
	}

	result, args, err := query.FormatQueryForRows(context.Background(), baseQuery, inputs)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	fmt.Println(result)
	fmt.Println(args)

	// Output:
	// INSERT INTO users (name, age) VALUES (?, ?), (?, ?)
	// [Alice 30 Bob 25]
}
