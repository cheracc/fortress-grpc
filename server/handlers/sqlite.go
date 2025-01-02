package handlers

import (
	"database/sql"
	"time"

	"github.com/cheracc/fortress-grpc"
	_ "github.com/mattn/go-sqlite3"
)

type SqliteHandler struct {
	db *sql.DB
	*fortress.Logger
}

type PlayerFilter struct {
	playerId string
	googleId string
	name     string
}

func NewSqliteHandler(logger *fortress.Logger) *SqliteHandler {
	handler := &SqliteHandler{}
	db, err := sql.Open("sqlite3", "./data.db")
	if err != nil {
		logger.Fatal(err.Error())
	}

	handler.Logger = logger
	handler.db = db
	return handler
}

func (h *SqliteHandler) InitializeDatabase() {
	// "CREATE DATABASE fortress;"
	// "CREATE TABLE 'players' ( "
	// 	 	"'player_id'	TEXT"
	// 		"'google_id'	TEXT"
	//		"'name'			TEXT"
	//		"'session_token' VARCHAR(255)"
	//		"'avatar_url	VARCHAR(255)"
	//		"'created_at'	"
	//		"'updated_at'	"
	//		"'last_read'	"

	stmt, err := h.db.Prepare("CREATE TABLE IF NOT EXISTS players (" +
		"player_id	TEXT, " +
		"google_id	TEXT, " +
		"name TEXT, " +
		"session_token TEXT, " +
		"avatar_url	TEXT, " +
		"created_at INTEGER, " +
		"updated_at INTEGER, " +
		"last_read INTEGER)")
	if err != nil {
		h.Fatal(err.Error())
	}
	stmt.Exec()
	defer stmt.Close()

	h.Log("initialized database and table")
}

func (h *SqliteHandler) LookupPlayerFromDb(f PlayerFilter) *fortress.Player {
	sql := "SELECT player_id, google_id, name, session_token, avatar_url, created_at, updated_at, last_read FROM players WHERE "
	sqlWherePlayerId := "player_id like '%" + f.playerId + "%' "
	sqlWhereGoogleId := "google_id like '%" + f.googleId + "%' "
	sqlWhereName := "name like '%" + f.name + "%'"

	if f.playerId != "" {
		sql = sql + sqlWherePlayerId
	}
	if f.googleId != "" {
		if f.playerId != "" {
			sql = sql + " OR "
		}
		sql = sql + sqlWhereGoogleId
	}
	if f.name != "" {
		if f.playerId != "" || f.googleId != "" {
			sql = sql + " OR "
		}
		sql = sql + sqlWhereName
	}

	rows, err := h.db.Query(sql)

	if err != nil {
		h.Fatal(err.Error())
	}
	if rows.Err() != nil {
		h.Fatal(rows.Err().Error())
	}
	defer rows.Close()

	players := make([]fortress.Player, 0)

	for rows.Next() {
		var playerId, googleId, name, sessionToken, avatarUrl string
		var created, updated, read int64
		err = rows.Scan(&playerId, &googleId, &name, &sessionToken, &avatarUrl, &created, &updated, &read)
		if err != nil {
			h.Fatal(err.Error())
		}

		player := fortress.LoadPlayer(playerId, googleId, name, sessionToken, avatarUrl, time.Unix(created, 0), time.Unix(updated, 0), time.Unix(read, 0))

		players = append(players, player)
	}

	if rows.Err() != nil {
		h.Fatal(rows.Err().Error())
	}

	if len(players) < 1 {
		h.Logf("SQL: no players found for player filter playerid:%s googleid:%s name:%s", f.playerId, f.googleId, f.name)
		return nil
	}
	if len(players) > 1 {
		h.Errorf("SQL: multiple players found for filter playerid:%s googleid:%s name:%s (these should all be unique??) only the first was returned", f.playerId, f.googleId, f.name)
		return &players[0]
	}

	return &players[0]
}

func (h *SqliteHandler) UpdatePlayerToDb(p *fortress.Player) {
	// "UPDATE players set google_id = ?, name = ?, session_token = ?, avatar_url = ?, created_at = ?, updated_at = ?, last_read = ? WHERE player_id = ?"
	stmt, err := h.db.Prepare("UPDATE players set google_id = ?, name = ?, session_token = ?, avatar_url = ?, created_at = ?, updated_at = ?, last_read = ? WHERE player_id = ?")
	if err != nil {
		h.Fatal(err.Error())
	}

	var res sql.Result
	res, err = stmt.Exec(p.GetGoogleId(), p.GetName(), p.GetSessionToken(), p.GetAvatarUrl(), p.CreatedAt.Unix(), p.UpdatedAt.Unix(), p.LastRead.Unix(), p.GetPlayerId())
	if err != nil {
		h.Fatal(err.Error())
	}
	defer stmt.Close()
	rows, _ := res.RowsAffected()

	if rows == 0 {
		h.Errorf("SQL: update for player %s did not match any rows (this should not happen) - creating a new record instead", p.GetPlayerId())
		h.CreateNewPlayerDbRecord(p)
		return
	}
	h.Logf("Updated database record for player %s(%s)", p.GetName(), p.GetPlayerId())
}

func (h *SqliteHandler) CreateNewPlayerDbRecord(p *fortress.Player) {
	// "INSERT INTO players (player_id, google_id, name, session_token, avatar_url, created_at, updated_at, last_read) VALUES (?, ?, ?, ?, ?, ?, ?, ?)"
	stmt, _ := h.db.Prepare("INSERT INTO players (player_id, google_id, name, session_token, avatar_url, created_at, updated_at, last_read) VALUES (?, ?, ?, ?, ?, ?, ?, ?)")
	stmt.Exec(p.GetPlayerId(), p.GetGoogleId(), p.GetName(), p.GetSessionToken(), p.GetAvatarUrl(), p.CreatedAt.Unix(), p.UpdatedAt.Unix(), p.LastRead.Unix())
	defer stmt.Close()

	h.Logf("Added new database record for player %s(%s)", p.GetName(), p.GetPlayerId())
}

// TODO needs to be implemented - just check to see if it exists in the name column case-insensitive
func (h *SqliteHandler) IsNameUnique(name string) bool {
	return true
}

func (h *SqliteHandler) CloseDb() {
	if err := h.db.Close(); err != nil {
		h.Error(err.Error())
	}
}
