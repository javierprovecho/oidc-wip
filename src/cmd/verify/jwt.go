package verify

import (
	"encoding/json"
	"fmt"
	"os"

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

		if len(args) == 1 {
			issuer = parse.GetIssuerFromUri(args[0])
		}
		if issuerFromEnv := os.Getenv("ISSUER_URI"); issuerFromEnv != "" {
			issuer = parse.GetIssuerFromUri(issuerFromEnv)
		}
		if audienceFromEnv := os.Getenv("AUDIENCE"); audienceFromEnv != "" {
			audience = audienceFromEnv
		}

		fmt.Println(parse.GetIssuer(string(content)))

		if issuer != "" {
			claims, isVerified := verify.VerifyTokenWithIssuer(issuer, audience, string(content))
			marshalledJson, _ := json.MarshalIndent(claims, "", "    ")
			fmt.Fprint(os.Stdout, marshalledJson)
			fmt.Fprintf(os.Stderr, "verified: %t", isVerified)
		} else {
			claims, isVerified := verify.VerifyToken(audience, string(content))
			marshalledJson, _ := json.MarshalIndent(claims, "", "    ")
			fmt.Fprint(os.Stdout, marshalledJson)
			fmt.Fprintf(os.Stderr, "verified: %t", isVerified)

		}

	},
}

var filename string

func init() {
	VerifyCmd.AddCommand(jwtCmd)

	jwtCmd.PersistentFlags().StringVarP(&filename, "filename", "f", "-", "path to file containing jwt token, use '-' for stdin")
}
