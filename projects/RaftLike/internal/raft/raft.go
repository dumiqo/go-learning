package raft

type Raft struct {
	Type Type
}

type Type int

const (
	Follower Type = iota
	Candidate
	Leader
)
