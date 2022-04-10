package cmd

import (
	"git-svn-bridge/repo"
	"git-svn-bridge/store"
	"github.com/spf13/cobra"
	"strings"
)

var createRepoName string
var createRepoCmd = &cobra.Command{
	Use:   "create <url-to-svn-repo>",
	Short: "Create repository configuration",
	Args:  cobra.ExactArgs(1),
	Run:   createRepo,
}

func init() {
	createRepoCmd.PersistentFlags().StringVarP(&createRepoName, "name", "n", "", "repository name")
	repoCmd.AddCommand(createRepoCmd)
}

func createRepo(_ *cobra.Command, args []string) {
	svnUrl := args[0]
	if createRepoName == "" {
		index := strings.LastIndex(svnUrl, "/")
		createRepoName = svnUrl[index+1:]
	}

	store.StoreRepo(repo.CreateRepo(createRepoName, svnUrl))
}
