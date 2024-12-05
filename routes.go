package main

import (
	"fmt"
	"math"
	"net/http"
)

func initRoutes() {
	http.HandleFunc("/", routeRoot)
	http.HandleFunc("/nextVote", routeNextVote)
	http.HandleFunc("/submitVote", routeSubmitVote)
}

// Middleware TODO
//		Log routes
//		Rate limiting
//		User loading

// Base route, return HTML template
func routeRoot(w http.ResponseWriter, req *http.Request) {

	// TODO move to a middleware
	user, err := database.GetUser(req.RemoteAddr)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("Failed to get user!"))
		fmt.Printf("Failed to get User \"%s\"\n", err)
		return
	}

	fmt.Printf("Req from %v\n", user)

	w.Write([]byte("TODO return main page"))
}

func routeNextVote(w http.ResponseWriter, req *http.Request) {
	// TODO validate user

	// TODO return 2xx when user has completed their queue

	w.Write([]byte("TODO return clips"))
}

func routeSubmitVote(w http.ResponseWriter, req *http.Request) {
	// TODO validate user

	// TODO Get CurrentVote from UUID
	// TODO check CurrentVote time.
	//		Users should spend at min 1 minute per vote
	//		Depending on video lengths.
	// 		If a user votes too fast, return 200 and
	//		pretend we accepted it.

	// TODO return 2xx when user has completed their queue

	w.Write([]byte("TODO submit vote"))
}

// TODO /voteClips

// TODO /myVotes

// TODO /totalVotes

func requestToUser(req *http.Request) (user User) {
	return User{math.MaxUint, req.RemoteAddr}
}
