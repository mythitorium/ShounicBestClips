package main

import (
	"fmt"
	"time"

	"github.com/getsentry/sentry-go"
)

const minCullVotes = 20
const minCullRatio = .40
const intervalCullTask = 1 * time.Hour

type VoteState = map[string]*VideoStats

type VideoStats struct {
	totalVotes int
	totalScore int
}

func (cv *VideoStats) ShouldCull() bool {
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

		UpdateUnculledClipTotal()

		time.Sleep(intervalCullTask)
	}
}

func cullVideos(database *Database) error {
	tx, err := database.Begin()
	if err != nil {
		return err
	}

	videos := make(VoteState)

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

		videos[url] = &VideoStats{}
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

		if videos[url] == nil {
			videos[url] = &VideoStats{}
		}

		videos[url].totalScore += score
		videos[url].totalVotes += 1
	}

	err = rows.Close()
	if err != nil {
		return err
	}

	// Reset table to allow videos to return to the queue
	// if we adjust the numbers in prod to be more lenient.
	_, err = tx.Exec("DELETE FROM culled_videos")
	if err != nil {
		tx.Rollback()
		return err
	}

	fmt.Println("Culling Debug:\n\tUrl Votes Score")
	for url, stats := range videos {
		fmt.Printf("\t%s %v ", url, stats)
		fmt.Printf(
			"(%d/%d) = %f > %f\n",
			stats.totalScore,
			stats.totalVotes,
			float32(stats.totalScore)/float32(stats.totalVotes),
			minCullRatio,
		)
		if stats.ShouldCull() {
			_, err = tx.Exec(
				"INSERT OR IGNORE INTO culled_videos VALUES (?)",
				url,
			)
			if err != nil {
				tx.Rollback()
				return err
			}
		}
	}

	return tx.Commit()
}
