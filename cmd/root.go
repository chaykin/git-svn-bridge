package cmd

import (
	"git-svn-bridge/log"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "git-svn-bridge command",
	Short: "Git-svn-bridge is a bi-directional bridge between GIT and SVN repositories",
	Long:  "A bi-directional bridge that allows work with SVN repo like it is a (almost) orginal GIT repo",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("could not execute root command: %w", err)
	}
}
