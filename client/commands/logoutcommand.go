package commands

import (
	"fmt"

	"github.com/cheracc/fortress-grpc"
)

type LogoutCommand struct {
	LogoutFunc func()
}

func (c LogoutCommand) Execute(player *fortress.Player, _ string) (string, error) {
	if c.LogoutFunc == nil {
		return "", fmt.Errorf("logout function is not defined")
	}
	name := player.GetName()
	c.LogoutFunc()

	return fmt.Sprintf("Logged out %s", name), nil
}

func (c LogoutCommand) GetName() string {
	return "logout"
}
