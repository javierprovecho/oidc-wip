package cmd

import (
	"github.com/javierprovecho/oidc-wip/src/cmd/verify"
)

func init() {
	rootCmd.AddCommand(verify.VerifyCmd)
}
