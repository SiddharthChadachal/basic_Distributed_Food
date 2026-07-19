package internal

import (
	"fmt"
	"time"
)

type Election struct {
	Node *Node
	RPC  *RPC
}

func (election *Election) StartChecking() {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		<-ticker.C

		if election.Node.IsLeader() {
			continue
		} else {
			if election.HasTimedOut() {
				election.StartElection()
			}
		}
	}
}

func (election *Election) HasTimedOut() bool {
	election.Node.Mutex.Lock()
	defer election.Node.Mutex.Unlock()
	return time.Since(election.Node.LatestHeartBeat) > 5*time.Second
}

func (election *Election) StartElection() {

	election.Node.Mutex.Lock()
	if election.Node.ElectionInProgress {
		election.Node.Mutex.Unlock()
		return
	}

	election.Node.ElectionInProgress = true
	election.Node.Mutex.Unlock()

	defer func() {
		election.Node.Mutex.Lock()
		election.Node.ElectionInProgress = false
		election.Node.Mutex.Unlock()
	}()

	fmt.Println("Starting Election...")

	election.Node.BecomeCandidate()

	highestId := election.Node.ID
	for _, peer := range election.Node.Peers {

		if election.PingPeer(peer) {

			if peer.ID > highestId {
				highestId = peer.ID
			}

		}

	}

	if election.Node.ID == highestId {
		fmt.Println("I ", election.Node.ID, ", am becoming leader !")
		election.Node.BecomeLeader()
	} else {
		fmt.Println("I ", election.Node.ID, ", will serve ", highestId)
		election.Node.BecomeFollower(highestId)
	}
}

func (election *Election) PingPeer(peer Peer) bool {
	return election.RPC.Ping(peer)
}
