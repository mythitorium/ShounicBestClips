package main

import (
	"fmt"
	"net/http"
)

type CustomMux struct{ *http.ServeMux }
type RouteFunc func(http.ResponseWriter, *http.Request)
type UserRouteFunc func(http.ResponseWriter, *http.Request, User)

// Basic HTTP route with logging
func (mux *CustomMux) newRoute(pattern string, handler RouteFunc) {
	mux.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
		mux.logRequest(r)
		handler(w, r)
	})
}

// HTTP route with User and logging.
func (mux *CustomMux) newUserRoute(pattern string, handler UserRouteFunc) {
	mux.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
		mux.logRequest(r)

		user, err := mux.loadUser(w, r)
		if err != nil {
			return
		}

		handler(w, r, user)
	})
}

func (mux *CustomMux) logRequest(r *http.Request) {
	fmt.Printf("[%s] %s %s\n", r.RemoteAddr, r.Method, r.RequestURI)
}

// Load the user from database.
// If there is an error, a 500 error will automatically be written to the ResponseWriter.
func (mux *CustomMux) loadUser(w http.ResponseWriter, r *http.Request) (user User, err error) {
	user, err = database.GetUser(r.RemoteAddr)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("Failed to get user!"))
		// TODO log to Sentry
		fmt.Printf("Failed to get User \"%s\"\n", err)
	}

	return
}
