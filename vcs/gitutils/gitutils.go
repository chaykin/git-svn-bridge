package gitutils

import (
	"fmt"
	"git-svn-bridge/conf"
	"git-svn-bridge/repo"
	"git-svn-bridge/shell"
	"git-svn-bridge/store"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"os/exec"
	"strconv"
	"strings"
)

const GitCentralRepoName = "git-central-repo"

func IsRefExists(gitRepo *git.Repository, refName string) bool {
	exists := false

	branches, err := gitRepo.Branches()
	if err != nil {
		panic(fmt.Errorf("could not list branches: %w", err))
	}

	err = branches.ForEach(func(ref *plumbing.Reference) error {
		if refName == ref.Name().String() {
			exists = true
		}
		return nil
	})
	if err != nil {
		panic(fmt.Errorf("an error occurred while getting branches: %w", err))
	}

	return exists
}

func GetBranchName(ref string) string {
	index := strings.LastIndex(ref, "/")
	name := ref[index+1:]
	if strings.Contains(ref, "/tags/") {
		name = "tags/" + name
	}
	return name
}

func PullAndRebase(repoPath, remote, branch string) {
	command := fmt.Sprintf("git pull --rebase %s %s", remote, branch)
	executeCommand(repoPath, command)
}

func GetGitAuthor(repo *repo.Repo, repoPath string) string {
	command := "git log -n 1 --format=%an"
	commandArgs := strings.Split(command, " ")
	systemCmd := exec.Command(commandArgs[0], commandArgs[1:]...)
	systemCmd.Dir = repoPath

	out, err := systemCmd.CombinedOutput()
	if err != nil {
		exitError, ok := err.(*exec.ExitError)
		exitCode := "?"
		if ok {
			exitCode = strconv.Itoa(exitError.ExitCode())
		}

		panic(fmt.Errorf("could not execute command '%s': %w.\nExit code: %s\nCommand output:\n%s", command, err, exitCode, out))
	}

	userName := strings.TrimRight(string(out), "\n")
	for _, u := range store.GetAllUsers(repo) {
		if u.GetGitUserName() == userName {
			return u.GetGitUserName()
		}
	}

	panic(fmt.Errorf("could not find user '%s' for repo '%s'", userName))
}

func BuildCommitMessage(repoPath, branch string) string {
	command := fmt.Sprintf("git log --pretty=format:%s HEAD..%s", conf.GetConfig().CommitMessageFormat, branch)
	commandArgs := strings.Split(command, " ")
	systemCmd := exec.Command(commandArgs[0], commandArgs[1:]...)
	systemCmd.Dir = repoPath

	out := shell.ExecuteCommand(systemCmd)
	return strings.TrimRight(string(out), "\n")
}

func MergeNoFF(repoPath, message, branch string) {
	command := fmt.Sprintf("git merge --no-ff --no-log -m MESSAGE %s", branch)

	commandArgs := strings.Split(command, " ")
	commandArgs[5] = message

	systemCmd := exec.Command(commandArgs[0], commandArgs[1:]...)
	systemCmd.Dir = repoPath

	shell.ExecuteCommand(systemCmd)
}

func Merge(repoPath, branch string) {
	if branch == "master" {
		branch = "trunk"
	}

	command := fmt.Sprintf("git merge origin/%s", branch)
	executeCommand(repoPath, command)
}

func GetMergeBase(repoPath, oldSha, newSha string) string {
	command := fmt.Sprintf("git merge-base %s %s", newSha, oldSha)
	return executeCommand(repoPath, command)
}

func AbortMerge(repoPath string) {
	executeCommand(repoPath, "git merge --abort")
}

func AbortRebase(repoPath string) {
	executeCommand(repoPath, "git rebase --abort")
}

func Fetch(repoPath, remote, branch string) {
	command := fmt.Sprintf("git fetch %s %s:%s", remote, branch, branch)
	executeCommand(repoPath, command)
}

func RemoveBranch(gitRepo *git.Repository, repoPath, ref string) {
	refExists := IsRefExists(gitRepo, ref)
	if refExists {
		branchName := GetBranchName(ref)
		command := fmt.Sprintf("git branch -D %s", branchName)
		executeCommand(repoPath, command)
	}
}

func executeCommand(cmdDir, cmd string) string {
	commandArgs := strings.Split(cmd, " ")
	systemCmd := exec.Command(commandArgs[0], commandArgs[1:]...)
	if cmdDir != "" {
		systemCmd.Dir = cmdDir
	}

	return shell.ExecuteCommand(systemCmd)
}
