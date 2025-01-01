package commands

import (
	"os"
	"time"

	"github.com/cheracc/fortress-grpc"
)

// StopCommand represents a command used to stop the server
type StopCommand struct {
	// CloseDatabaseFunc is used to close the sqlite database before stopping the server
	CloseDatabaseFunc func() error
}

// Execute closes the sqlite database and stops the server (gracefully?)
func (c *StopCommand) Execute(player *fortress.Player, args []string) (string, error) {
	c.CloseDatabaseFunc()
	go func() {
		time.Sleep(time.Second)
		os.Exit(0)
	}()
	return "", nil
}
