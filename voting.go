package main

// Represents a vote sent to a user.
// Contains a random UUID to prevent vote manipulation by modifying responses.
type CurrentVote struct {
	uuid      string
	startTime uint
	timeoutAt uint
	a         string
	b         string
}
