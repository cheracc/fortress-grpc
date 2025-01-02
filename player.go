package fortress

import (
	"sync"
	"time"

	"github.com/google/uuid"
)

// a player object. used by both client and server
type Player struct {
	*sync.RWMutex
	playerId     string
	googleId     string
	name         string
	sessionToken string
	avatarURL    string

	CreatedAt time.Time
	UpdatedAt time.Time
	LastRead  time.Time
}

func NewPlayer() *Player {
	now := time.Now().UTC()
	player := Player{
		RWMutex:   &sync.RWMutex{},
		playerId:  generatePlayerID(),
		CreatedAt: now,
		UpdatedAt: now,
		LastRead:  now,
	}

	return &player
}

func LoadPlayer(pId string, gId string, name string, sToken string, aUrl string, created time.Time, updated time.Time, read time.Time) Player {
	p := Player{&sync.RWMutex{}, pId, gId, name, sToken, aUrl, created, updated, read}
	p.setAccessed()
	return p
}

func generatePlayerID() string {
	return uuid.NewString()
}

func (p *Player) setUpdated() {
	p.Lock()
	p.UpdatedAt = time.Now().UTC()
	p.Unlock()
}

func (p *Player) setAccessed() {
	p.Lock()
	p.LastRead = time.Now().UTC()
	p.Unlock()
}

func (p *Player) GetPlayerId() string {
	p.RLock()
	playerId := p.playerId
	p.RUnlock()
	p.setAccessed()
	return playerId
}

func (p *Player) SetPlayerId(id string) {
	p.Lock()
	p.playerId = id
	p.Unlock()
	p.setUpdated()
}

func (p *Player) GetGoogleId() string {
	p.RLock()
	googleId := p.googleId
	p.RUnlock()
	p.setAccessed()
	return googleId
}

func (p *Player) SetGoogleId(id string) {
	p.Lock()
	p.googleId = id
	p.Unlock()
	p.setUpdated()
}

func (p *Player) GetName() string {
	p.RLock()
	name := p.name
	p.RUnlock()
	p.setAccessed()
	return name
}

func (p *Player) SetName(name string) {
	p.Lock()
	p.name = name
	p.Unlock()
	p.setUpdated()
}

func (p *Player) GetSessionToken() string {
	p.RLock()
	sessionToken := p.sessionToken
	p.RUnlock()
	p.setAccessed()
	return sessionToken
}

func (p *Player) SetSessionToken(tokenString string) {
	p.Lock()
	p.sessionToken = tokenString
	p.Unlock()
	p.setUpdated()
}

func (p *Player) GetAvatarUrl() string {
	p.RLock()
	avatarURL := p.avatarURL
	p.RUnlock()
	p.setAccessed()
	return avatarURL
}

func (p *Player) SetAvatarUrl(url string) {
	p.setUpdated()
	p.Lock()
	p.avatarURL = url
	p.Unlock()
}

func (p *Player) GetUpdatedAt() time.Time {
	p.RLock()
	updatedAt := p.UpdatedAt
	p.RUnlock()
	p.setAccessed()
	return updatedAt
}

func (p *Player) GetCreatedAt() time.Time {
	p.RLock()
	createdAt := p.CreatedAt
	p.RUnlock()
	p.setAccessed()
	return createdAt
}

func (p *Player) SetCreatedAt(time time.Time) {
	p.setUpdated()
	p.Lock()
	p.CreatedAt = time
	p.Unlock()
}
