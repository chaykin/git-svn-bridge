package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   "git-svn-bridge command",
	Short: "Git-svn-bridge is a bi-directional bridge between GIT and SVN repositories",
	Long:  "A bi-directional bridge than allows work with SVN repo like it is a (almost) orginal GIT repo",
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		_, err = fmt.Fprintln(os.Stderr, err)
		if err != nil {
			panic(fmt.Errorf("could not print to stderr: %w", err))
		}
	}
}
