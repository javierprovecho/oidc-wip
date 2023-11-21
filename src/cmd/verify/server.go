package verify

import (
	"os"

	"github.com/javierprovecho/oidc-wip/src/pkg/parse"
	"github.com/javierprovecho/oidc-wip/src/pkg/server"
	"github.com/spf13/cobra"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "start a verification server",
	RunE: func(cmd *cobra.Command, args []string) error {
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
		return server.Server(issuer, audience, namespace, serviceAccount, pod)
	},
}

func init() {
	VerifyCmd.AddCommand(serverCmd)
}
