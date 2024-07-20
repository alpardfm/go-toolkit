package tokens

import (
	"reflect"

	"github.com/alpardfm/go-toolkit/codes"
	"github.com/alpardfm/go-toolkit/errors"
	"github.com/dgrijalva/jwt-go/v4"
)

func NewJWTToken[claimsType jwt.Claims](claims claimsType, secretKey []byte) (string, error) {
	jwtClaim := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	jwtToken, err := jwtClaim.SignedString(secretKey)
	if err != nil {
		return "", errors.NewWithCode(codes.CodeJWTSignedStringError, err.Error())
	}

	return jwtToken, nil
}

func ValidateJWTToken[claimsType jwt.Claims](tokenString string, secretKey []byte, claims claimsType) (*jwt.Token, error) {
	typeOfClaims := reflect.TypeOf(claims)
	if typeOfClaims.Kind() != reflect.Pointer {
		return nil, errors.NewWithCode(codes.CodeInvalidValue, "claims must be a pointer")
	}

	keyFunc := func(token *jwt.Token) (any, error) {
		if _, isOk := token.Method.(*jwt.SigningMethodHMAC); !isOk {
			return nil, errors.NewWithCode(codes.CodeJWTInvalidMethod, "invalid method algorithm")
		}

		return secretKey, nil
	}

	jwtToken, err := jwt.ParseWithClaims(tokenString, claims, keyFunc)
	if err != nil {
		return nil, errors.NewWithCode(codes.CodeJWTParseWithClaimsError, err.Error())
	}

	return jwtToken, nil
}

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
