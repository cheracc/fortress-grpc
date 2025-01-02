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

const (
	purgeInterval = 1 * time.Minute
)

type ChatHandler struct {
	fgrpc.UnimplementedChatServer
	*AuthHandler
	*fortress.Logger
	*ChatChannel
}

// ChatChannel represents a chat channel. Its methods are thread-safe
type ChatChannel struct {
	*sync.RWMutex
	members []*ChannelMember
}

func (c *ChatChannel) addMember(m *ChannelMember, playerName string) {
	c.sendMessage(fmt.Sprintf("%s has joined chat", playerName), "SERVER") // send this before adding the player so the new player doesn't get it (they already get their own msg)
	c.Lock()
	c.members = append(c.members, m)
	c.Unlock()
}
func (c *ChatChannel) sendMessage(message string, playerName string) {
	c.RLock()
	for _, m := range c.members {
		m.stream.Send(&fgrpc.ChatMessage{Message: message, SendingPlayerName: playerName})
	}
	c.RUnlock()
}
func (c *ChatChannel) removeMember(playerId string) {
	toRemove := -1
	c.Lock()
	for i, m := range c.members {
		if m.playerId == playerId {
			toRemove = i
			m.closed = true
			break
		}
	}
	c.Unlock()

	if toRemove >= 0 {
		c.Lock()
		c.members = slices.Delete(c.members, toRemove, toRemove+1)
		c.Unlock()
		c.sendMessage(fmt.Sprintf("%s has left chat", playerId), "SERVER")
	}
}

type ChannelMember struct {
	playerId string
	stream   fgrpc.Chat_JoinChannelServer
	closed   bool
}

func NewChatHandler(logger *fortress.Logger, auth *AuthHandler) *ChatHandler {
	h := &ChatHandler{fgrpc.UnimplementedChatServer{},
		auth,
		logger,
		&ChatChannel{&sync.RWMutex{}, make([]*ChannelMember, 0)}}

	go func() {
		time.Sleep(20 * time.Second)
		for {
			h.PurgeInactives()
			time.Sleep(purgeInterval)
		}
	}()

	return h
}

func (h *ChatHandler) JoinChannel(req *fgrpc.ChatRequest, stream fgrpc.Chat_JoinChannelServer) error {
	playerId := h.GetPlayerIdFromTokenString(req.GetSessionToken())

	if playerId == "" {
		return h.Error("invalid or expired session token")
	}

	member := &ChannelMember{playerId, stream, false}
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

// PurgeInactives will check each channel member to see if they are still online, then remove those who are not
func (h *ChatHandler) PurgeInactives() {
	inactives := []string{}
	h.RLock()
	for _, m := range h.members {
		if !h.IsOnline(PlayerFilter{playerId: m.playerId}) {
			inactives = append(inactives, m.playerId)
		}
	}
	h.RUnlock()
	for _, id := range inactives {
		h.Logf("Removing %s from chat for inactivity", id)
		h.removeMember(id)
	}
}
