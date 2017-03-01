package main

import (
	"fmt"
	"net/http"
	"net/rpc"
	"time"
)


type Server struct {
	Conf          Config
	State         ServerState
	LastHeartbeat time.Time

	StateParam StateParameters

	ElectionTimestamp time.Time
	GotNumVotes       int
}

var self Server

func init() {
	configuration := ReadConfig()
	stateParams := readSavedState()

	self = Server{
		Conf:        configuration,
		State:       FOLLOWER,
		StateParam:  stateParams,
		GotNumVotes: 0,
	}

	self.LastHeartbeat = time.Now()
}

func main() {

	clientComm := new(ClientComm)
	dataCenterComm := new(DataCenterComm)

	clientComm.value = 1000
	rpc.Register(clientComm)
	rpc.Register(dataCenterComm)
	rpc.HandleHTTP()

	//go EstablishConnections()
	go runStateMachine()
	go waitUserInput()

	err := http.ListenAndServe(self.Conf.MyAddress, nil)
	if err != nil {
		fmt.Println(err.Error())
	}
}

func runStateMachine() {
	for {
		fmt.Println(self.State)
		switch self.State {
		case FOLLOWER:
			followerBehavior()
		case CANDIDATE:
			candidateBehavior()
		case LEADER:
		}
		time.Sleep(1000 * time.Millisecond)
	}
}

func followerBehavior() {
	if time.Since(self.LastHeartbeat) > self.Conf.Timeout {
		// Follower timeout, convert to candidate
		self.State = CANDIDATE
		startElection()
		return
	}
	
	checkAndUpdateLogs()
}

func startElection() {
	// Increment currentTerm
	self.StateParam.CurrentTerm++

	// Vote for self
	self.StateParam.VotedFor = self.Conf.ProcessID
	self.GotNumVotes = 1

	// Reset election timer
	self.ElectionTimestamp = time.Now()

	// TODO: Send RequestVote RPCs to all other servers
	
}

func receivedMajorityVotes() bool {
	if self.GotNumVotes > self.Conf.NumMajority {
		return true
	}
	return false
}

func candidateBehavior() {

	// If votes received from majority of servers: become leader
	if receivedMajorityVotes() {
		self.State = LEADER
		/*
		Upon election: send initial empty AppendEntries RPCs
		(heartbeat) to each server; repeat during idle periods to
		prevent election timeouts (§5.2)
		*/
		go sendAppendEntries()
		return
	}

	// TODO: If AppendEntries RPC received from new leader: convert to follower
	// Do this in AppendEntries RPC handler

	// If election timeout elapses: start new election
	if time.Since(self.ElectionTimestamp) > self.Conf.Timeout {
		startElection()
	}
}

func leaderBehavior() {
	
	// TODO: If command received from client: append entry to local log,
	// respond after entry applied to state machine (§5.3)
	// Handled in ClientRPC handler

	// If last log index ≥ nextIndex for a follower: send
	// AppendEntries RPC with log entries starting at nextIndex

	// If successful: update nextIndex and matchIndex for
	// follower (§5.3)

	// If AppendEntries fails because of log inconsistency:
	// decrement nextIndex and retry (§5.3)

	// If there exists an N such that N > commitIndex, a majority
	// of matchIndex[i] ≥ N, and log[N].term == currentTerm:
	// set commitIndex = N (§5.3, §5.4).
	
	checkAndUpdateLogs()
}

func sendAppendEntries() {

}

func checkAndUpdateLogs() {
	// If commitIndex > lastApplied: increment lastApplied,
	if self.StateParam.CommitIndex > self.StateParam.LastApplied {
		// apply log[lastApplied] to state machine (§5.3)
		self.Conf.RemainingTickets -= self.StateParam.Logs[self.StateParam.LastApplied]
	}
	
	// TODO: If RPC request or response contains term T > currentTerm:
	// set currentTerm = T, convert to follower (§5.1)
	// Handled in AppendEntriesRPC handler
}