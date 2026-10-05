package cli

import (
	"github.com/spf13/cobra"
)

// run some pre checks and handle args here. this doesnt do much for now.
func Start(cmd *cobra.Command, args []string) {
	startTCP()
}
