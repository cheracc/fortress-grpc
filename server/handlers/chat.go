package handlers

import (
	"context"
	"fmt"
	"slices"
	"sync"
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
	*sync.Mutex
	members []ChannelMember
}

func (c *ChatChannel) GetMembers() *[]ChannelMember {
	return &c.members
}
func (c *ChatChannel) addMember(m *ChannelMember, playerName string) {
	c.sendMessage(fmt.Sprintf("%s has joined chat", playerName), "SERVER")
	c.members = append(c.members, *m)
}
func (c *ChatChannel) sendMessage(message string, playerName string) {
	for _, m := range *c.GetMembers() {
		m.stream.Send(&fgrpc.ChatMessage{Message: message, SendingPlayerName: playerName})
	}
}
func (c *ChatChannel) RemoveMember(index int) {
	c.members = slices.Delete(c.members, index, index+1)
}

type ChannelMember struct {
	*sync.Mutex
	playerId string
	stream   fgrpc.Chat_JoinChannelServer
	closed   bool
}

func (m *ChannelMember) SetClosed() {
	m.closed = true
}
func (m *ChannelMember) IsClosed() bool {
	return m.closed
}

func NewChatHandler(logger *fortress.Logger, auth *AuthHandler) *ChatHandler {
	h := &ChatHandler{fgrpc.UnimplementedChatServer{}, auth, logger, ChatChannel{&sync.Mutex{}, []ChannelMember{}}}

	go func() {
		for {
			toRemove := []int{}
			for i, m := range *h.GetMembers() {
				if !h.IsOnline(m.playerId) {
					m.SetClosed()
					toRemove = append(toRemove, i)
				}
			}
			time.Sleep(3 * time.Second) // if changing this, also change the sleep time in JoinChannel

			for _, i := range toRemove {
				h.RemoveMember(i)
			}
		}
	}()

	return h
}

func (h *ChatHandler) JoinChannel(req *fgrpc.ChatRequest, stream fgrpc.Chat_JoinChannelServer) error {
	playerId := h.GetPlayerIdFromTokenString(req.GetSessionToken())

	if playerId == "" {
		return h.Error("invalid or expired session token")
	}

	member := &ChannelMember{&sync.Mutex{}, playerId, stream, false}
	h.addMember(member, h.GetPlayerNameFromId(playerId, false))
	for {
		if member.IsClosed() {
			h.Logf("Closing chat stream for player %s due to inactivity", playerId)
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

	h.Logf("Sending chat message: %s from player %s", msg.GetMessage(), playerId)
	h.ChatChannel.sendMessage(msg.GetMessage(), h.GetPlayerNameFromId(playerId, false))
	return nil, nil
}
