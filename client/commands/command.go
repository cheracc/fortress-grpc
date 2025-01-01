package commands

import "github.com/cheracc/fortress-grpc"

type Command interface {
	Execute(*fortress.Player, string) (string, error)
	GetName() string
}
