package plock

import (
	"fmt"
	"github.com/juju/fslock"
)

func Lock() {
	lock := fslock.New("pid.lock")
	err := lock.Lock()
	if err != nil {
		panic(fmt.Errorf("could not acquire a lock: %w", err))
	}
}
