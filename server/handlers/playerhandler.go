package handlers

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/cheracc/fortress-grpc"
	fgrpc "github.com/cheracc/fortress-grpc/grpc"
	"github.com/google/uuid"
)

const (
	dbSaveInterval = 1 * time.Minute
	inactiveTime   = 3 * time.Minute
)

type PlayerHandler struct {
	fgrpc.UnimplementedPlayerServer
	*SqliteHandler
	*AuthHandler
	*fortress.Logger
	*OnlinePlayers
}

// OnlinePlayers contains a map of recently active players keyed by their playerId
// The methods of OnlinePlayers are thread-safe
type OnlinePlayers struct {
	*sync.RWMutex
	players map[string]*fortress.Player
}

func (o *OnlinePlayers) GetOnlinePlayers() []*fortress.Player {
	onlinePlayers := make([]*fortress.Player, 0)

	o.RLock()
	for _, p := range o.players {
		onlinePlayers = append(onlinePlayers, p)
	}
	o.RUnlock()
	return onlinePlayers
}
func (o *OnlinePlayers) RemoveOnlinePlayer(playerId string) {
	o.Lock()
	delete(o.players, playerId)
	o.Unlock()
}
func (o *OnlinePlayers) AddOnlinePlayer(player *fortress.Player) error {
	if player == nil {
		return fmt.Errorf("player passed to addonlineplayer is nil")
	}
	o.Lock()
	o.players[player.GetPlayerId()] = player
	o.Unlock()
	return nil
}
func (o OnlinePlayers) IsOnline(filter PlayerFilter) bool {
	o.RLock()
	defer o.RUnlock()
	_, found := o.players[filter.playerId]
	if found {
		return true
	}
	return o.GetOnlinePlayer(filter) != nil
}
func (o *OnlinePlayers) PurgeInactives(inactiveTime time.Duration) int {
	inactives := []string{}
	o.RLock()
	for id, p := range o.players {
		if p.GetInactiveTime() > inactiveTime {
			inactives = append(inactives, id)
		}
	}
	o.RUnlock()
	for _, id := range inactives {
		o.RemoveOnlinePlayer(id)
	}
	return len(inactives)
}

// GetOnlinePlayer returns a reference to the Player object that matches any field in PlayerFilter,
// it returns nil if there is no match. This function does not check the database.
func (o *OnlinePlayers) GetOnlinePlayer(filter PlayerFilter) *fortress.Player {
	for _, v := range o.GetOnlinePlayers() {
		if (filter.playerId != "" && v.GetPlayerId() == filter.playerId) ||
			(filter.googleId != "" && v.GetGoogleId() == filter.googleId) ||
			(filter.name != "" && v.GetName() == filter.name) {
			return v
		}
	}
	return nil
}

func NewPlayerHandler(sqliteHandler *SqliteHandler, logger *fortress.Logger) *PlayerHandler {
	handler := &PlayerHandler{
		fgrpc.UnimplementedPlayerServer{},
		sqliteHandler,
		nil,
		logger,
		&OnlinePlayers{&sync.RWMutex{}, make(map[string]*fortress.Player)}}

	go func() {
		time.Sleep(20 * time.Second)
		for {
			handler.saveModifiedPlayersToDb()
			removed := handler.PurgeInactives(inactiveTime)
			if removed > 0 {
				handler.Logf("%d inactive players have been logged out", removed)
			}
			time.Sleep(dbSaveInterval)
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
	p := h.GetOnlinePlayer(filter)
	if p != nil {
		return p, false
	}

	// check the db
	if checkDatabase {
		p := h.SqliteHandler.LookupPlayerFromDb(filter)
		if p != nil {
			h.AddOnlinePlayer(p)
			h.Logf("Loaded player %s(%s) from database.", p.GetName(), p.GetPlayerId())
			return p, false
		}
	}

	//not found, create new
	p = fortress.NewPlayer()

	if _, err := uuid.Parse(filter.playerId); err == nil { // if the requested id was a uuid, set the new player to that id, otherwise, keep the generated one
		p.SetPlayerId(filter.playerId)
	}
	h.registerNewPlayer(p)
	h.Logf("Created a new player with ID %s", p.GetPlayerId())
	return p, true

}

// registerNewPlayer adds the given player to onlinePlayers and saves them to db
func (h *PlayerHandler) registerNewPlayer(player *fortress.Player) {
	h.AddOnlinePlayer(player)
	h.SqliteHandler.CreateNewPlayerDbRecord(player)
}

// saveUpdatedPlayersToDb saves any player that has been modified in the last minute to the database
func (h *PlayerHandler) saveModifiedPlayersToDb() {
	for _, p := range h.GetOnlinePlayers() {
		if time.Since(p.GetUpdatedAt()) < dbSaveInterval+(10*time.Second) {
			h.SqliteHandler.UpdatePlayerToDb(p)
		}
	}
}

// GetPlayerNameFromId returns the name associated with the given playerId, optionally checking the database. Returns "" if no player is found
func (h *PlayerHandler) GetPlayerNameFromId(playerId string, checkDatabase bool) string {
	p, _ := h.GetPlayer(PlayerFilter{playerId: playerId}, checkDatabase)
	return p.GetName()
}
