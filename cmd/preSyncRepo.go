package cmd

import (
	"fmt"
	"git-svn-bridge/log"
	"git-svn-bridge/sync"
	"github.com/spf13/cobra"
)

var preSyncRepoName string
var preSyncRepoCmd = &cobra.Command{
	Use:   "pre-sync <ref> <old-sha> <new-sha>",
	Short: "Checks if sync SVN and GIT repositories is available",
	Run:   preSyncRepo,
}

func init() {
	preSyncRepoCmd.PersistentFlags().StringVarP(&preSyncRepoName, "repo", "r", "", "repository name")
	if err := preSyncRepoCmd.MarkPersistentFlagRequired("repo"); err != nil {
		log.Fatalf("could not init pre-sync command: %w", err)
	}

	rootCmd.AddCommand(preSyncRepoCmd)
}

func preSyncRepo(_ *cobra.Command, args []string) {
	defer log.StdErrOnPanicf(fmt.Errorf("could not pre-sync repository '%s'", preSyncRepoName))

	sync.NewPreSync(preSyncRepoName).PreSync(args[0], args[1], args[2])
}
