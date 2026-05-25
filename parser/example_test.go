package parser_test

import (
	"fmt"

	"github.com/alpardfm/go-toolkit/parser"
)

func ExampleInitParser() {
	// Initialize a parser with standard JSON configuration
	p := parser.InitParser(parser.Options{
		JSONOptions: parser.JSONOptions{},
	})

	// Marshal a struct to JSON
	type User struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}

	data, err := p.JSONParser().Marshal(User{Name: "Alice", Age: 30})
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Println(string(data))

	// Unmarshal JSON back to a struct
	var user User
	err = p.JSONParser().Unmarshal(data, &user)
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Println(user.Name, user.Age)

	// Output:
	// {"name":"Alice","age":30}
	// Alice 30
}
