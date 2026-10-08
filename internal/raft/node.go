package raft

import (
	"context"
	"fmt"
	"sync"
	"time"

	pb "distributed-library-kv/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type State int

const (
	Follower State = iota
	Candidate
	Leader
)

type Node struct {
	mu sync.Mutex

	ID string

	State State

	CurrentTerm int
	VotedFor    string

	LastHeartbeat time.Time

	Log         []pb.LogEntry
	CommitIndex int
}

func NewNode(id string) *Node {
	return &Node{
		ID:            id,
		State:         Follower,
		CurrentTerm:   0,
		VotedFor:      "",
		LastHeartbeat: time.Now(),
		Log:           make([]pb.LogEntry, 0),
		CommitIndex:   0,
	}
}

func (n *Node) BecomeCandidate() {
	n.mu.Lock()
	defer n.mu.Unlock()

	n.State = Candidate
	n.CurrentTerm++
	n.VotedFor = n.ID
	n.LastHeartbeat = time.Now()
}
func (n *Node) BecomeLeader() {
	n.mu.Lock()
	defer n.mu.Unlock()

	n.State = Leader
	n.VotedFor = n.ID
	n.LastHeartbeat = time.Now()
}

func (n *Node) BecomeFollower(term int) {
	n.mu.Lock()
	defer n.mu.Unlock()

	n.State = Follower
	n.CurrentTerm = term
	n.VotedFor = ""
	n.LastHeartbeat = time.Now()
}
func (n *Node) HandleRequestVote(term int, candidateID string) (int, bool) {
	n.mu.Lock()
	defer n.mu.Unlock()

	// Reject an old term.
	if term < n.CurrentTerm {
		return n.CurrentTerm, false
	}

	// A newer term makes this node a follower.
	if term > n.CurrentTerm {
		n.CurrentTerm = term
		n.State = Follower
		n.VotedFor = ""
	}

	// Give the vote if we have not voted yet,
	// or if we already voted for this candidate.
	if n.VotedFor == "" || n.VotedFor == candidateID {
		n.VotedFor = candidateID
		n.LastHeartbeat = time.Now()

		return n.CurrentTerm, true
	}

	// Already voted for another candidate.
	return n.CurrentTerm, false
}
func (n *Node) StartElection(nodes []string) bool {
	n.BecomeCandidate()

	n.mu.Lock()
	term := n.CurrentTerm
	candidateID := n.ID
	n.mu.Unlock()

	votes := 1

	selfAddress := ""

	switch n.ID {
	case "node1":
		selfAddress = "node1:50051"
	case "node2":
		selfAddress = "node2:50052"
	case "node3":
		selfAddress = "node3:50053"
	}

	for _, address := range nodes {

		if address == "" || address == selfAddress {
			continue
		}
		ctx, cancel := context.WithTimeout(
			context.Background(),
			2*time.Second,
		)

		conn, err := grpc.NewClient(
			address,
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)

		if err != nil {
			cancel()
			continue
		}

		client := pb.NewKVServiceClient(conn)

		response, err := client.RequestVote(
			ctx,
			&pb.VoteRequest{
				Term:        int32(term),
				CandidateId: candidateID,
			},
		)

		conn.Close()
		cancel()

		if err != nil {
			continue
		}

		// If another node has a newer term,
		// step down immediately.
		if int(response.GetTerm()) > term {
			n.BecomeFollower(int(response.GetTerm()))
			return false
		}

		if response.GetVoteGranted() {
			votes++
		}
	}

	// Become leader only after receiving a majority.
	if votes >= 2 {
		n.mu.Lock()

		// Make sure we are still a candidate
		// in the same election term.
		if n.State == Candidate && n.CurrentTerm == term {
			n.State = Leader
			n.LastHeartbeat = time.Now()
			n.mu.Unlock()

			return true
		}

		n.mu.Unlock()
	}

	n.mu.Lock()
	if n.State == Candidate {
		n.State = Follower
	}
	n.mu.Unlock()

	return false
}
func (n *Node) IsLeader() bool {
	n.mu.Lock()
	defer n.mu.Unlock()

	return n.State == Leader
}
func (n *Node) StartHeartbeat(nodes []string) {
	go func() {
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()

		for range ticker.C {
			if !n.IsLeader() {
				continue
			}

			for _, address := range nodes {
				go n.sendHeartbeat(address)
			}
		}
	}()
}
func (n *Node) sendHeartbeat(address string) {
	ctx, cancel := context.WithTimeout(
		context.Background(),
		2*time.Second,
	)
	defer cancel()

	conn, err := grpc.NewClient(
		address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return
	}
	defer conn.Close()

	client := pb.NewKVServiceClient(conn)

	n.mu.Lock()
	term := n.CurrentTerm
	n.mu.Unlock()

	_, _ = client.Heartbeat(
		ctx,
		&pb.HeartbeatRequest{
			Term:     int32(term),
			LeaderId: n.ID,
		},
	)
}
func (n *Node) HandleAppendEntries(
	term int,
	leaderID string,
	entries []*pb.LogEntry,
	leaderCommit int,
) (int, bool) {

	n.mu.Lock()
	defer n.mu.Unlock()

	// Reject old term.
	if term < n.CurrentTerm {
		return n.CurrentTerm, false
	}

	// Update term if leader has newer term.
	if term > n.CurrentTerm {
		n.CurrentTerm = term
		n.VotedFor = ""
	}

	// Become follower when receiving leader message.
	n.State = Follower
	n.LastHeartbeat = time.Now()

	// Add entries to local log.
	for _, entry := range entries {
		if entry != nil {
			n.Log = append(n.Log, *entry)
		}
	}

	// Update commit index.
	if leaderCommit > n.CommitIndex {
		n.CommitIndex = leaderCommit

		if n.CommitIndex > len(n.Log) {
			n.CommitIndex = len(n.Log)
		}
	}

	return n.CurrentTerm, true
}
func (n *Node) GetLog() []pb.LogEntry {
	n.mu.Lock()
	defer n.mu.Unlock()

	result := make([]pb.LogEntry, len(n.Log))
	copy(result, n.Log)

	return result
}
func (n *Node) ReplicateEntry(
	nodes []string,
	command string,
	key string,
	value string,
) bool {

	n.mu.Lock()

	if n.State != Leader {
		n.mu.Unlock()
		return false
	}

	entry := &pb.LogEntry{
		Term:    int32(n.CurrentTerm),
		Command: command,
		Key:     key,
		Value:   value,
	}

	n.Log = append(n.Log, *entry)

	term := n.CurrentTerm
	n.mu.Unlock()

	// Leader counts as one vote.
	acknowledgements := 1

	resultCh := make(chan bool, len(nodes))

	for _, address := range nodes {
		if address == "" {
			continue
		}

		go func(address string) {

			ctx, cancel := context.WithTimeout(
				context.Background(),
				2*time.Second,
			)
			defer cancel()

			conn, err := grpc.NewClient(
				address,
				grpc.WithTransportCredentials(
					insecure.NewCredentials(),
				),
			)

			if err != nil {
				resultCh <- false
				return
			}
			defer conn.Close()

			client := pb.NewKVServiceClient(conn)

			response, err := client.AppendEntries(
				ctx,
				&pb.AppendEntriesRequest{
					Term:         int32(term),
					LeaderId:     n.ID,
					Entries:      []*pb.LogEntry{entry},
					LeaderCommit: int32(n.CommitIndex),
				},
			)

			if err != nil {
				fmt.Printf("AppendEntries to %s failed: %v\n", address, err)
				resultCh <- false
				return
			}

			fmt.Printf(
				"AppendEntries to %s: success=%v term=%d\n",
				address,
				response.GetSuccess(),
				response.GetTerm(),
			)

			resultCh <- response.GetSuccess()

		}(address)
	}

	// Wait for follower responses.
	for range nodes {
		select {
		case success := <-resultCh:
			if success {
				acknowledgements++
			}

			// Majority achieved.
			if acknowledgements >= 2 {
				n.mu.Lock()

				n.CommitIndex = len(n.Log)

				n.mu.Unlock()

				return true
			}

		case <-time.After(3 * time.Second):
			return false
		}
	}

	return false
}
func (n *Node) ElectionTimeoutExpired(timeout time.Duration) bool {
	n.mu.Lock()
	defer n.mu.Unlock()

	return time.Since(n.LastHeartbeat) > timeout
}
func (n *Node) ShouldStartElection() bool {
	n.mu.Lock()
	defer n.mu.Unlock()

	// Never start an election if already leader.
	if n.State == Leader {
		return false
	}

	// Start election if we have not heard from a leader
	// for more than 5 seconds.
	return time.Since(n.LastHeartbeat) > 5*time.Second
}
