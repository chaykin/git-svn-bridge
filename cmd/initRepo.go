package cmd

import (
	"errors"
	"fmt"
	"git-svn-bridge/conf"
	"git-svn-bridge/repo"
	"git-svn-bridge/store"
	"git-svn-bridge/vcs/gitsvn"
	"git-svn-bridge/vcs/gitutils"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/spf13/cobra"
	"path/filepath"
)

var initRepoCmd = &cobra.Command{
	Use:   "init <repo-name>",
	Short: "Initialize repo",
	Long:  "Initialize repo: create GIT server repo, bridge repo, config hooks, etc",
	Args:  cobra.ExactArgs(1),
	RunE:  initRepo,
}

func init() {
	rootCmd.AddCommand(initRepoCmd)
}

func initRepo(_ *cobra.Command, args []string) error {
	repoName := args[0]
	if !store.HasRepo(repoName) {
		return errors.New("There is no repo with name " + repoName)
	}

	repo := store.GetRepo(repoName)

	err := initGitRepo(&repo)
	if err != nil {
		return fmt.Errorf("could not init git repo '%s': %w", repoName, err)
	}
	err = initBridgeRepo(&repo)
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

func initBridgeRepo(repo *repo.Repo) error {
	systemUser := store.GetUser(repo, conf.GetConfig().SystemGitUserName)
	bridgeRepoPath := repo.GetBridgeRepoPath()

	gitSvnExecutor := gitsvn.CreateExecutor(systemUser)
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

	err = createRemote(gitRepo, gitutils.GitCentralRepoName, repo.GetGitRepoPath())
	if err != nil {
		return fmt.Errorf("could not create remote for git repo '%s': %w", repo.GetName(), err)
	}

	err = gitRepo.Push(&git.PushOptions{RemoteName: gitutils.GitCentralRepoName})
	if err != nil {
		return fmt.Errorf("could not push from bridge repo '%s': %w", repo.GetName(), err)
	}

	return nil
}

func createRemote(gitRepo *git.Repository, name, path string) error {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return fmt.Errorf("could not resolve absolute path for '%s': %w", path, err)
	}

	_, err = gitRepo.CreateRemote(&config.RemoteConfig{
		Name: name,
		URLs: []string{absPath},
	})

	return err
}
