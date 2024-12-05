package main

import (
	"database/sql"
	"fmt"
	"net"

	_ "github.com/mattn/go-sqlite3"
)

func LoadDatabase(file string) (db *Database, err error) {
	conn, err := sql.Open("sqlite3", file)
	db = &Database{conn}

	if err == nil {
		err = db.setup()
	}

	return
}

type Database struct{ *sql.DB }

// Setup the database.
// Ran every time we load the database.
func (db *Database) setup() (err error) {
	var setupQueries = []string{
		// TODO ? video: title, uploader, docSubmitter, upload date
		"CREATE TABLE IF NOT EXISTS videos (id INTEGER PRIMARY KEY, url TEXT)",
		"CREATE TABLE IF NOT EXISTS votes (user_id INTEGER, video_id INTEGER, score INTEGER)",
		"CREATE TABLE IF NOT EXISTS users (id INTEGER PRIMARY KEY, ip TEXT UNIQUE)",
	}

	// Transaction so we can undo if we error
	tran, err := db.Begin()
	if err != nil {
		return
	}

	// Run all setupQueries
	for _, query := range setupQueries {
		_, err = db.Exec(query)
		if err != nil {
			return
		}
	}

	// Commit transaction
	return tran.Commit()
}

func (db *Database) GetUser(remoteAddr string) (user User, err error) {
	remoteAddr, _, err = net.SplitHostPort(remoteAddr)
	if err != nil {
		return
	}

	row, err := db.Query(
		"INSERT OR IGNORE INTO users(ip) VALUES (?); SELECT * FROM users WHERE ip=?;",
		remoteAddr,
		remoteAddr,
	)
	if err != nil {
		return
	}
	defer row.Close()

	if !row.Next() {
		err = fmt.Errorf("GetUser gave 0 results for %s", remoteAddr)
		return
	}

	err = row.Scan(&user.id, &user.ip)
	return
}

func (db *Database) GetNewVoteForUser(user User) (err error)            { return }
func (db *Database) SubmitUserVote(user User, video string) (err error) { return }
func (db *Database) IsUserQueueComplete(user User) (bool, err error)    { return }
