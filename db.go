package main

import (
	"database/sql"
	"fmt"
	"math/rand/v2"
	"net"
	"time"

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
		"CREATE TABLE IF NOT EXISTS videos (id INTEGER PRIMARY KEY, url TEXT UNIQUE)",
		"CREATE TABLE IF NOT EXISTS votes (user_id INTEGER, video_id INTEGER, score INTEGER)",
		"CREATE TABLE IF NOT EXISTS users (id INTEGER PRIMARY KEY, ip TEXT UNIQUE)",
		"CREATE TABLE IF NOT EXISTS active_votes ( " +
			"id INTEGER PRIMARY KEY NOT NULL, " +
			"user_id INTEGER UNIQUE NOT NULL, " +
			"start_time INTEGER NOT NULL, " +
			"a TEXT, " +
			"b TEXT " +
			")",

		// TEST
		"INSERT OR IGNORE INTO videos (url) VALUES ('one'),('two'),('three'),('four'),('five')",
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

// Get the next vote for a user
// If a vote already exists, it will be deleted.
// If there are < 2 options, `vote` will be nil
func (db *Database) GetNextVoteForUser(user User) (vote *VoteOptions, err error) {
	a, b, err := db.findNextPair(user)
	if a == "" || b == "" || err != nil {
		// Return nil vote, we don't have enough
		// voting options for this user
		return
	}

	vote = &VoteOptions{time.Now(), rand.Int64(), a, b}
	row, err := db.Query(
		"DELETE FROM active_votes WHERE user_id=?;"+
			"INSERT INTO active_votes VALUES (?, ?, ?, ?, ?) RETURNING id;",
		user.id,
		vote.Id,
		user.id,
		vote.startTime,
		a,
		b,
	)
	if err != nil {
		return
	}
	defer row.Close()

	if row.Next() {
		err = row.Scan(&vote.Id)
	}

	return
}

// Get new vote options for the user
// Empty a or b strings means not enough available voting options
func (db *Database) findNextPair(user User) (a string, b string, err error) {
	row, err := db.Query(
		"SELECT url FROM videos WHERE id NOT IN (SELECT video_id FROM votes WHERE user_id = ?) ORDER BY random() LIMIT 2",
		user.id,
	)
	if err != nil {
		return
	}
	defer row.Close()

	if !row.Next() {
		// 0 videos available
		return
	}

	err = row.Scan(&a)
	if err != nil || !row.Next() {
		return
	}

	err = row.Scan(&b)
	return
}

func (db *Database) cleanExpiredVotes() (count int64, err error) {
	res, err := db.Exec(
		// TODO fixme
		"DELETE FROM active_votes WHERE start_time + ? > unixepoch()",
		VoteExpireTime,
	)
	if err != nil {
		return
	}
	return res.RowsAffected()
}

func (db *Database) SubmitUserVote(user User, voteId VoteOptions, choice string) (err error) {
	// TODO check if vote is expired

	return
}
func (db *Database) IsUserQueueComplete(user User) (bool, err error) { return }
