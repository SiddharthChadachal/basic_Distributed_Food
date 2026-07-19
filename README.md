# Mini Distributed Food Cluster

A distributed systems project in Go demonstrating leader election, heartbeats, and failure detection using HTTP-based RPC.

## Features

- Leader Election (Highest ID)
- Heartbeat Mechanism
- Automatic Leader Failover
- Failure Detection
- Manual Re-election Endpoint
- Multi-node Simulation
- HTTP RPC Communication

## Architecture

```
          +---------+
          | Leader  |
          +---------+
           /   |   \
          /    |    \
         /     |     \
+---------+ +---------+ +---------+
|Follower | |Follower | |Follower |
+---------+ +---------+ +---------+
```

## Run

Terminal 1

```bash
go run cmd/node/main.go --id=1
```

Terminal 2

```bash
go run cmd/node/main.go --id=2
```

Terminal 3

```bash
go run cmd/node/main.go --id=3
```

## Trigger Re-election

```bash
curl -X POST http://localhost:8802/reelect
```

## Future Improvements

- UUID-based node IDs
- Cluster Bootstrap
- Service Discovery
- voting system for leader election
- Gossip architecture
- Node Weights to provide realtime 
