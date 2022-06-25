package shell

import (
	"bufio"
	"bytes"
	"fmt"
	"os/exec"
	"strconv"
)

func ExecuteCommand(cmd *exec.Cmd) string {
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		panic(fmt.Errorf("could not get stdout pipe for command '%s': %w", cmd, err))
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		panic(fmt.Errorf("could not get stderr pipe for command '%s': %w", cmd, err))
	}

	if err := cmd.Start(); err != nil {
		panic(fmt.Errorf("could not start execute command '%s': %w", cmd, err))
	}

	var combinedBuf bytes.Buffer
	stdOutDoneChan := make(chan bool)
	stdErrDoneChan := make(chan bool)
	go readOutput(bufio.NewScanner(stdout), &combinedBuf, stdOutDoneChan)
	go readOutput(bufio.NewScanner(stderr), &combinedBuf, stdErrDoneChan)

	<-stdOutDoneChan
	<-stdErrDoneChan

	if err := cmd.Wait(); err != nil {
		exitError, ok := err.(*exec.ExitError)
		exitCode := "?"
		if ok {
			exitCode = strconv.Itoa(exitError.ExitCode())
		}

		panic(fmt.Errorf("could not execute command '%s': %w.\nExit code: %s\nCommand output:\n%s", cmd, err, exitCode, combinedBuf.String()))
	}

	return combinedBuf.String()
}

func readOutput(scanner *bufio.Scanner, buf *bytes.Buffer, doneChan chan<- bool) {
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		l := scanner.Text()
		buf.WriteString(l)
		fmt.Println(l)
	}

	doneChan <- true
}
