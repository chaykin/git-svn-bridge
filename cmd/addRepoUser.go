package cmd

import (
	"fmt"
	"git-svn-bridge/log"
	"git-svn-bridge/store"
	"git-svn-bridge/usr"
	"github.com/spf13/cobra"
)

var addRepoUserCmd = &cobra.Command{
	Use:   "add-user <repo-name> <svn-user-name>",
	Short: "Add user information for the repo",
	Long:  "Add user credentials for SVN and GIT repositories",
	Args:  cobra.ExactArgs(2),
	Run:   addRepoUser,
}

func init() {
	rootCmd.AddCommand(addRepoUserCmd)
}

func addRepoUser(_ *cobra.Command, args []string) {
	repoName := args[0]
	svnUserName := args[1]

	defer log.StdErrOnPanicf(fmt.Errorf("could not add user '%s' for repo '%s'", svnUserName, repoName))

	if !store.HasRepo(repoName) {
		panic(fmt.Errorf("there is no repo with name %s", repoName))
	}

	fmt.Printf("Adding/overwriting SVN usr: %s\n", svnUserName)

	pass := readPassFromInput()
	//TODO: check mail valid
	mail := readFieldFromInput("e-mail")
	gitUserName := readFieldFromInput("Git user name")

	repo := store.GetRepo(repoName)
	user := usr.CreateUser(&repo, svnUserName, pass, gitUserName, mail)
	store.StoreUser(user)
}
