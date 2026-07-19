package internal

import (
	"fmt"

	"time"
)

func (p Peer) GetAddress() string {
	return fmt.Sprintf("%s:%d", p.Host, p.Port)
}

type Config struct {
	ID                int
	Peers             []Peer
	Host              string
	Port              int
	HeartbeatInterval time.Duration
	ElectionTimeout   time.Duration
}

func LoadConfig(id int) Config {

	allNodes := []Peer{
		{ID: 1, Host: "localhost", Port: 8801},
		{ID: 2, Host: "localhost", Port: 8802},
		{ID: 3, Host: "localhost", Port: 8803},
	}

	var self Peer
	var peers []Peer

	for _, node := range allNodes {
		if node.ID == id {
			self = node
		} else {
			peers = append(peers, node)
		}

	}

	return Config{
		ID:                self.ID,
		Peers:             peers,
		Host:              self.Host,
		Port:              self.Port,
		HeartbeatInterval: 1 * time.Second,
		ElectionTimeout:   5 * time.Second,
	}
}
