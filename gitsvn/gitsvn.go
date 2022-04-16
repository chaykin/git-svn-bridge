package gitsvn

import (
	"bufio"
	"fmt"
	"git-svn-bridge/conf"
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

func (c *CommandExecutor) Init(repoPath string) error {
	userRepo := c.user.GetRepo()
	command := fmt.Sprintf("-t tags -b branches -T trunk %s %s", userRepo.GetSvnUrl(), repoPath)

	return c.executeCommand("init", command)
}

func (c *CommandExecutor) Fetch(repoPath string) error {
	err := c.createAuthorsFile()
	if err != nil {
		return err
	}

	config := conf.GetConfig()

	authorsFile, err := getAuthorsFilePath(config)
	if err != nil {
		return fmt.Errorf("could not resolve authors file path: %w", err)
	}

	command := fmt.Sprintf("--authors-file=%s --log-window-size %d", authorsFile, config.LogWindowsSize)
	return c.executeCommandEx(repoPath, "fetch", command)
}

func (c *CommandExecutor) createAuthorsFile() error {
	userRepo := c.user.GetRepo()
	users := store.GetAllUsers(userRepo)

	authorsFile, err := os.Create(conf.GetConfig().AuthorsFile)
	if err != nil {
		return fmt.Errorf("could not create authors file: %w", err)
	}

	writer := bufio.NewWriter(authorsFile)
	for _, user := range users {
		author := fmt.Sprintf("%s = %s <%s>", user.GetSvnUserName(), user.GetGitUserFullName(), user.GetEmail())
		_, err = fmt.Fprintln(writer, author)
		if err != nil {
			return fmt.Errorf("could not write to authors file: %w", err)
		}
	}

	err = writer.Flush()
	if err != nil {
		return fmt.Errorf("could not flush authors file: %w", err)
	}

	err = authorsFile.Close()
	if err != nil {
		return fmt.Errorf("could not close authors file: %w", err)
	}

	return nil
}

func getAuthorsFilePath(config *conf.Config) (string, error) {
	authorsFile := config.AuthorsFile
	if !filepath.IsAbs(authorsFile) {
		wd, err := os.Getwd()
		fmt.Println(wd)
		if err != nil {
			return "", err
		}
		authorsFile = filepath.Join(wd, authorsFile)
	}
	return authorsFile, nil
}

func (c *CommandExecutor) executeCommand(cmdName, cmdArgs string) error {
	return c.executeCommandEx("", cmdName, cmdArgs)
}

func (c *CommandExecutor) executeCommandEx(cmdDir, cmdName, cmdArgs string) error {
	userRepo := c.user.GetRepo()
	repoName := userRepo.GetName()

	command := fmt.Sprintf(commandPattern, cmdName, c.user.GetSvnUserName(), cmdArgs)
	commandArgs := strings.Split(command, " ")
	systemCmd := exec.Command(commandArgs[0], commandArgs[1:]...)
	if cmdDir != "" {
		systemCmd.Dir = cmdDir
	}

	stdin, err := systemCmd.StdinPipe()
	if err != nil {
		return fmt.Errorf("could not pipe stdin for repo '%s': %w", repoName, err)
	}

	_, err = io.WriteString(stdin, c.user.GetSvnPassword())
	if err != nil {
		return fmt.Errorf("could not write to stdin for repo '%s': %w", repoName, err)
	}
	err = stdin.Close()
	if err != nil {
		return fmt.Errorf("could not close stdin for repo '%s': %w", repoName, err)
	}

	out, err := systemCmd.CombinedOutput()
	if err != nil {
		fmt.Printf("%s\n", out)
		return fmt.Errorf("could not execute command '%s' for repo '%s': %w", cmdName, repoName, err)
	}

	return nil
}
