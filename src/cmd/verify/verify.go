package verify

import (
	"github.com/spf13/cobra"
)

var (
	audience string
	issuer   string
)

var VerifyCmd = &cobra.Command{
	Use: "verify",
}

func init() {
	VerifyCmd.PersistentFlags().StringVar(&audience, "audience", "", "jwt audience to be verified")
	VerifyCmd.PersistentFlags().StringVar(&issuer, "issuer", "", "jwt issuer to be verified")
}
