package cmd

import (
	"errors"
	"fmt"
	"git-svn-bridge/gitsvn"
	"git-svn-bridge/repo"
	"git-svn-bridge/store"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/spf13/cobra"
)

var initRepoCmd = &cobra.Command{
	Use:   "init <repo-name> <git-user-name>",
	Short: "Initialize repo",
	Long:  "Initialize repo: create GIT server repo, bridge repo, config hooks, etc",
	Args:  cobra.ExactArgs(2),
	RunE:  initRepo,
}

func init() {
	rootCmd.AddCommand(initRepoCmd)
}

func initRepo(_ *cobra.Command, args []string) error {
	repoName := args[0]
	gitUserName := args[1]
	if !store.HasRepo(repoName) {
		return errors.New("There is no repo with name " + repoName)
	}

	repo := store.GetRepo(repoName)

	err := initGitRepo(&repo)
	if err != nil {
		return fmt.Errorf("could not init git repo '%s': %w", repoName, err)
	}
	err = initBridgeRepo(&repo, gitUserName)
	if err != nil {
		return fmt.Errorf("could not init bridge repo '%s': %w", repoName, err)
	}

	//TODO add hooks to git repo

	return nil
}

func initGitRepo(repo *repo.Repo) error {
	gitRepoPath := repo.GetGitRepoPath()
	gitRepo, err := git.PlainInit(gitRepoPath, true)
	if err != nil {
		return fmt.Errorf("could not init git repo '%s': %w", repo.GetName(), err)
	}

	err = createRemote(gitRepo, "bridge", repo.GetBridgeRepoPath())
	if err != nil {
		return fmt.Errorf("could not create remote for git repo '%s': %w", repo.GetName(), err)
	}

	return nil
}

func initBridgeRepo(repo *repo.Repo, gitUserName string) error {
	user := store.GetUser(repo, gitUserName)
	bridgeRepoPath := repo.GetBridgeRepoPath()

	gitSvnExecutor := gitsvn.CreateExecutor(user)
	err := gitSvnExecutor.Init(bridgeRepoPath)
	if err != nil {
		return fmt.Errorf("could not init bridge repo '%s': %w", repo.GetName(), err)
	}

	err = gitSvnExecutor.Fetch(bridgeRepoPath)
	if err != nil {
		return fmt.Errorf("could not fetch bridge repo '%s': %w", repo.GetName(), err)
	}

	gitRepo, err := git.PlainOpen(bridgeRepoPath)
	if err != nil {
		return fmt.Errorf("could not open bridge repo '%s': %w", repo.GetName(), err)
	}

	err = createRemote(gitRepo, "git-central-repo", repo.GetGitRepoPath())
	if err != nil {
		return fmt.Errorf("could not create remote for git repo '%s': %w", repo.GetName(), err)
	}

	err = gitRepo.Push(&git.PushOptions{RemoteName: "git-central-repo"})
	if err != nil {
		return fmt.Errorf("could not push from bridge repo '%s': %w", repo.GetName(), err)
	}

	return nil
}

func createRemote(gitRepo *git.Repository, name, path string) error {
	_, err := gitRepo.CreateRemote(&config.RemoteConfig{
		Name: name,
		URLs: []string{path},
	})

	return err
}
