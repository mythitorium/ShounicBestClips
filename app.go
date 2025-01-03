package main

import (
	"fmt"
	"net/http"
	"time"
)

var database *Database

// TODO vArgs
var argDBFile = "votes.db?_mutex=full"
var argBindAddr = ":8081"
var argMaxVoteTime = 4 * time.Hour

// TODO: Make this less stupid -myth
var votingDeadlineUnix = "1736496000"

func main() {
	var err error

	fmt.Printf("Loading database %s\n", argDBFile)
	database, err = LoadDatabase(argDBFile)
	if err != nil {
		panic(err)
	}
	defer database.Close()

	serveMux := CustomMux{http.NewServeMux()}
	initRoutes(serveMux)

	fmt.Printf("Starting http server on %s\n", argBindAddr)
	if err = http.ListenAndServe(argBindAddr, serveMux); err != nil {
		panic(err)
	}
}
