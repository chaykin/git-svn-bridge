package cmd

import (
	"errors"
	"fmt"
	"git-svn-bridge/conf"
	"git-svn-bridge/gitsvn"
	"git-svn-bridge/store"
	"github.com/go-git/go-billy/v5/osfs"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/storage/filesystem"
	"github.com/spf13/cobra"
	"path/filepath"
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

	bridgeRepoPath := filepath.Join(conf.GetConfig().ReposRoot, "bridge", repoName)

	err := initGitRepo(repoName, bridgeRepoPath)
	if err != nil {
		return fmt.Errorf("could not init git repo '%s': %w", repoName, err)
	}
	err = initBridgeRepo(repoName, gitUserName, bridgeRepoPath)
	if err != nil {
		return fmt.Errorf("could not init bridge repo '%s': %w", repoName, err)
	}

	return nil
}

func initGitRepo(repoName, bridgeRepoPath string) error {
	gitRepoPath := filepath.Join(conf.GetConfig().ReposRoot, "git", repoName)

	storage := filesystem.NewStorage(osfs.New(gitRepoPath), nil)
	gitRepo, err := git.Init(storage, nil)
	if err != nil {
		return fmt.Errorf("could not init git repo '%s': %w", repoName, err)
	}

	_, err = gitRepo.CreateRemote(&config.RemoteConfig{
		Name: "bridge",
		URLs: []string{bridgeRepoPath},
	})

	if err != nil {
		return fmt.Errorf("could not create remote for git repo '%s': %w", repoName, err)
	}

	return nil
}

func initBridgeRepo(repoName, gitUserName, bridgeRepoPath string) error {
	repo := store.GetRepo(repoName)
	user := store.GetUser(repo, gitUserName)

	gitSvnExecutor := gitsvn.CreateExecutor(user)
	err := gitSvnExecutor.Init(bridgeRepoPath)
	if err != nil {
		return fmt.Errorf("could not init bridge repo '%s': %w", repoName, err)
	}

	err = gitSvnExecutor.Fetch(bridgeRepoPath)
	if err != nil {
		return fmt.Errorf("could not fetch bridge repo '%s': %w", repoName, err)
	}

	return nil
}
