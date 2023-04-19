package cmd

import (
	"github.com/javierprovecho/oidc-wip/src/cmd/fetch"
)

func init() {
	rootCmd.AddCommand(fetch.FetchCmd)
}
