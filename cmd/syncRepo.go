package cmd

import (
	"fmt"
	"git-svn-bridge/log"
	"git-svn-bridge/sync"
	"github.com/spf13/cobra"
)

var syncRepoName string
var syncRepoCmd = &cobra.Command{
	Use:   "sync <ref1>..<refN> ",
	Short: "Sync SVN and GIT repositories",
	Long:  "Sync SVN and GIT repositories with bridge GIT repository in the middle",
	Run:   syncRepo,
}

func init() {
	syncRepoCmd.PersistentFlags().StringVarP(&syncRepoName, "repo", "r", "", "repository name")
	if err := syncRepoCmd.MarkPersistentFlagRequired("repo"); err != nil {
		log.Fatalf("could not init sync command: %w", err)
	}

	rootCmd.AddCommand(syncRepoCmd)
}

func syncRepo(_ *cobra.Command, args []string) {
	defer log.OnPanicf(fmt.Errorf("could not sync repository '%s'", syncRepoName))

	man := sync.New(syncRepoName)
	if len(args) == 0 {
		man.SyncAllRefs()
	} else {
		man.SyncRefs(args)
	}
}
