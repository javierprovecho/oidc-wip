package fetch

import (
	"fmt"

	"github.com/spf13/cobra"
)

var gceCmd = &cobra.Command{
	Use:   "gce",
	Short: "fetch a token from gce metadata service",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Hugo Static Site Generator v0.9 -- HEAD")
	},
}

func init() {
	FetchCmd.AddCommand(gceCmd)
}
