package internal

import (
	"sync"
	"time"
)

type Peer struct {
	ID   int
	Host string
	Port int
}

type Node struct {
	ID                 int
	Port               int
	role               Role
	Peers              []Peer
	LeaderId           int
	LatestHeartBeat    time.Time
	Mutex              sync.Mutex
	ElectionInProgress bool
}

func NewNode(config Config) *Node {
	return &Node{
		ID:                 config.ID,
		Port:               config.Port,
		Peers:              config.Peers,
		LeaderId:           -1,
		role:               Follower,
		LatestHeartBeat:    time.Now(),
		ElectionInProgress: false,
	}
}

func (this *Node) GetLeader() int {
	return this.LeaderId
}

func (this *Node) BecomeLeader() {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()
	this.LeaderId = this.ID
	this.role = Leader
}

func (this *Node) BecomeFollower(leaderId int) {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()
	this.LeaderId = leaderId
	this.role = Follower
}

// if ping to all the peers fail
// then change role to candidate

func (this *Node) BecomeCandidate() {
	this.role = Candidate
}

func (this *Node) UpdateHeartBeat() {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()
	this.LatestHeartBeat = time.Now()
}

func (this *Node) IsLeader() bool {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()
	return this.ID == this.LeaderId && this.role == Leader
}

func (this *Node) IsCandidate() bool {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()
	return this.ID != this.LeaderId && this.role == Candidate
}

func (this *Node) IsFollower() bool {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()
	return this.ID != this.LeaderId && this.role == Follower
}

func (this *Node) Serve() {
	// cooking functionality for follower
}

func (this *Node) FetchFood() {
	// receiving serving functionality for leader to serve requests
}
