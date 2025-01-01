package commands

import (
	"fmt"

	"github.com/cheracc/fortress-grpc"
)

// NameCommand represents a command which can be used to change a player's name
type NameCommand struct {
	// RenamePlayerFunc is the function that will be used to rename the player
	RenamePlayerFunc func(*fortress.Player, string) error
}

// Execute does name verification and calls RenamePlayerFunc
func (c *NameCommand) Execute(player *fortress.Player, args []string) (string, error) {
	if len(args) <= 0 {
		return "", fmt.Errorf("tried to rename but no name given")
	}
	if len(args) > 1 {
		return "", fmt.Errorf("too many arguments. Syntax: name <new name>")
	}
	c.RenamePlayerFunc(player, args[0])
	return "", nil
}
