package plock

import (
	"git-svn-bridge/log"
	"github.com/juju/fslock"
)

func Lock() {
	lock := fslock.New("pid.lock")
	if err := lock.Lock(); err != nil {
		log.Fatalf("could not acquire a lock: %w", err)
	}
}
