package commands

import "github.com/cheracc/fortress-grpc"

// An Executable is used to execute commands sent to the server by a user
type Executable interface {
	// Execute accepts a Player and a string array of optional arguments
	// Execute then returns a response as a string and an error
	Execute(*fortress.Player, []string) (string, error)
}
