package cmd

import (
	"github.com/aidenfine/kewl-db/pkg/cli"
	"github.com/spf13/cobra"
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "start the kewl-db server with a single node.",
	Run: func(cmd *cobra.Command, args []string) {
		cli.Start(cmd, args)
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
}
