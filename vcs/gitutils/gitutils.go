package gitutils

import (
	"errors"
	"fmt"
	"git-svn-bridge/conf"
	"git-svn-bridge/repo"
	"git-svn-bridge/store"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"os/exec"
	"strings"
)

const GitCentralRepoName = "git-central-repo"

func IsRefExists(gitRepo *git.Repository, refName string) (bool, error) {
	exists := false

	branches, err := gitRepo.Branches()
	if err != nil {
		return false, err
	}

	err = branches.ForEach(func(ref *plumbing.Reference) error {
		if refName == ref.Name().String() {
			exists = true
		}
		return nil
	})
	if err != nil {
		return false, err
	}

	return exists, nil
}

func GetBranchName(ref string) string {
	index := strings.LastIndex(ref, "/")
	return ref[index+1:]
}

func PullAndRebase(repoPath, remote, branch string) error {
	command := fmt.Sprintf("git pull --rebase %s %s", remote, branch)
	return executeCommand(repoPath, command)
}

func GetGitAuthor(repo *repo.Repo, repoPath string) (string, error) {
	commandArgs := strings.Split("git log -n 1 --format=%an", " ")
	systemCmd := exec.Command(commandArgs[0], commandArgs[1:]...)
	systemCmd.Dir = repoPath

	out, err := systemCmd.CombinedOutput()
	if err != nil {
		fmt.Printf("%s\n", out)
		return "", fmt.Errorf("could not get GIT author for repo '%s': %w", repoPath, err)
	}

	userName := strings.TrimRight(string(out), "\n")
	for _, u := range store.GetAllUsers(repo) {
		if u.GetGitUserFullName() == userName {
			return u.GetGitUserName(), nil
		}
	}

	return "", errors.New("could not find user " + userName)
}

func BuildCommitMessage(repoPath, branch string) (string, error) {
	command := fmt.Sprintf("git log --pretty=format:%s HEAD..%s", conf.GetConfig().CommitMessageFormat, branch)
	commandArgs := strings.Split(command, " ")
	systemCmd := exec.Command(commandArgs[0], commandArgs[1:]...)
	systemCmd.Dir = repoPath

	out, err := systemCmd.CombinedOutput()
	if err != nil {
		fmt.Printf("%s\n", out)
		return "", fmt.Errorf("could not build commit message for repo '%s': %w", repoPath, err)
	}

	return strings.TrimRight(string(out), "\n"), nil
}

func MergeNoFF(repoPath, message, branch string) error {
	command := fmt.Sprintf("git merge --no-ff --no-log -m MESSAGE %s", branch)

	commandArgs := strings.Split(command, " ")
	commandArgs[5] = message

	systemCmd := exec.Command(commandArgs[0], commandArgs[1:]...)
	systemCmd.Dir = repoPath

	out, err := systemCmd.CombinedOutput()
	if err != nil {
		fmt.Printf("%s\n", out)
		return fmt.Errorf("could not merge (--no-ff) for repo '%s': %w", repoPath, err)
	}

	return nil
}

func Merge(repoPath, branch string) error {
	if branch == "master" {
		branch = "trunk"
	}

	command := fmt.Sprintf("git merge origin/%s", branch)
	return executeCommand(repoPath, command)
}

func Fetch(repoPath, remote, branch string) error {
	command := fmt.Sprintf("git fetch %s %s:%s", remote, branch, branch)
	return executeCommand(repoPath, command)
}

func executeCommand(cmdDir, cmd string) error {
	commandArgs := strings.Split(cmd, " ")
	systemCmd := exec.Command(commandArgs[0], commandArgs[1:]...)
	if cmdDir != "" {
		systemCmd.Dir = cmdDir
	}

	out, err := systemCmd.CombinedOutput()
	if err != nil {
		fmt.Printf("%s\n", out)
		return fmt.Errorf("could not execute command '%s': %w", cmd, err)
	}

	return nil
}
