package main

import (
	"log"
	"time"
	"fmt"
	"net/rpc"
)

func runStateMachine() {
	for {
		log.Println(self.State)
		switch self.State {
		case FOLLOWER:
			followerBehavior()
		case CANDIDATE:
			candidateBehavior()
		case LEADER:
			leaderBehavior()
		}
		time.Sleep(1000 * time.Millisecond)
	}
}

func followerBehavior() {
	if time.Since(self.LastHeartbeat) > self.Conf.Timeout {
		// Follower timeout, convert to candidate
		self.ChangeState(CANDIDATE)
		log.Println(self.State)
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
	
	// TODO: Send RequestVote RPCs to all other servers
	fmt.Println(len(self.Conf.Peers))
	
	done := make(chan bool)
	
	for _, peer := range self.Conf.Peers {
		if !tryEstablishConnection(&peer) {
			continue
		}
		
		go func(peer Peer) {
			// Asynchronous call
			requestVoteRequest := RequestVoteRequest{
				Term:         self.StateParam.CurrentTerm,
				CandidateId:  self.Conf.ProcessID,
				LastLogIndex: self.StateParam.LastApplied,
				
				// TODO: last candidate's log's term
				LastLogTerm: 0,
			}
			
			lenOfLogs := len(self.StateParam.Logs)
			if lenOfLogs > 0 {
				requestVoteRequest.LastLogTerm = self.StateParam.Logs[lenOfLogs-1].Term
			}
			
			requestVoteReply := new(RequestVoteReply)
			
			if peer.Comm == nil {
				fmt.Println("client == nil")
			}
			
			err := peer.Comm.Call("DataCenterComm.RequestVoteHandler", requestVoteRequest, &requestVoteReply)
			if err != nil {
				peer.Comm = nil
				peer.Connected = false
				return
			}
			
			if requestVoteReply.VoteGranted {
				self.GotNumVotes++
				log.Printf("GotNumVotes: %v", self.GotNumVotes)
			}
			
			if receivedMajorityVotes() {
				done <- true
			}
		}(peer)
	}
	
	go func() {
		<-done
		self.ChangeState(LEADER)
	}()
}

func tryEstablishConnection(peer *Peer) bool {
	if !peer.Connected {
		
		client, err := rpc.DialHTTP("tcp", peer.Address)
		
		if err == nil {
			peer.Comm = client
			peer.Connected = true
			return true
		} else {
			log.Printf("cannot reach peer, %s\n", peer.Address)
			return false
		}
	}
	return true
}

func receivedMajorityVotes() bool {
	if self.GotNumVotes > self.Conf.NumMajority {
		return true
	}
	return false
}

func candidateBehavior() {
	
	// If AppendEntries RPC received from new leader: convert to follower
	// This is handled in AppendEntries RPC handler
	
	// If election timeout elapses: start new election
	if time.Since(self.LastHeartbeat) > self.Conf.Timeout {
		startElection()
	}
}

func leaderBehavior() {
	
	// TODO: If command received from client: append entry to local log,
	// respond after entry applied to state machine (§5.3)
	// Handled in ClientRPC handler
	
	sendAppendEntriesToAll()
	
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

func sendAppendEntriesToAll() {
	
	for _, peer := range self.Conf.Peers {
		if !tryEstablishConnection(&peer) {
			continue
		}
		
		// Send an append entry rpc to peer
		go sendAppendEntriesToPeer(peer)
	}
}

func sendAppendEntriesToPeer(peer Peer) {
	// Asynchronous call
	appendEntriesRequest := AppendEntriesRequest{
		Term: self.StateParam.CurrentTerm,
		LeaderId: self.Conf.ProcessID,
		PrevLogTerm: self.StateParam.GetLastLogEntryTerm(),
	}
	
	// TODO: update appendEntriesRequest.Entries according to peer condition
	
	appendEntriesReply := new(AppendEntriesReply)
	
	err := peer.Comm.Call("DataCenterComm.AppendEntriesHandler", appendEntriesRequest, &appendEntriesReply)
	if err != nil {
		peer.Comm = nil
		peer.Connected = false
		return
	}
}

func checkAndUpdateLogs() {
	// If commitIndex > lastApplied: increment lastApplied,
	if self.StateParam.CommitIndex > self.StateParam.LastApplied {
		// apply log[lastApplied] to state machine (§5.3)
		self.Conf.RemainingTickets -= self.StateParam.Logs[self.StateParam.LastApplied].Num
	}
	
	// TODO: If RPC request or response contains term T > currentTerm:
	// set currentTerm = T, convert to follower (§5.1)
	// Handled in AppendEntriesRPC handler
}
