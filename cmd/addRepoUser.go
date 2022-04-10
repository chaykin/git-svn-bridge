package cmd

import (
	"errors"
	"fmt"
	"git-svn-bridge/store"
	"git-svn-bridge/usr"
	"github.com/spf13/cobra"
)

var addRepoUserCmd = &cobra.Command{
	Use:   "add-user <repo-name> <svn-user-name>",
	Short: "Add user information for the repo",
	Long:  "Add user credentials for SVN and GIT repositories",
	Args:  cobra.ExactArgs(2),
	RunE:  addRepoUser,
}

func init() {
	repoCmd.AddCommand(addRepoUserCmd)
}

func addRepoUser(_ *cobra.Command, args []string) error {
	repoName := args[0]
	svnUserName := args[1]

	if !store.HasRepo(repoName) {
		return errors.New("There is no repo with name " + repoName)
	}

	fmt.Printf("Adding/overwriting SVN usr: %s\n", svnUserName)

	pass := readPassFromInput()
	mail := readFieldFromInput("e-mail")
	//TODO: check mail valid
	gitUserName := readFieldFromInput("Git user name")

	repo := store.GetRepo(repoName)
	user := usr.CreateUser(repo, svnUserName, pass, gitUserName, mail)
	store.StoreUser(user)

	return nil
}
