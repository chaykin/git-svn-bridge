package cmd

import (
	"fmt"
	"git-svn-bridge/conf"
	"git-svn-bridge/gitHook"
	"git-svn-bridge/log"
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
	Run:   initRepo,
}

func init() {
	rootCmd.AddCommand(initRepoCmd)
}

func initRepo(_ *cobra.Command, args []string) {
	repoName := args[0]

	defer log.StdErrOnPanicf(fmt.Errorf("could not init repository '%s'", repoName))

	if !store.HasRepo(repoName) {
		panic(fmt.Errorf("there is no repository with name %s", repoName))
	}

	repository := store.GetRepo(repoName)
	initGitRepo(&repository)
	initBridgeRepo(&repository)

	hookFactory := gitHook.New(&repository)
	hookFactory.CreateUpdateHook()
	//TODO add post-update hook to git repository
}

func initGitRepo(repo *repo.Repo) {
	gitRepoPath := repo.GetGitRepoPath()
	gitRepo, err := git.PlainInit(gitRepoPath, true)
	if err != nil {
		panic(fmt.Errorf("could not init git repo '%s': %w", gitRepoPath, err))
	}

	createRemote(gitRepo, "bridge", repo.GetBridgeRepoPath())
}

func initBridgeRepo(repo *repo.Repo) {
	systemUser := store.GetUser(repo, conf.GetConfig().SystemGitUserName)
	bridgeRepoPath := repo.GetBridgeRepoPath()

	gitSvnExecutor := gitsvn.CreateExecutor(systemUser)
	gitSvnExecutor.Init(bridgeRepoPath)
	gitSvnExecutor.Fetch(bridgeRepoPath)

	gitRepo, err := git.PlainOpen(bridgeRepoPath)
	if err != nil {
		panic(fmt.Errorf("could not open bridge repo '%s': %w", bridgeRepoPath, err))
	}

	createRemote(gitRepo, gitutils.GitCentralRepoName, repo.GetGitRepoPath())

	if err := gitRepo.Push(&git.PushOptions{RemoteName: gitutils.GitCentralRepoName}); err != nil {
		panic(fmt.Errorf("could not push bridge repo '%s': %w", bridgeRepoPath, err))
	}
}

func createRemote(gitRepo *git.Repository, name, path string) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		panic(fmt.Errorf("could not resolve absolute path for '%s': %w", path, err))
	}

	_, err = gitRepo.CreateRemote(&config.RemoteConfig{
		Name: name,
		URLs: []string{absPath},
	})

	if err != nil {
		panic(fmt.Errorf("could not create remote '%s': %w", path, err))
	}
}
