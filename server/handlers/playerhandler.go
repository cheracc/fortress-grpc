package handlers

import (
	"context"
	"time"

	"github.com/cheracc/fortress-grpc"
	fgrpc "github.com/cheracc/fortress-grpc/grpc"
	"github.com/google/uuid"
)

type PlayerHandler struct {
	fgrpc.UnimplementedPlayerServer
	*SqliteHandler
	*AuthHandler
	*fortress.Logger
	players []*fortress.Player
}

func NewPlayerHandler(sqliteHandler *SqliteHandler, logger *fortress.Logger) *PlayerHandler {
	handler := &PlayerHandler{fgrpc.UnimplementedPlayerServer{}, sqliteHandler, nil, logger, []*fortress.Player{}}

	go func() {
		time.Sleep(1 * time.Minute)
		for {
			handler.SaveAllPlayersToDb()
			time.Sleep(time.Minute)
		}
	}()
	return handler
}

func (h *PlayerHandler) SetAuthHandler(auth *AuthHandler) {
	h.AuthHandler = auth
}

func (h *PlayerHandler) GetPlayerData(ctx context.Context, playerInfo *fgrpc.PlayerInfo) (*fgrpc.PlayerMessage, error) {
	tokenString := playerInfo.SessionToken

	id := h.GetPlayerIdFromTokenString(tokenString)
	if id == "" {
		return &fgrpc.PlayerMessage{}, h.Errorf("Server was unable to verify or decode session token %s", tokenString)
	}
	// if reaching this point, the session token sent is valid and we can send the requested player data
	p, _ := h.GetPlayer(PlayerFilter{playerId: id}, true)

	payload := &fgrpc.PlayerMessage{PlayerId: p.GetPlayerId(), Name: p.GetName(), CreatedAt: p.GetCreatedAt().Unix()}
	return payload, nil
}

func (h *PlayerHandler) GetOnlinePlayers() *[]*fortress.Player {
	return &h.players
}

// this error gets passed back to the user/client
func (h *PlayerHandler) RenamePlayer(player *fortress.Player, newName string) error {
	if !h.SqliteHandler.IsNameUnique(newName) {
		return h.Errorf("the name %s name is already in use", newName)
	}

	// do some other checking - like against a list of banned words, limit characters allowed, etc

	player.SetName(newName)
	return nil
}

// GetPlayer checks currently online players for one that matches any field in the provided playerfilter. If checkDatabase is true, it will also check
// the sqlite database using the same filter. If it does not find a player, it creates a new one. If it either loaded from the database
// or created it, it adds that player to onlinePlayers. it returns the player as well as whether that player was newly created
func (h *PlayerHandler) GetPlayer(filter PlayerFilter, checkDatabase bool) (*fortress.Player, bool) {
	for _, p := range *h.GetOnlinePlayers() {
		if p.GetPlayerId() == filter.playerId || p.GetName() == filter.name || p.GetGoogleId() == filter.googleId {
			return p, false
		}
	}
	// check the db
	if checkDatabase {
		p := h.SqliteHandler.LookupPlayerFromDb(filter)
		if p != nil {
			h.updateOnlinePlayer(p)
			h.Logf("Loaded player %s(%s) from database.", p.GetName(), p.GetPlayerId())
			return p, false
		}
	}

	//not found, create new
	p := fortress.NewPlayer()

	if _, err := uuid.Parse(filter.playerId); err == nil { // if the requested id was a uuid, set the new player to that id, otherwise, keep the generated one
		p.SetPlayerId(filter.playerId)
	}
	h.registerNewPlayer(p, true)
	h.Logf("Created a new player with ID %s", p.GetPlayerId())
	return p, true

}

func (h *PlayerHandler) registerNewPlayer(player *fortress.Player, saveToDb bool) {
	h.updateOnlinePlayer(player)
	if saveToDb {
		h.SqliteHandler.CreateNewPlayerDbRecord(player)
	}
}

func (h *PlayerHandler) SaveAllPlayersToDb() {
	for _, p := range *h.GetOnlinePlayers() {
		if time.Since(p.UpdatedAt) < 1*time.Minute {
			h.SqliteHandler.UpdatePlayerToDb(p)
		}
	}
}

func (h *PlayerHandler) updateOnlinePlayer(player *fortress.Player) {
	var toRemove int = -1
	for i, p := range h.players {
		if p.GetPlayerId() == player.GetPlayerId() {
			toRemove = i
			break
		}
	}
	if toRemove >= 0 {
		h.players = append(h.players[:toRemove], h.players[toRemove+1:]...) // remove the player if it was found
		h.Logf("Updating online player: %s(%s)", player.GetName(), player.GetPlayerId())
	} else {
		h.Logf("Loaded player: %s(%s)", player.GetName(), player.GetPlayerId())
	}
	h.players = append(h.players, player)
}

// GetPlayerNameFromId returns the name associated with the given playerId, optionally checking the database. Returns "" if no player is found
func (h *PlayerHandler) GetPlayerNameFromId(playerId string, checkDatabase bool) string {
	p, _ := h.GetPlayer(PlayerFilter{playerId: playerId}, checkDatabase)
	return p.GetName()
}
