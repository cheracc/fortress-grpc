package fortress

import (
	"time"

	"github.com/google/uuid"
)

// a player object. used by both client and server
type Player struct {
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
	}

	return &player
}

func LoadPlayer(pId string, gId string, name string, sToken string, aUrl string, created time.Time, updated time.Time, read time.Time) Player {
	p := Player{pId, gId, name, sToken, aUrl, created, updated, read}
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
	defer p.setAccessed()
	return p.playerId
}

func (p *Player) SetPlayerId(id string) {
	p.playerId = id
	p.setUpdated()
}

func (p *Player) GetGoogleId() string {
	defer p.setAccessed()
	return p.googleId
}

func (p *Player) SetGoogleId(id string) {
	p.googleId = id
	p.setUpdated()
}

func (p *Player) GetName() string {
	defer p.setAccessed()
	return p.name
}

func (p *Player) SetName(name string) {
	p.name = name
	p.setUpdated()
}

func (p *Player) GetSessionToken() string {
	defer p.setAccessed()
	return p.sessionToken
}

func (p *Player) SetSessionToken(tokenString string) {
	p.sessionToken = tokenString
	p.setUpdated()
}

func (p *Player) GetAvatarUrl() string {
	defer p.setAccessed()
	return p.avatarURL
}

func (p *Player) SetAvatarUrl(url string) {
	p.setUpdated()
	p.avatarURL = url
}

func (p *Player) GetCreatedAt() time.Time {
	defer p.setAccessed()
	return p.CreatedAt
}

func (p *Player) SetCreatedAt(time time.Time) {
	p.setUpdated()
	p.CreatedAt = time
}
