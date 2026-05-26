package tokens

import (
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type testClaims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

func TestNewJWTToken(t *testing.T) {
	secretKey := []byte("test-secret-key-for-signing-jwt")

	tests := []struct {
		name    string
		claims  testClaims
		key     []byte
		wantErr bool
	}{
		{
			name: "valid claims",
			claims: testClaims{
				UserID: "user-123",
				RegisteredClaims: jwt.RegisteredClaims{
					Issuer:    "test-app",
					ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
				},
			},
			key:     secretKey,
			wantErr: false,
		},
		{
			name: "empty claims",
			claims: testClaims{
				RegisteredClaims: jwt.RegisteredClaims{},
			},
			key:     secretKey,
			wantErr: false,
		},
		{
			name: "claims with all registered fields",
			claims: testClaims{
				UserID: "user-456",
				RegisteredClaims: jwt.RegisteredClaims{
					Issuer:    "my-app",
					Subject:   "auth",
					Audience:  jwt.ClaimStrings{"api"},
					ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
					IssuedAt:  jwt.NewNumericDate(time.Now()),
					NotBefore: jwt.NewNumericDate(time.Now()),
					ID:        "unique-id-123",
				},
			},
			key:     secretKey,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := NewJWTToken(tt.claims, tt.key)
			if tt.wantErr {
				if err == nil {
					t.Error("NewJWTToken() expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("NewJWTToken() unexpected error: %v", err)
			}
			if token == "" {
				t.Error("NewJWTToken() returned empty token")
			}
			// JWT has 3 parts
			parts := strings.Split(token, ".")
			if len(parts) != 3 {
				t.Errorf("NewJWTToken() token has %d parts, want 3", len(parts))
			}
		})
	}
}

func TestValidateJWTToken(t *testing.T) {
	secretKey := []byte("test-secret-key-for-signing-jwt")
	wrongKey := []byte("wrong-secret-key-for-validation")

	// Create a valid token
	validClaims := testClaims{
		UserID: "user-123",
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "test-app",
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	validToken, err := NewJWTToken(validClaims, secretKey)
	if err != nil {
		t.Fatalf("setup: NewJWTToken() error: %v", err)
	}

	// Create an expired token
	expiredClaims := testClaims{
		UserID: "user-expired",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(-time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now().Add(-2 * time.Hour)),
		},
	}
	expiredToken, err := NewJWTToken(expiredClaims, secretKey)
	if err != nil {
		t.Fatalf("setup: NewJWTToken() error: %v", err)
	}

	tests := []struct {
		name    string
		token   string
		key     []byte
		claims  *testClaims
		wantErr bool
		errMsg  string
	}{
		{
			name:    "valid token with correct key",
			token:   validToken,
			key:     secretKey,
			claims:  &testClaims{},
			wantErr: false,
		},
		{
			name:    "valid token with wrong key",
			token:   validToken,
			key:     wrongKey,
			claims:  &testClaims{},
			wantErr: true,
			errMsg:  "CodeJWTParseWithClaimsError",
		},
		{
			name:    "expired token",
			token:   expiredToken,
			key:     secretKey,
			claims:  &testClaims{},
			wantErr: true,
			errMsg:  "CodeJWTParseWithClaimsError",
		},
		{
			name:    "malformed token",
			token:   "not.a.valid.jwt.token",
			key:     secretKey,
			claims:  &testClaims{},
			wantErr: true,
			errMsg:  "CodeJWTParseWithClaimsError",
		},
		{
			name:    "empty token",
			token:   "",
			key:     secretKey,
			claims:  &testClaims{},
			wantErr: true,
			errMsg:  "CodeJWTParseWithClaimsError",
		},
		{
			name:    "non-pointer claims",
			token:   validToken,
			key:     secretKey,
			claims:  nil, // will use non-pointer below
			wantErr: true,
			errMsg:  "claims must be a pointer",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var jwtToken *jwt.Token
			var err error

			if tt.name == "non-pointer claims" {
				// Test with non-pointer claims
				jwtToken, err = ValidateJWTToken(tt.token, tt.key, testClaims{})
			} else {
				jwtToken, err = ValidateJWTToken(tt.token, tt.key, tt.claims)
			}

			if tt.wantErr {
				if err == nil {
					t.Error("ValidateJWTToken() expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("ValidateJWTToken() unexpected error: %v", err)
			}
			if jwtToken == nil {
				t.Error("ValidateJWTToken() returned nil token")
			}
			if !jwtToken.Valid {
				t.Error("ValidateJWTToken() token is not valid")
			}
		})
	}
}

func TestGetClaimsOfJWTToken(t *testing.T) {
	secretKey := []byte("test-secret-key-for-signing-jwt")

	// Create and validate a token
	originalClaims := testClaims{
		UserID: "user-789",
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "test-app",
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		},
	}
	tokenStr, err := NewJWTToken(originalClaims, secretKey)
	if err != nil {
		t.Fatalf("setup: NewJWTToken() error: %v", err)
	}

	jwtToken, err := ValidateJWTToken(tokenStr, secretKey, &testClaims{})
	if err != nil {
		t.Fatalf("setup: ValidateJWTToken() error: %v", err)
	}

	t.Run("extract valid claims", func(t *testing.T) {
		claims, err := GetClaimsOfJWTToken[*testClaims](*jwtToken)
		if err != nil {
			t.Fatalf("GetClaimsOfJWTToken() error: %v", err)
		}
		if claims.UserID != "user-789" {
			t.Errorf("GetClaimsOfJWTToken() UserID = %v, want %v", claims.UserID, "user-789")
		}
		if claims.Issuer != "test-app" {
			t.Errorf("GetClaimsOfJWTToken() Issuer = %v, want %v", claims.Issuer, "test-app")
		}
	})

	t.Run("wrong claims type", func(t *testing.T) {
		type otherClaims struct {
			jwt.RegisteredClaims
		}
		_, err := GetClaimsOfJWTToken[*otherClaims](*jwtToken)
		if err == nil {
			t.Error("GetClaimsOfJWTToken() expected error for wrong type, got nil")
		}
	})
}

func TestJWTRoundTrip(t *testing.T) {
	secretKey := []byte("round-trip-secret-key-32-bytes!!")

	claims := testClaims{
		UserID: "roundtrip-user",
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "roundtrip-app",
			Subject:   "auth",
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	// Create token
	tokenStr, err := NewJWTToken(claims, secretKey)
	if err != nil {
		t.Fatalf("NewJWTToken() error: %v", err)
	}

	// Validate token
	jwtToken, err := ValidateJWTToken(tokenStr, secretKey, &testClaims{})
	if err != nil {
		t.Fatalf("ValidateJWTToken() error: %v", err)
	}

	// Extract claims
	extracted, err := GetClaimsOfJWTToken[*testClaims](*jwtToken)
	if err != nil {
		t.Fatalf("GetClaimsOfJWTToken() error: %v", err)
	}

	// Verify
	if extracted.UserID != "roundtrip-user" {
		t.Errorf("round trip UserID = %v, want %v", extracted.UserID, "roundtrip-user")
	}
	if extracted.Issuer != "roundtrip-app" {
		t.Errorf("round trip Issuer = %v, want %v", extracted.Issuer, "roundtrip-app")
	}
	if extracted.Subject != "auth" {
		t.Errorf("round trip Subject = %v, want %v", extracted.Subject, "auth")
	}
}
