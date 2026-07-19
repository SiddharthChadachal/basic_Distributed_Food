package internal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type RPC struct {
	Node     *Node
	Election *Election
}

type HeartbeatRequest struct {
	LeaderId int `json:"leader_id"`
}

func (r *RPC) StartServer() {

	mux := http.NewServeMux()
	mux.HandleFunc("/heartbeat", r.HandleHeartBeat)
	mux.HandleFunc("/status", r.HandleStatus)
	mux.HandleFunc("/reelection", r.HandleReElection)

	addr := fmt.Sprintf(":%d", r.Node.Port)

	fmt.Println("Starting server on")

	go func() {
		if err := http.ListenAndServe(addr, mux); err != nil {
			fmt.Println("[ERROR] ", err)
		}
	}()

}

func (r *RPC) HandleHeartBeat(w http.ResponseWriter, req *http.Request) {

	if req.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var hb HeartbeatRequest

	if err := json.NewDecoder(req.Body).Decode(&hb); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	r.Node.LeaderId = hb.LeaderId
	r.Node.UpdateHeartBeat()
	w.WriteHeader(http.StatusOK)
	fmt.Printf(
		"[Node %d] Received heartbeat from leader %d\n",
		r.Node.ID,
		hb.LeaderId,
	)
}

func (r *RPC) HandleStatus(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-type", "application/json")
	json.NewEncoder(w).Encode(r.Node)
}

func (r *RPC) SendHeartBeat(peer Peer) error {

	body := HeartbeatRequest{
		LeaderId: r.Node.ID,
	}

	data, err := json.Marshal(body)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("http://%s:%d/heartbeat", peer.Host, peer.Port)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(data))

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	return nil
}

func (r *RPC) BroadcastHeartBeat() {
	for _, peer := range r.Node.Peers {

		err := r.SendHeartBeat(peer)
		if err != nil {
			fmt.Println("[ERROR] Unable to send HeartBeat", err)
		}
		fmt.Printf(
			"[Leader %d] Sending heartbeat to node %d\n",
			r.Node.ID,
			peer.ID,
		)
	}
}

func (r *RPC) Ping(peer Peer) bool {
	client := http.Client{
		Timeout: time.Second,
	}
	url := fmt.Sprintf("http://%s:%d/status", peer.Host, peer.Port)

	resp, err := client.Get(url)
	if err != nil {
		fmt.Println("[ERROR] ", err)
		return false
	}

	defer resp.Body.Close()
	return resp.StatusCode == http.StatusOK
}

func (r *RPC) HandleReElection(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(w, "Bruh wrong Method like come on ", http.StatusMethodNotAllowed)
		return
	}

	if !r.Node.IsLeader() {
		http.Error(w, "You cant raise your voice against leader !", http.StatusForbidden)
		return
	}

	fmt.Println("\nMOVE ASIDE RE-ELECTION")
	fmt.Println("\n========== MANUAL RE-ELECTION ==========")

	fmt.Println("LEADER GOING DOWN !")
	r.Node.LeaderId = -1
	go r.Election.StartElection()
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Manual Election"))
}
