package commands

import (
	"os"
	"time"

	"github.com/cheracc/fortress-grpc"
)

type QuitCommand struct {
}

func (c QuitCommand) Execute(player *fortress.Player, args string) (string, error) {
	defer func() {
		time.Sleep(1 * time.Millisecond)
		os.Exit(0)
	}()
	return "", nil
}

func (c QuitCommand) GetName() string {
	return "quit"
}
