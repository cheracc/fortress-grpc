package commands

import (
	"fmt"

	"github.com/cheracc/fortress-grpc"
)

// ListCommand represents a command that lists all online players to the user sending the command
type ListCommand struct {
	// GetOnlinePlayersFunc is the function that will return the list of online players, this is set using CommandHandler.RegisterCommand()
	GetOnlinePlayersFunc func() *[]*fortress.Player
}

// The Execute function which calls GetOnlinePlayersFunc and formats the response
func (c *ListCommand) Execute(player *fortress.Player, args []string) (string, error) {
	playerListString := " "
	var i int = -1
	var p *fortress.Player
	for i, p = range *c.GetOnlinePlayersFunc() {
		name := p.GetName()
		if name == "" {
			name = "no-name"
		}
		if i > 0 {
			playerListString = playerListString + ", "
		}
		playerListString = fmt.Sprintf(playerListString+"%s(%s)", name, p.GetPlayerId())
	}

	var output string
	var err error = nil
	if i >= 0 {
		var areOrIs string
		var pluralS string = ""
		if i == 0 {
			areOrIs = "is"
		} else {
			areOrIs = "are"
			pluralS = "s"
		}
		output = fmt.Sprintf("Online Players:\n\r    %s\n\rThere %s %d player%s online", playerListString, areOrIs, i+1, pluralS)
	} else {
		output = "There are no players online."
		err = fmt.Errorf("no players online but a player requested this")
	}

	return output, err
}
