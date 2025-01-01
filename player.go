package fortress

import (
	"sync"
	"time"

	"github.com/google/uuid"
)

// a player object. used by both client and server
type Player struct {
	*sync.Mutex
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
		playerId:  generatePlayerID(),
		CreatedAt: now,
		UpdatedAt: now,
		LastRead:  now,
		Mutex:     &sync.Mutex{},
	}

	return &player
}

func LoadPlayer(pId string, gId string, name string, sToken string, aUrl string, created time.Time, updated time.Time, read time.Time) Player {
	p := Player{&sync.Mutex{}, pId, gId, name, sToken, aUrl, created, updated, read}
	p.setAccessed()
	return p
}

func generatePlayerID() string {
	return uuid.NewString()
}

func (p *Player) setUpdated() {
	p.UpdatedAt = time.Now().UTC()
}

func (p *Player) setAccessed() {
	p.LastRead = time.Now().UTC()
}

func (p *Player) GetPlayerId() string {
	p.Lock()
	defer p.Unlock()
	defer p.setAccessed()
	return p.playerId
}

func (p *Player) SetPlayerId(id string) {
	p.Lock()
	defer p.Unlock()
	p.playerId = id
	p.setUpdated()
}

func (p *Player) GetGoogleId() string {
	p.Lock()
	defer p.Unlock()
	defer p.setAccessed()
	return p.googleId
}

func (p *Player) SetGoogleId(id string) {
	p.Lock()
	defer p.Unlock()
	p.googleId = id
	p.setUpdated()
}

func (p *Player) GetName() string {
	p.Lock()
	defer p.Unlock()
	defer p.setAccessed()
	return p.name
}

func (p *Player) SetName(name string) {
	p.Lock()
	defer p.Unlock()
	p.name = name
	p.setUpdated()
}

func (p *Player) GetSessionToken() string {
	p.Lock()
	defer p.Unlock()
	defer p.setAccessed()
	return p.sessionToken
}

func (p *Player) SetSessionToken(tokenString string) {
	p.Lock()
	defer p.Unlock()
	p.sessionToken = tokenString
	p.setUpdated()
}

func (p *Player) GetAvatarUrl() string {
	p.Lock()
	defer p.Unlock()
	defer p.setAccessed()
	return p.avatarURL
}

func (p *Player) SetAvatarUrl(url string) {
	p.Lock()
	defer p.Unlock()
	p.setUpdated()
	p.avatarURL = url
}

func (p *Player) GetCreatedAt() time.Time {
	p.Lock()
	defer p.Unlock()
	defer p.setAccessed()
	return p.CreatedAt
}

func (p *Player) SetCreatedAt(time time.Time) {
	p.Lock()
	defer p.Unlock()
	p.setUpdated()
	p.CreatedAt = time
}

func (p *Player) TimeSinceLastRead() time.Duration {
	p.Lock()
	defer p.Unlock()
	return time.Since(p.LastRead)
}
