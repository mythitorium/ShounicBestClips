package main

import (
	"math"
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
	// TODO return 2xx when user has completed their queue

	w.Write([]byte("TODO return clips"))
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

// TODO /voteClips

// TODO /myVotes

// TODO /totalVotes

func requestToUser(req *http.Request) (user User) {
	return User{math.MaxUint, req.RemoteAddr}
}
