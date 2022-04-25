package utils

import (
	"gopkg.in/square/go-jose.v2/jwt"
)

func ExtractPayload(tokenString string) (map[string]interface{}, error) {
	var payload map[string]interface{}

	token, err := jwt.ParseSigned(tokenString)

	if err != nil {
		return nil, err
	}

	err = token.UnsafeClaimsWithoutVerification(&payload)

	if err != nil {
		return nil, err
	}

	return payload, nil
}
