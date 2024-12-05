package main

import (
	"fmt"
	"net/http"
	"time"
)

type CustomMux struct{ *http.ServeMux }
type RouteFunc func(http.ResponseWriter, *http.Request)
type UserRouteFunc func(http.ResponseWriter, *http.Request, User)

// Basic HTTP route with logging
func (mux *CustomMux) NewRoute(pattern string, handler RouteFunc) {
	mux.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
		cw := &CustomResponseWriter{w, 200}

		start := time.Now()
		handler(cw, r)
		end := time.Since(start) / time.Millisecond

		fmt.Printf(
			"[%s] %dms %d %s %s\n",
			r.RemoteAddr,
			end,
			cw.statusCode,
			r.Method,
			r.RequestURI,
		)
	})
}

// HTTP route with User and logging.
func (mux *CustomMux) NewUserRoute(pattern string, handler UserRouteFunc) {
	mux.NewRoute(pattern, func(w http.ResponseWriter, r *http.Request) {
		user, err := mux.loadUser(w, r)
		if err != nil {
			return
		}

		handler(w, r, user)
	})
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

// Custom Writer so we can pull the statusCode for logging
type CustomResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *CustomResponseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}
