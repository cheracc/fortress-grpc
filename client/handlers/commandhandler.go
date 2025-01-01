package handlers

import (
	"github.com/cheracc/fortress-grpc"
	"github.com/cheracc/fortress-grpc/client/commands"
)

type CommandHandler struct {
	*fortress.Logger
	commands []commands.Command
}

func NewCommandHandler(logger *fortress.Logger) *CommandHandler {
	return &CommandHandler{logger, make([]commands.Command, 0)}
}

func (h *CommandHandler) GetCommandOrNil(commandName string) commands.Command {
	for _, c := range h.commands {
		if c.GetName() == commandName {
			return c
		}
	}
	return nil
}

func (h *CommandHandler) RegisterCommand(cmd commands.Command) {
	h.commands = append(h.commands, cmd)
}
