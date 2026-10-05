package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "kewl-db",
	Short: "kewl-db is pretty kewl",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Welcome to kewl-db.")
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
