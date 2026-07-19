package internal

import "time"

type HeartBeat struct {
	Node *Node
	RPC  *RPC
}

func (this *HeartBeat) SendHeartBeat() {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		<-ticker.C

		if this.Node.IsLeader() {
			this.RPC.BroadcastHeartBeat()
		}
	}
}
