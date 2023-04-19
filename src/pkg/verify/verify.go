package verify

import (
	"context"
	"log"
	"net/url"
	"time"

	"github.com/auth0/go-jwt-middleware/v2/jwks"
	"github.com/auth0/go-jwt-middleware/v2/validator"
	"github.com/javierprovecho/oidc-wip/src/pkg/parse"
)

func VerifyToken(audience, token string) (validator.RegisteredClaims, bool) {
	issuer, err := parse.GetIssuer(token)
	if err != nil {
		log.Print(err)
		return validator.RegisteredClaims{}, false
	}

	return VerifyTokenWithIssuer(issuer, audience, token)
}

func VerifyTokenWithIssuer(issuer, audience, token string) (validator.RegisteredClaims, bool) {

	issuerURL, err := url.Parse(issuer)
	if err != nil {
		log.Fatalf("failed to parse the issuer url: %v", err)
	}

	provider := jwks.NewCachingProvider(issuerURL, 5*time.Minute)

	// Set up the validator.
	jwtValidator, err := validator.New(
		provider.KeyFunc,
		validator.RS256,
		issuerURL.String(),
		[]string{audience},
	)
	if err != nil {
		log.Fatalf("failed to set up the validator: %v", err)
	}

	claims, err := jwtValidator.ValidateToken(context.Background(), token)
	if err != nil {
		log.Fatalf("failed to validate token: %v", err)
		return validator.RegisteredClaims{}, false
	}

	return claims.(*validator.ValidatedClaims).RegisteredClaims, true
}
