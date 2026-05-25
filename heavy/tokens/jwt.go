package tokens

import (
	"reflect"

	"github.com/alpardfm/go-toolkit/codes"
	"github.com/alpardfm/go-toolkit/errors"
	"github.com/golang-jwt/jwt/v5"
)

// NewJWTToken creates a new JWT token signed with HS256 using the provided claims and secret key.
func NewJWTToken[claimsType jwt.Claims](claims claimsType, secretKey []byte) (string, error) {
	jwtClaim := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	jwtToken, err := jwtClaim.SignedString(secretKey)
	if err != nil {
		return "", errors.NewWithCode(codes.CodeJWTSignedStringError, err.Error())
	}

	return jwtToken, nil
}

// ValidateJWTToken parses and validates a JWT token string using the provided secret key and claims type.
// The claims parameter must be a pointer type.
func ValidateJWTToken[claimsType jwt.Claims](tokenString string, secretKey []byte, claims claimsType) (*jwt.Token, error) {
	typeOfClaims := reflect.TypeOf(claims)
	if typeOfClaims.Kind() != reflect.Pointer {
		return nil, errors.NewWithCode(codes.CodeInvalidValue, "claims must be a pointer")
	}

	keyFunc := func(token *jwt.Token) (any, error) {
		return secretKey, nil
	}

	jwtToken, err := jwt.ParseWithClaims(tokenString, claims, keyFunc, jwt.WithValidMethods([]string{"HS256", "HS384", "HS512"}))
	if err != nil {
		return nil, errors.NewWithCode(codes.CodeJWTParseWithClaimsError, err.Error())
	}

	return jwtToken, nil
}

// GetClaimsOfJWTToken extracts the typed claims from a validated JWT token.
func GetClaimsOfJWTToken[claimsType jwt.Claims](jwtToken jwt.Token) (claimsType, error) {
	claims, isOk := jwtToken.Claims.(claimsType)
	if !isOk {
		typeOfClaims := reflect.TypeOf(claims)
		if typeOfClaims.Kind() != reflect.Pointer {
			return claims, errors.NewWithCode(codes.CodeInvalidValue, "claims must be a pointer")
		}

		return claims, errors.NewWithCode(codes.CodeJWTInvalidClaimsType, "claims type is not equals")
	}

	return claims, nil
}
