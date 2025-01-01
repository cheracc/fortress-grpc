package handlers

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"time"

	"github.com/cheracc/fortress-grpc"
	fgrpc "github.com/cheracc/fortress-grpc/grpc"
	"github.com/golang-jwt/jwt/v5"
	"google.golang.org/grpc/peer"
)

// The AuthHandler handles all user authorizations. It holds a reference to the PlayerHandler
type AuthHandler struct {
	fgrpc.UnimplementedAuthServer                  // the gRPC server that accepts auth requests
	*PlayerHandler                                 // a pointer to the PlayerHandler to lookup players
	OauthHandler                                   // a separate handler to deal with oAuth stuff
	*fortress.Logger                               // the logger
	*ecdsa.PrivateKey                              // the private key used for signing auth tokens
	authenticatingPlayers         map[string]*Auth // the players that are currently in the process of authenticating
}

// JwtTokenClaims contains the data that we save on the session token
type JwtTokenClaims struct {
	PlayerID string `json:"player-id"`
	jwt.RegisteredClaims
}

// Auth contains the information used during the authorization of the player. The session token will be the final field to be populated
type Auth struct {
	googleId     string
	oauthState   string
	sessionToken string
	avatarUrl    string
	ipAddress    string
}

// isComplete returns whether the Auth is fully populated
func (a *Auth) isComplete() bool {
	return (a.googleId != "") && (a.oauthState != "") && (a.sessionToken != "")
}

// Authorize is the gRPC receiving function for authorization requests.
func (h *AuthHandler) Authorize(ctx context.Context, playerInfo *fgrpc.PlayerInfo) (*fgrpc.AuthInfo, error) {
	ip, _ := peer.FromContext(ctx)
	unverifiedPlayerId := playerInfo.GetId() // if this is the first time a user is connecting, this is just a randomly generated uuid from the client
	h.Logf("Recieved authorization request from %s", ip.Addr.String())

	authInfo := fgrpc.AuthInfo{PlayerID: unverifiedPlayerId} // set up the response payload

	receivedSessionToken := playerInfo.GetSessionToken() // the token the user sent - either a random uuid oauthstate token or a signed session token
	if receivedSessionToken == "" {                      // this is user's first attempt to auth
		var state string
		authInfo.LoginURL, state = h.OauthHandler.GenerateLoginURL()
		authInfo.SessionToken = state

		h.authenticatingPlayers[state] = &Auth{oauthState: state, ipAddress: ip.Addr.String()}
		return &authInfo, nil
	}

	if len(receivedSessionToken) < 40 { // this is an oauthstate token, this user is in the process of authenticating
		auth := h.authenticatingPlayers[receivedSessionToken]
		if auth == nil { // this player had an auth token, but we have no record of it. it may be very old. just send them a new one and a link
			var state string
			authInfo.LoginURL, state = h.OauthHandler.GenerateLoginURL()
			authInfo.SessionToken = state

			h.authenticatingPlayers[state] = &Auth{oauthState: state, ipAddress: ip.Addr.String()}
			return &authInfo, nil
		}
		if auth.isComplete() { // this authorization is complete (logged in with google). load their player and send them a session token
			player, _ := h.GetPlayer(PlayerFilter{googleId: auth.googleId}, true) // fetching by google id since its trusted
			player.SetAvatarUrl(auth.avatarUrl)
			player.SetSessionToken(auth.sessionToken)
			authInfo.PlayerID = player.GetPlayerId()
			authInfo.SessionToken = auth.sessionToken
			authInfo.LoginURL = ""

			h.authenticatingPlayers[auth.oauthState] = nil // remove their auth
			return &authInfo, nil
		}
		// from here they are still authenticating, just send their same info back
		authInfo.SessionToken = playerInfo.SessionToken
		authInfo.LoginURL = ""
		return &authInfo, nil
	}

	// if reaching here, the user has an actual session token, but we don't know if its valid
	hasValidToken := h.verifySessionToken(receivedSessionToken) // see if this token is valid (and not expired)
	if hasValidToken {                                          // already had valid token
		playerId := h.GetPlayerIdFromTokenString(receivedSessionToken)   // use the playerId from the session token, as it is signed
		player, _ := h.GetPlayer(PlayerFilter{playerId: playerId}, true) // log this player in if they are not already
		h.Logf("Player %s(%s) refreshed session token", player.GetName(), playerId)
		h.AuthorizePlayer(player)

		authInfo.PlayerID = player.GetPlayerId()
		authInfo.SessionToken = player.GetSessionToken()
		authInfo.LoginURL = ""

		return &authInfo, nil
	} else {
		h.Logf("Received token from %s[unverified] invalid or expired", unverifiedPlayerId) // this user had a token but its invalid or expired
		var state string
		authInfo.LoginURL, state = h.GenerateLoginURL()
		authInfo.SessionToken = state

		h.authenticatingPlayers[state] = &Auth{oauthState: state, ipAddress: ip.Addr.String()}
		return &authInfo, nil
	}
}

// NewAuthHandler constructs a new AuthHandler using the given gRCP Server, PlayerHandler, and Logger
func NewAuthHandler(playerHandler *PlayerHandler, logger *fortress.Logger) *AuthHandler {
	randomKey, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)

	handler := &AuthHandler{fgrpc.UnimplementedAuthServer{}, playerHandler, NewOauthHandler(logger), logger, randomKey, make(map[string]*Auth)}

	handler.OauthHandler.StartListener(handler)

	logger.Log("Started AuthHandler")
	return handler
}

// AuthorizePlayer sets a new fresh session token on the player and saves them to the database
func (h *AuthHandler) AuthorizePlayer(player *fortress.Player) error {
	player.SetSessionToken(h.generateToken(player.GetPlayerId()))
	h.PlayerHandler.SqliteHandler.UpdatePlayerToDb(player)
	h.updateOnlinePlayer(player)

	return nil
}

// generateToken generates a new JWT session token using the provided playerId
func (h *AuthHandler) generateToken(playerId string) string {
	expirationTime := time.Now().UTC().Add(5 * time.Minute)

	claims := &JwtTokenClaims{
		PlayerID: playerId,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)

	if h.PrivateKey == nil {
		h.Fatal("Tried to create new session token but there is no private key loaded.")
	}

	tokenString, err := token.SignedString(h.PrivateKey)
	if err != nil {
		h.Logf("could not create new jwt for pid %s: %s", playerId, err)
		return ""
	}
	return tokenString
}

// verifySessionToken verifies the signature of the token and its expiry, and returns whether both are valid
func (h *AuthHandler) verifySessionToken(token string) bool {
	tkn, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		return &h.PrivateKey.PublicKey, nil
	})

	if err != nil {
		h.Errorf("error parsing jwt token: %s: %s", token, err)
		return false
	}

	if !tkn.Valid {
		h.Errorf("session token %s is invalid (privatekey %s)", token, h.PrivateKey)
		return false
	}
	return true
}

// GetPlayerIdFromTokenString accepts a session token string and extracts the PlayerId from it.
// It returns "" if it was unable to parse the tokenString.
func (h *AuthHandler) GetPlayerIdFromTokenString(tokenString string) string {
	_, claims := h.getTokenFromString(tokenString)
	return claims.PlayerID
}

func (h *AuthHandler) getTokenFromString(tokenString string) (*jwt.Token, *JwtTokenClaims) {
	claims := &JwtTokenClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (any, error) {
		return &h.PrivateKey.PublicKey, nil
	})
	if err != nil {
		h.Errorf("could not parse jwt token %s: %w", tokenString, err)
		return nil, &JwtTokenClaims{}
	}
	return token, claims
}

// IsValidToken is an exported function to perform token validations, it simply calls verifySessionToken() with the given token string
func (h *AuthHandler) IsValidToken(tokenString string) bool {
	valid := h.verifySessionToken(tokenString)
	return valid
}
