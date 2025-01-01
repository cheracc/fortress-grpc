package handlers

import (
	"context"
	"strings"

	"github.com/cheracc/fortress-grpc"
	fgrpc "github.com/cheracc/fortress-grpc/grpc"
	"github.com/cheracc/fortress-grpc/server/handlers/commands"
)

// A CommandHandler handles commands sent to the server by a user. It holds references to the AuthHandler and PlayerHandler, as well as the Logger
type CommandHandler struct {
	// the gRPC server that handles commands
	fgrpc.UnimplementedCommandServer
	// the array of pointers to the registered commands
	commands []*Command
	// the authorization handler
	*AuthHandler
	// the player handler
	*PlayerHandler
	// the logger
	*fortress.Logger
}

// A Command
type Command struct {
	// the name of the command
	Name string
	// the actual command interface
	Exec commands.Executable
}

// NewCommandHandler constructs a new command handler
func NewCommandHandler(auth *AuthHandler, playerHandler *PlayerHandler, logger *fortress.Logger) *CommandHandler {
	handler := &CommandHandler{fgrpc.UnimplementedCommandServer{}, make([]*Command, 0), auth, playerHandler, logger}
	return handler
}

// RegisterCommand adds the command to CommandHandler.commands
func (h *CommandHandler) RegisterCommand(command *Command) {
	h.commands = append(h.commands, command)
	h.Logf("Registered command %s", command.Name)
}

// receives commands from remote users and sends response as a json payload
func (h *CommandHandler) Command(ctx context.Context, commandInfo *fgrpc.CommandInfo) (*fgrpc.CommandReturn, error) {
	h.Logf("Loaded commands: %v", h.commands)
	h.Logf("Received command execution request %s %s from player %s", commandInfo.GetCommandName(), commandInfo.GetCommandArguments(), commandInfo.GetPlayerInfo().GetId())

	c := h.lookupCommand(commandInfo.GetCommandName())
	if c == nil {
		return &fgrpc.CommandReturn{Success: false, JsonPayload: "command not recognized: %s" + commandInfo.GetCommandName()}, h.Errorf("no command found: %s", commandInfo.GetCommandName())
	}

	// found the command, validate the user
	tokenString := commandInfo.GetPlayerInfo().GetSessionToken()

	if !h.AuthHandler.IsValidToken(tokenString) {
		return &fgrpc.CommandReturn{Success: false, JsonPayload: ""}, h.Error("invalid session token")
	}
	playerId := h.AuthHandler.GetPlayerIdFromTokenString(tokenString)
	if playerId != commandInfo.GetPlayerInfo().GetId() {
		h.Logf("session token id: %s, sent id: %s", playerId, commandInfo.GetPlayerInfo().GetId())
		return &fgrpc.CommandReturn{Success: false, JsonPayload: ""}, h.Error("session token does not match sending player")
	}

	player, _ := h.PlayerHandler.GetPlayerByID(playerId, false)
	if player == nil {
		return &fgrpc.CommandReturn{Success: false, JsonPayload: ""}, h.Errorf("could not find an online player with id %s", playerId)
	}

	h.Logf("Executing command %s for player %s(%s)", c.Name, player.GetName(), player.GetPlayerId())
	response, err := c.Execute(player, getArgs(commandInfo))

	return &fgrpc.CommandReturn{Success: err != nil, JsonPayload: response}, err
}

// lookupCommand fetches the command with the given name if it exists
func (h *CommandHandler) lookupCommand(name string) *Command {
	h.Logf("Looking up command %s", name)
	for _, c := range h.commands {
		h.Logf("Checking command %s", c.Name)
		if c.Name == name {
			return c
		}
	}
	return nil
}

// getArgs extracts the arguments from the sent CommandInfo as a string slice
func getArgs(ci *fgrpc.CommandInfo) []string {
	args := strings.Split(ci.CommandArguments, " ")

	return args
}

// Execute calls the function contained in the command
func (c *Command) Execute(player *fortress.Player, args []string) (string, error) {
	s, err := c.Exec.Execute(player, args)
	if err != nil {
		return "", err
	}

	return s, nil
}
