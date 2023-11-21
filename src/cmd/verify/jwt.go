package verify

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/auth0/go-jwt-middleware/v2/validator"
	"github.com/javierprovecho/oidc-wip/src/pkg/parse"
	"github.com/javierprovecho/oidc-wip/src/pkg/verify"

	"github.com/spf13/cobra"
)

var jwtCmd = &cobra.Command{
	Use:   "jwt",
	Short: "verify a jwt token",
	Run: func(cmd *cobra.Command, args []string) {
		var content []byte

		switch filename {
		case "-":
			content, _ = os.ReadFile("/dev/stdin")
		default:
			content, _ = os.ReadFile(filename)
		}

		var namespace, serviceAccount, pod string

		if len(args) == 1 {
			issuer = parse.GetIssuerFromUri(args[0])
			namespace, serviceAccount, pod = parse.GetSubFromURI(args[0])
		}
		if issuerFromEnv := os.Getenv("ISSUER_URI"); issuerFromEnv != "" {
			issuer = parse.GetIssuerFromUri(issuerFromEnv)
			namespace, serviceAccount, pod = parse.GetSubFromURI(issuerFromEnv)
		}
		if audienceFromEnv := os.Getenv("AUDIENCE"); audienceFromEnv != "" {
			audience = audienceFromEnv
		}

		fmt.Println(parse.GetIssuerFromToken(string(content)))

		var claims validator.RegisteredClaims
		var isIssuerAndAudienceVerified bool

		if issuer != "" {
			claims, isIssuerAndAudienceVerified = verify.VerifyTokenWithIssuer(string(content), issuer, audience)
		} else {
			claims, isIssuerAndAudienceVerified = verify.VerifyToken(string(content), audience)
		}

		isSubVerified := verify.VerifyTokenWithSub(string(content), namespace, serviceAccount, pod)

		marshalledJson, _ := json.MarshalIndent(claims, "", "    ")
		fmt.Fprintln(os.Stdout, string(marshalledJson))
		fmt.Fprintf(os.Stderr, "verified issuer and audience: %t\n", isIssuerAndAudienceVerified)
		fmt.Fprintf(os.Stderr, "verified sub: %t\n", isSubVerified)

	},
}

var filename string

func init() {
	VerifyCmd.AddCommand(jwtCmd)

	jwtCmd.PersistentFlags().StringVarP(&filename, "filename", "f", "-", "path to file containing jwt token, use '-' for stdin")
}
