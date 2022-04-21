package cmd

import (
	"git-svn-bridge/sync"
	"github.com/spf13/cobra"
)

var syncRepoName string
var syncRepoCmd = &cobra.Command{
	Use:   "sync <ref1>..<refN> ",
	Short: "Create repository configuration",
	RunE:  syncRepo,
}

func init() {
	syncRepoCmd.PersistentFlags().StringVarP(&syncRepoName, "repo", "r", "", "repository name")
	err := syncRepoCmd.MarkPersistentFlagRequired("repo")
	if err != nil {
		panic(err)
	}

	rootCmd.AddCommand(syncRepoCmd)
}

func syncRepo(_ *cobra.Command, args []string) error {
	man, err := sync.New(syncRepoName)
	if err != nil {
		return err
	}

	if len(args) == 0 {
		return man.SyncAllRefs()
	} else {
		return man.SyncRefs(args)
	}
}
