package cmd

import (
	"github.com/spf13/cobra"
)

var repoCmd = &cobra.Command{
	Use:   "repo",
	Short: "Working with repos configuration",
	Long:  "Create/Update repositories configuration",
}

func init() {
	rootCmd.AddCommand(repoCmd)
}
