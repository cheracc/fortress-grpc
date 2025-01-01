package main

import (
	"bufio"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/cheracc/fortress-grpc"
	"github.com/cheracc/fortress-grpc/client/commands"
	"github.com/cheracc/fortress-grpc/client/handlers"
)

func main() {
	logger := fortress.NewLogger()
	remote := handlers.NewRemote(logger)
	cmd := handlers.NewCommandHandler(logger)
	chat := handlers.NewChatHandler(remote)

	cmd.RegisterCommand(commands.LogoutCommand{LogoutFunc: remote.Logout})
	cmd.RegisterCommand(commands.SayCommand{SayFunc: chat.SendChatMessageToServer})
	cmd.RegisterCommand(commands.QuitCommand{})

	go refreshTokenEveryMinute(remote)
	go joinChatOnceLoggedIn(remote, chat)

	inputReader := bufio.NewReader(os.Stdin)
	for {
		line, _ := inputReader.ReadString('\n')
		line = trimString(line)
		cmdName, args, _ := strings.Cut(line, " ")
		if cmdName == "" {
			continue
		}

		if c := cmd.GetCommandOrNil(cmdName); c != nil {
			c.Execute(remote.Player, args)
			continue
		}

		response := remote.SendCommand(cmdName, args)
		if response != "" {
			logger.ToConsole(response)
		}
	}
}

func joinChatOnceLoggedIn(remote *handlers.Remote, chat *handlers.Chat) {
	time.Sleep(10 * time.Second)
	for {
		if remote.HasSessionToken() {
			if !chat.HasOpenChannel() {
				chat.JoinChat()
				continue
			}
			break
		}
	}
}

func refreshTokenEveryMinute(remote *handlers.Remote) {
	for {
		remote.Authorize()
		if !remote.HasSessionToken() {
			time.Sleep(5 * time.Second)
		} else {
			time.Sleep(1 * time.Minute)
		}
	}
}

func trimString(s string) string {
	if runtime.GOOS == "windows" {
		return strings.TrimRight(s, "\r\n")
	} else {
		return strings.TrimRight(s, "\n")
	}
}
