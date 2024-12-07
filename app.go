package main

import (
	"fmt"
	"net/http"
	"time"
)

var database *Database

// TODO vArgs
var argDBFile = "votes.db"
var argBindAddr = ":8081"

func main() {
	var err error

	fmt.Printf("Loading database %s\n", argDBFile)
	database, err = LoadDatabase(argDBFile)
	if err != nil {
		panic(err)
	}
	defer database.Close()

	// Start cleanup loop
	go cleanupActiveVotesLoop()

	serveMux := CustomMux{http.NewServeMux()}
	initRoutes(serveMux)

	fmt.Printf("Starting http server on %s\n", argBindAddr)
	if err = http.ListenAndServe(argBindAddr, serveMux); err != nil {
		panic(err)
	}
}

func cleanupActiveVotesLoop() {
	for {
		database.cleanExpiredVotes()
		time.Sleep(10 * time.Minute)
	}
}
