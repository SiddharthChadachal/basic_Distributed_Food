package internal

type Role int

const (
	Follower Role = iota
	Candidate
	Leader
)
