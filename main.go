package main

import (
	"git-svn-bridge/cmd"
	"git-svn-bridge/log"
	"git-svn-bridge/plock"
)

func main() {
	log.InitLogging()
	defer log.CloseLog()

	plock.Lock()
	cmd.Execute()
}
