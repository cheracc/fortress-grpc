package handlers

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/cheracc/fortress-grpc"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

const oauthGoogleUrlAPI = "https://www.googleapis.com/oauth2/v2/userinfo?access_token="

var googleOauthConfig = &oauth2.Config{
	RedirectURL:  "http://localhost:8000/auth/google/callback",
	ClientID:     os.Getenv("GOOGLE_OAUTH_CLIENT_ID"),
	ClientSecret: os.Getenv("GOOGLE_OAUTH_CLIENT_SECRET"),
	Scopes:       []string{"openid"},
	Endpoint:     google.Endpoint,
}

type OauthHandler struct {
	httpServer  *http.Server
	mux         *http.ServeMux
	oauthStates map[string]string
	*AuthHandler
	*fortress.Logger
}

type GoogleJson struct {
	Id      string `json:"id"`
	Picture string `json:"picture"`
}

func (h *OauthHandler) StartListener(authHandler *AuthHandler) {
	h.AuthHandler = authHandler
	h.mux.HandleFunc("/auth/google/callback", h.OauthGoogleCallback)

	go func() {
		h.Logf("Starting Oauth http server, listening on %s", h.httpServer.Addr)
		if err := h.httpServer.ListenAndServe(); err != nil {
			h.Fatalf("%v", err)
		}
	}()
}

func NewOauthHandler(logger *fortress.Logger) OauthHandler {
	oauthHandler := OauthHandler{}
	oauthHandler.Logger = logger
	mux := http.NewServeMux()
	oauthHandler.mux = mux

	server := &http.Server{
		Addr:    "localhost:8000",
		Handler: oauthHandler.mux,
	}

	oauthHandler.httpServer = server

	oauthHandler.oauthStates = make(map[string]string)
	logger.Log("Finished configuring Oauth Handler")
	return oauthHandler
}

func (h *OauthHandler) OauthGoogleCallback(w http.ResponseWriter, r *http.Request) {
	h.Log("Received callback from Google:")
	receivedState := r.FormValue("state") // the random state we attached to the login url

	auth := h.authenticatingPlayers[receivedState]
	if auth == nil { // don't recognize this state
		h.Log("Invalid oath state received")
		return
	}

	h.Log("Received state is valid, requesting user data from Google")
	data := h.fetchGoogleDataFromResponse(r.FormValue("code"))
	h.Logf("Received data: %s", data)

	// set the google id to what we received
	auth.googleId = data.Id
	auth.avatarUrl = data.Picture

	// authorize this player to login
	player, _ := h.GetPlayerByGoogleID(data.Id)
	h.AuthorizePlayer(player)
	auth.sessionToken = player.GetSessionToken()

	fmt.Fprintf(w, "<h1>You may close this tab</h1>")
}

func (h *OauthHandler) fetchGoogleDataFromResponse(code string) *GoogleJson {
	data := h.getUserDataFromGoogle(code)
	authInfo := GoogleJson{}

	err := json.Unmarshal(data, &authInfo)
	if err != nil {
		h.Log(err.Error())
	}

	return &authInfo
}

func (h *OauthHandler) getUserDataFromGoogle(code string) []byte {
	// Use 'code' to get token and get user info from Google.

	token, err := googleOauthConfig.Exchange(context.Background(), code)
	if err != nil {
		h.Errorf("code exchange wrong: %s", err.Error())
		return nil
	}
	response, err := http.Get(oauthGoogleUrlAPI + token.AccessToken)
	if err != nil {
		h.Errorf("failed getting user info: %s", err.Error())
		return nil
	}
	defer response.Body.Close()
	contents, err := io.ReadAll(response.Body)
	if err != nil {
		h.Errorf("failed read response: %s", err.Error())
		return nil
	}
	return contents
}

//build link to sign in with google

func (h *OauthHandler) GenerateLoginURL() (string, string) {
	oauthState := generateRandomState()
	url := googleOauthConfig.AuthCodeURL(oauthState)

	return url, oauthState
}

func generateRandomState() string {
	b := make([]byte, 16)
	rand.Read(b)

	return base64.URLEncoding.EncodeToString(b)
}
