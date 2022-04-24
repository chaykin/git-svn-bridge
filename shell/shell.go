package shell

import (
	"fmt"
	"os/exec"
	"strconv"
)

func ExecuteCommand(cmd *exec.Cmd) string {
	out, err := cmd.CombinedOutput()
	if err != nil {
		exitError, ok := err.(*exec.ExitError)
		exitCode := "?"
		if ok {
			exitCode = strconv.Itoa(exitError.ExitCode())
		}

		//TODO cmd.String()???
		panic(fmt.Errorf("could not execute command '%s': %w.\nExit code: %s\nCommand output:\n%s", cmd, err, exitCode, out))
	}

	return string(out)
}
