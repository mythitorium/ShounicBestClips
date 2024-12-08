package main

import (
	"time"
)

// Represents a vote sent to a user.
// Contains a random UUID to prevent vote manipulation by modifying responses.
type VoteOptions struct {
	startTime time.Time
	A         string `json:"a"`
	B         string `json:"b"`
}
