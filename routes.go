package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func initRoutes(serveMux CustomMux) {
	serveMux.NewRoute("/", routeRoot)
	serveMux.NewUserRoute("/nextVote", routeNextVote)
	serveMux.NewUserRoute("/submitVote", routeSubmitVote)
}

// Middleware TODO
//		Rate limiting

// Base route, return HTML template
func routeRoot(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte("TODO return main page"))
}

func routeNextVote(w http.ResponseWriter, req *http.Request, user User) {
	options, err := database.GetNextVoteForUser(user)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("Failed to fetch from database."))
		// TODO log to Sentry
		fmt.Printf("Failed to get new votes for user %v \"%s\"\n", user, err)
		return
	}

	// User has completed their queue
	if options == nil {
		w.WriteHeader(204) // NO_CONTENT
		w.Write([]byte("No more items to vote on!"))
		return
	}

	bytes, err := json.Marshal(options)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("Failed to write JSON data?"))
		// TODO log to Sentry
		fmt.Printf("Failed to write json data %v\n", options)
		return
	}

	w.Write(bytes)
}

func routeSubmitVote(w http.ResponseWriter, req *http.Request, user User) {
	// TODO Get CurrentVote from UUID
	// TODO check CurrentVote time.
	//		Users should spend at min 1 minute per vote
	//		Depending on video lengths.
	// 		If a user votes too fast, return 200 and
	//		pretend we accepted it.

	// TODO return 2xx when user has completed their queue

	w.Write([]byte("TODO submit vote"))
}

// TODO /myVotes

// TODO /totalVotes
