package tokens_test

import (
	"fmt"
	"strings"

	"github.com/alpardfm/go-toolkit/heavy/tokens"
	"github.com/golang-jwt/jwt/v5"
)

func ExampleNewJWTToken() {
	// Define custom claims
	type MyClaims struct {
		UserID string `json:"user_id"`
		jwt.RegisteredClaims
	}

	claims := MyClaims{
		UserID: "user-123",
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer: "my-app",
		},
	}

	// Create a JWT token signed with HS256
	secretKey := []byte("my-secret-key-for-signing")
	token, err := tokens.NewJWTToken(claims, secretKey)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	// JWT tokens have 3 parts separated by dots
	parts := strings.Split(token, ".")
	fmt.Println(len(parts))

	// Output:
	// 3
}
