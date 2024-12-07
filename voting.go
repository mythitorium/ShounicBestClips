package main

import (
	"time"
)

var VoteExpireTime = 1 * time.Hour

// Represents a vote sent to a user.
// Contains a random UUID to prevent vote manipulation by modifying responses.
type VoteOptions struct {
	startTime time.Time
	Id        int64  `json:"id"`
	A         string `json:"a"`
	B         string `json:"b"`
}

func (vote *VoteOptions) IsExpired() bool {
	expiresAt := vote.startTime.Add(VoteExpireTime)
	return time.Now().After(expiresAt)
}
