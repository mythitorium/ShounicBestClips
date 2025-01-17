package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/getsentry/sentry-go"
)

var database *Database

var envDBFile = getEnvOrDefault("CLIPS_DB", "votes.db?_mutex=full&_journal_mode=wal")
var envBindAddr = getEnvOrDefault("CLIPS_BIND", ":8081")
var envBehindProxy = os.Getenv("CLIPS_BEHIND_PROXY")

// Provided by build flags
var commitSHA string

// this makes the front end misreport?
//var votingDeadlineUnix int64 = 1737953940

var votingDeadlineUnix int64 = time.Date(2025, time.January, 27, 4, 59, 0, 0, time.UTC).Unix()
var voteCooldown time.Duration = 5

var totalUnculledClipsInDb int64

func main() {
	var err error

	if commitSHA != "" {
		fmt.Printf("Starting buildSHA: %s\n", commitSHA[:7])
	}

	err = sentry.Init(sentry.ClientOptions{
		Release:       commitSHA,
		SampleRate:    0.1,
		EnableTracing: true,
	})
	if err != nil {
		fmt.Printf("sentry.Init: %s\n", err)
	}
	defer sentry.Flush(2 * time.Second)

	fmt.Printf("Loading database %s\n", envDBFile)
	database, err = LoadDatabase(envDBFile)
	if err != nil {
		panic(err)
	}
	defer database.Close()

	serveMux := CustomMux{http.NewServeMux()}
	initRoutes(serveMux)

	go taskCullVideos()

	fmt.Printf("Starting http server on %s\n", envBindAddr)
	if err = http.ListenAndServe(envBindAddr, serveMux); err != nil {
		panic(err)
	}
}

func getEnvOrDefault(key string, defValue string) (value string) {
	value, exists := os.LookupEnv(key)
	if !exists {
		value = defValue
	}
	return
}

func UpdateUnculledClipTotal() {
	totalUnculledClipsInDb = database.GetTotalClips()
	fmt.Println("The total number of unculled, servable clips is now", totalUnculledClipsInDb)
}
