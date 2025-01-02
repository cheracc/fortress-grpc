package main

import (
	"github.com/cheracc/fortress-grpc"
	"github.com/cheracc/fortress-grpc/server/handlers"
	"github.com/cheracc/fortress-grpc/server/handlers/commands"
)

func main() {
	logger := fortress.NewLogger()
	grpcHandler := handlers.NewGrpcServer(logger)
	sqlite := handlers.NewSqliteHandler(logger)
	playerHandler := handlers.NewPlayerHandler(sqlite, logger)
	auth := handlers.NewAuthHandler(playerHandler, logger)
	playerHandler.SetAuthHandler(auth)
	chat := handlers.NewChatHandler(logger, auth)

	sqlite.InitializeDatabase()

	commandHandler := handlers.NewCommandHandler(auth, playerHandler, logger)
	commandHandler.RegisterCommand(&handlers.Command{Name: "stop", Exec: &commands.StopCommand{CloseDatabaseFunc: sqlite.CloseDb}})
	commandHandler.RegisterCommand(&handlers.Command{Name: "name", Exec: &commands.NameCommand{RenamePlayerFunc: playerHandler.RenamePlayer}})
	commandHandler.RegisterCommand(&handlers.Command{Name: "list", Exec: &commands.ListCommand{GetOnlinePlayersFunc: playerHandler.GetOnlinePlayers}})

	defer sqlite.CloseDb()

	grpcHandler.RegisterHandlers(playerHandler, auth, commandHandler, chat)
	grpcHandler.StartListener()

}
