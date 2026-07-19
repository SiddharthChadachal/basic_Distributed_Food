package main

import (
	"basic_distributed_food/internal"
	"flag"
)

func main() {
	id := flag.Int("id", 1, "Node ID")
	flag.Parse()

	config := internal.LoadConfig(*id)
	node := internal.NewNode(config)

	rpc := &internal.RPC{
		Node: node,
	}

	heartbeat := &internal.HeartBeat{
		Node: node,
		RPC:  rpc,
	}

	election := &internal.Election{
		Node: node,
		RPC:  rpc,
	}

	rpc.StartServer()

	go heartbeat.SendHeartBeat()
	go election.StartChecking()

	rpc.Election = election
	select {}
}
