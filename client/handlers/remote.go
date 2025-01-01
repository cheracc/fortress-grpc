package handlers

import (
	"context"
	"os/exec"
	"time"

	"github.com/cheracc/fortress-grpc"
	fgrpc "github.com/cheracc/fortress-grpc/grpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const server_addr string = "localhost:50051"

// a Remote handles communication with the remote server
type Remote struct {
	// the gRPC client for authorizations
	fgrpc.AuthClient
	// the gRPC client for handling commands
	fgrpc.CommandClient
	// the gRPC client for handling player data
	fgrpc.PlayerClient
	// the gRPC chat client
	fgrpc.ChatClient
	// the logger
	*fortress.Logger
	// the player that's currently loaded
	*fortress.Player
	browserWindow *exec.Cmd
}

// NewRemote constructs a new Remote with the given Logger
func NewRemote(logger *fortress.Logger) *Remote {
	conn, err := grpc.NewClient(server_addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Fatalf("failed to connect to fortressd at %v", server_addr)
	}
	//defer conn.Close()
	auth := fgrpc.NewAuthClient(conn)
	cmd := fgrpc.NewCommandClient(conn)
	pc := fgrpc.NewPlayerClient(conn)
	chat := fgrpc.NewChatClient(conn)

	return &Remote{auth, cmd, pc, chat, logger, fortress.NewPlayer(), nil}
}

// HasSessionToken returns whether there is an actual session token saved (this does not verify it)
func (r *Remote) HasSessionToken() bool {
	return len(r.GetSessionToken()) > 40
}

func (r *Remote) GetChatChannel() *ChatStream {
	ctx, cancelFunc := context.WithCancel(context.Background())

	stream, err := r.ChatClient.JoinChannel(ctx, &fgrpc.ChatRequest{SessionToken: r.GetSessionToken()})
	if err != nil {
		r.Error(err.Error())
		defer cancelFunc()
		return nil
	}
	return &ChatStream{stream, cancelFunc, nil}
}

func (r *Remote) GetPlayerData() {
	payload := &fgrpc.PlayerInfo{Id: r.GetPlayerId(), SessionToken: r.GetSessionToken()}
	response, err := r.PlayerClient.GetPlayerData(context.Background(), payload)
	if err != nil {
		r.Error(err.Error())
		return
	}
	r.SetName(response.GetName())
	r.SetCreatedAt(time.Unix(response.GetCreatedAt(), 0))
	r.Logf("Received player data for player %s(%s) Created at %s", r.GetName(), r.GetPlayerId(), r.GetCreatedAt())
}

// SendCommand sends user commands to the remote server and returns the response
func (r *Remote) SendCommand(cmd string, args string) string {
	playerInfo := &fgrpc.PlayerInfo{SessionToken: r.GetSessionToken(), Id: r.GetPlayerId()}
	r.Logf("Sending command to server: %s %s", cmd, args)
	payload, err := r.CommandClient.Command(context.Background(), &fgrpc.CommandInfo{PlayerInfo: playerInfo, CommandName: cmd, CommandArguments: args})
	if err != nil {
		r.Errorf("error calling SendCommand(): %v", err)
	}
	return payload.GetJsonPayload()
}

// Authorize() handles all client authorization
// On the first call to Authorize(), the server should respond with an oauth token (state) and a URL
//
//	The user needs to click on the URL to sign in to Google, once they do, the server completes the authorization of the account and waits for another Authorize() request
//	Authorize() will continue to be called every few seconds until the user completes the Google sign-in
//	Once logged in through Google, Authorize() will return a valid session token to be used henceforth
//	Authorize() is then called periodically to refresh the session token
func (r *Remote) Authorize() {
	authInfo, err := r.AuthClient.Authorize(context.Background(), &fgrpc.PlayerInfo{Id: r.GetPlayerId(), SessionToken: r.GetSessionToken()})
	if err != nil {
		r.Errorf("error calling Authorize(): %v", err)
		return
	}

	if authInfo.LoginURL != "" { // server sent a login url, so we will need to log in first
		r.SetPlayerId(authInfo.PlayerID)
		r.SetSessionToken(authInfo.SessionToken) // this should be an oauth token that we will return to confirm we are the same as who logged in with that link
		r.ToConsolef("Use the following link to log in: %s\n\r", authInfo.LoginURL)
		r.browserWindow = exec.Command("rundll32", "url.dll,FileProtocolHandler", authInfo.LoginURL)
		r.browserWindow.Start()
		return
	}

	if authInfo.SessionToken == "" {
		r.Warn("No session or oauth token received from server...")
		return
	}

	if len(authInfo.SessionToken) > 40 {
		if len(r.GetSessionToken()) < 40 { // received a session token but had an oauth token. this means we've successfully logged in
			r.SetSessionToken(authInfo.SessionToken)
			r.GetPlayerData()
			r.Logf("Logged in as %s(%s)", r.GetName(), r.GetPlayerId())
			r.browserWindow.Process.Kill()
		}
	}

	r.SetSessionToken(authInfo.SessionToken)
	if authInfo.PlayerID != "" {
		r.SetPlayerId(authInfo.PlayerID)
	}
}

// Logout clears player data and calls Authorize()
func (r *Remote) Logout() {
	r.SetPlayerId("")
	r.SetName("")
	r.SetSessionToken("")
	r.Authorize()
}
