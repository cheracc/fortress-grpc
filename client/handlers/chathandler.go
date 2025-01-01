package handlers

import (
	"context"
	"io"
	"time"

	fgrpc "github.com/cheracc/fortress-grpc/grpc"
)

const chatMonitorInterval = time.Second / 2

type Chat struct {
	*Remote
	*ChatStream
}

type ChatStream struct {
	fgrpc.Chat_JoinChannelClient
	ContextCancelFunc func()
	monitorCancelFunc func()
}

func (c *Chat) SendChatMessageToServer(playerName string, message string) {
	chatMessage := &fgrpc.ChatMessage{SendingPlayerName: playerName, Message: message, SessionToken: c.GetSessionToken()}

	_, err := c.SendMessage(context.Background(), chatMessage)
	if err != nil {
		c.Error(err.Error())
	}
}

func (c *Chat) PostMessageToConsole(message *fgrpc.ChatMessage) {
	if message != nil {
		sender := message.GetSendingPlayerName()
		msg := message.GetMessage()

		if msg == "" {
			return
		}

		c.ToConsolef("[CHAT] %s: %s", sender, msg)
	}
}

func (s *Chat) StartChannelMonitor() func() {
	ctx, cancel := context.WithCancel(context.Background())

	go func(ctx context.Context) {
		for {
			select {
			case <-ctx.Done():
				// cancel this jawn
				s.Log("Closing chat stream")
				return
			default:
				message, err := s.Recv()
				if err == io.EOF {
					s.Log("Chat stream closed by remote server")
					return
				}
				if err != nil {
					s.Error(err.Error())
				}
				s.PostMessageToConsole(message)
				time.Sleep(chatMonitorInterval)
			}
		}
	}(ctx)

	return cancel
}
func NewChatHandler(remote *Remote) *Chat {
	return &Chat{remote, nil}
}

func (c *Chat) JoinChat() {
	if c.ChatStream != nil {
		c.Error("tried to join a chat stream but we already have one")
		return
	}
	c.ChatStream = c.GetChatChannel()
	c.monitorCancelFunc = c.StartChannelMonitor()
	c.Log("Joined chat.")
}

func (c *Chat) HasOpenChannel() bool {
	return c.ChatStream != nil
}
