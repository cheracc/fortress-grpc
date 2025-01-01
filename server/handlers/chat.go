package handlers

import (
	"context"
	"fmt"
	"time"

	"github.com/cheracc/fortress-grpc"
	fgrpc "github.com/cheracc/fortress-grpc/grpc"
)

type ChatHandler struct {
	fgrpc.UnimplementedChatServer
	*AuthHandler
	*fortress.Logger
	ChatChannel
}

type ChatChannel struct {
	members []ChannelMember
}

func (c *ChatChannel) addMember(m ChannelMember, playerName string) {
	c.sendMessage(fmt.Sprintf("%s has joined chat", playerName), "SERVER")
	c.members = append(c.members, m)
}
func (c *ChatChannel) sendMessage(message string, playerName string) {
	for _, m := range c.members {
		m.stream.Send(&fgrpc.ChatMessage{Message: message, SendingPlayerName: playerName})
	}
}

type ChannelMember struct {
	playerId string
	stream   fgrpc.Chat_JoinChannelServer
	closed   bool
}

func NewChatHandler(logger *fortress.Logger, auth *AuthHandler) *ChatHandler {
	h := &ChatHandler{fgrpc.UnimplementedChatServer{}, auth, logger, ChatChannel{[]ChannelMember{}}}

	return h
}

func (h *ChatHandler) JoinChannel(req *fgrpc.ChatRequest, stream fgrpc.Chat_JoinChannelServer) error {
	playerId := h.GetPlayerIdFromTokenString(req.GetSessionToken())

	if playerId == "" {
		return h.Error("invalid or expired session token")
	}

	member := ChannelMember{playerId, stream, false}
	h.addMember(member, h.GetPlayerNameFromId(playerId, false))
	for {
		if member.closed {
			break
		}
		time.Sleep(2 * time.Second)
	}
	return nil
}

// ChatHandler.SendMessage is the gRPC server function that receives chat messages from clients
func (h *ChatHandler) SendMessage(ctx context.Context, msg *fgrpc.ChatMessage) (*fgrpc.Empty, error) {
	playerId := h.GetPlayerIdFromTokenString(msg.GetSessionToken())

	if playerId == "" {
		return nil, h.Errorf("invalid or expired session token")
	}

	h.ChatChannel.sendMessage(msg.GetMessage(), h.GetPlayerNameFromId(playerId, false))
	return nil, nil
}
