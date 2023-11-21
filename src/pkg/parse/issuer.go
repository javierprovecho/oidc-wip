package parse

import (
	"fmt"

	"gopkg.in/square/go-jose.v2/jwt"
)

func GetIssuerFromToken(token string) (string, error) {
	jwtToken, err := jwt.ParseSigned(token)
	if err != nil {
		return "", err
	}

	unsafeClaims := jwt.Claims{}
	if err := jwtToken.UnsafeClaimsWithoutVerification(&unsafeClaims); err != nil {
		return "", err
	}

	return fmt.Sprint(unsafeClaims.Issuer), nil
}
