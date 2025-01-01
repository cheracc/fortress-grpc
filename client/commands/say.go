package commands

import (
	"github.com/cheracc/fortress-grpc"
)

type SayCommand struct {
	SayFunc func(string, string)
}

func (c SayCommand) Execute(player *fortress.Player, args string) (string, error) {
	c.SayFunc(player.GetName(), args)

	return "", nil
}

func (c SayCommand) GetName() string {
	return "say"
}
