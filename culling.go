package main

import (
	"fmt"
	"time"

	"github.com/getsentry/sentry-go"
)

const minCullVotes = 20
const minCullRatio = .30
const intervalCullTask = 1 * time.Hour

type voteState = map[string]*CullVotes

type CullVotes struct {
	totalVotes int
	totalScore int
}

func (cv *CullVotes) ShouldCull() bool {
	return cv.totalVotes > minCullVotes &&
		float32(cv.totalScore)/float32(cv.totalVotes) < minCullRatio
}

func taskCullVideos() {
	var err error
	for {
		fmt.Println("Running cull task")

		if err = cullVideos(database); err != nil {
			fmt.Print(err.Error())
			sentry.CaptureException(err)
		}

		time.Sleep(intervalCullTask)
	}
}

func cullVideos(database *Database) error {
	tx, err := database.Begin()
	if err != nil {
		return err
	}

	videos := make(voteState)

	// Populate map
	rows, err := tx.Query("SELECT url FROM videos")
	if err != nil {
		return err
	}

	for rows.Next() {
		var url string
		err = rows.Scan(&url)
		if err != nil {
			return err
		}

		videos[url] = &CullVotes{}
	}

	err = rows.Close()
	if err != nil {
		return err
	}

	// Count the vote scores and vote counts
	rows, err = tx.Query("SELECT video_url, score FROM votes")
	if err != nil {
		return err
	}

	for rows.Next() {
		var url string
		var score int
		err = rows.Scan(&url, &score)
		if err != nil {
			return err
		}

		videos[url].totalScore += score
		videos[url].totalVotes += 1
	}

	err = rows.Close()
	if err != nil {
		return err
	}

	fmt.Println("Culling Debug:\n\tUrl Votes Score")
	for url, votes := range videos {
		fmt.Printf("\t%s %v\n", url, votes)
		if votes.ShouldCull() {
			_, err = tx.Exec(
				"INSERT OR IGNORE INTO culled_videos VALUES (?)",
				url,
			)
			if err != nil {
				return err
			}
		}
	}

	return tx.Commit()
}
