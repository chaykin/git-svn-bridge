package gitsvn

import (
	"bufio"
	"fmt"
	"git-svn-bridge/conf"
	"git-svn-bridge/shell"
	"git-svn-bridge/store"
	"git-svn-bridge/usr"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const commandPattern = "git svn %s --username %s %s"

type CommandExecutor struct {
	user usr.User
}

func CreateExecutor(user usr.User) CommandExecutor {
	return CommandExecutor{user}
}

func (c *CommandExecutor) Init(repoPath string) {
	command := fmt.Sprintf("-t tags -b branches -T trunk %s %s", c.user.GetRepo().GetSvnUrl(), repoPath)

	c.executeCommand("init", command)
}

func (c *CommandExecutor) Fetch(repoPath string) {
	c.createAuthorsFile()
	config := conf.GetConfig()

	authorsFile, err := filepath.Abs(config.AuthorsFile)
	if err != nil {
		panic(fmt.Errorf("could not resolve authors file path '%s': %w", config.AuthorsFile, err))
	}

	command := fmt.Sprintf("--authors-file=%s --log-window-size %d", authorsFile, config.LogWindowsSize)
	c.executeCommandEx(repoPath, "fetch", command)
}

func (c *CommandExecutor) Commit(repoPath string) {
	c.createAuthorsFile()
	config := conf.GetConfig()

	authorsFile, err := filepath.Abs(config.AuthorsFile)
	if err != nil {
		panic(fmt.Errorf("could not resolve authors file path '%s': %w", config.AuthorsFile, err))
	}

	command := fmt.Sprintf("--authors-file=%s", authorsFile)
	c.executeCommandEx(repoPath, "dcommit", command)
}

func (c *CommandExecutor) createAuthorsFile() {
	users := store.GetAllUsers(c.user.GetRepo())

	fileName := conf.GetConfig().AuthorsFile
	authorsFile, err := os.Create(fileName)
	if err != nil {
		panic(fmt.Errorf("could not create authors file '%s': %w", fileName, err))
	}

	writer := bufio.NewWriter(authorsFile)
	defer func() {
		if err := writer.Flush(); err != nil {
			panic(fmt.Errorf("could not flush authors file '%s': %w", fileName, err))
		}

		if err := authorsFile.Close(); err != nil {
			panic(fmt.Errorf("could not close authors file: '%s': %w", fileName, err))
		}
	}()

	for _, user := range users {
		author := fmt.Sprintf("%s = %s <%s>", user.GetSvnUserName(), user.GetGitUserFullName(), user.GetEmail())
		if _, err = fmt.Fprintln(writer, author); err != nil {
			panic(fmt.Errorf("could not write to authors file: '%s': %w", fileName, err))
		}
	}
}

func (c *CommandExecutor) executeCommand(cmdName, cmdArgs string) {
	c.executeCommandEx("", cmdName, cmdArgs)
}

func (c *CommandExecutor) executeCommandEx(cmdDir, cmdName, cmdArgs string) {
	command := fmt.Sprintf(commandPattern, cmdName, c.user.GetSvnUserName(), cmdArgs)
	commandArgs := strings.Split(command, " ")
	systemCmd := exec.Command(commandArgs[0], commandArgs[1:]...)
	if cmdDir != "" {
		systemCmd.Dir = cmdDir
	}

	stdin, err := systemCmd.StdinPipe()
	if err != nil {
		panic(fmt.Errorf("could not pipe stdin for command '%s': %w", command, err))
	}
	defer func() {
		if err := stdin.Close(); err != nil {
			panic(fmt.Errorf("could not close stdin for command '%s': %w", command, err))
		}
	}()

	if _, err = io.WriteString(stdin, c.user.GetSvnPassword()); err != nil {
		panic(fmt.Errorf("could not write to stdin for command '%s': %w", command, err))
	}

	shell.ExecuteCommand(systemCmd)
}
