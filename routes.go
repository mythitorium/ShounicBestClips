package main

import (
	"encoding/json"
	"fmt"
	"io"
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

	// Send new vote to client
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
	if err := req.ParseForm(); err != nil {
		w.WriteHeader(406)
		w.Write([]byte("Failed to parse form input."))
		return
	}

	choice := req.PostForm.Get("choice")
	if choice == "" {
		w.WriteHeader(400)
		w.Write([]byte("No choice given"))
		return
	}

	err := database.SubmitUserVote(user, choice)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("Failed to communicate with database."))
		// TODO log to Sentry
		fmt.Printf("Failed to submit vote from %v of \"%s\": %v\n", user, choice, err)
		return
	}

	routeNextVote(w, req, user)
}

// TODO /myVotes

// TODO /totalVotes

func readJsonRequest(output *any, req *http.Request) (err error) {
	data, err := io.ReadAll(req.Body)
	if err != nil {
		return err
	}
	defer req.Body.Close()
	return json.Unmarshal(data, output)
}
