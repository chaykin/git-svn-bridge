package svn

import (
	"fmt"
	"git-svn-bridge/shell"
	"git-svn-bridge/usr"
	"os/exec"
	"strings"
)

const commandPattern = "svn --username %s --password %s %s %s"

type CommandExecutor struct {
	user usr.User
}

func CreateExecutor(user usr.User) CommandExecutor {
	return CommandExecutor{user}
}

func (c *CommandExecutor) Branches() []string {
	command := fmt.Sprintf("%s/branches", c.user.GetRepo().GetSvnUrl())

	result := c.executeCommand("ls", command)
	return strings.Split(result, "/")
}

func (c *CommandExecutor) Tags() []string {
	command := fmt.Sprintf("%s/tags", c.user.GetRepo().GetSvnUrl())

	result := c.executeCommand("ls", command)
	return strings.Split(result, "\\n")
}

func (c *CommandExecutor) executeCommand(cmdName, cmdArgs string) string {
	command := fmt.Sprintf(commandPattern, c.user.GetSvnUserName(), c.user.GetSvnPassword(), cmdName, cmdArgs)
	commandArgs := strings.Split(command, " ")
	systemCmd := exec.Command(commandArgs[0], commandArgs[1:]...)

	return shell.ExecuteCommand(systemCmd)
}
