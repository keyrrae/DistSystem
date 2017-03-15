package main

import (
	"fmt"
	"log"
	"net/rpc"
	"time"
)

func runStateMachine() {
	for {
		log.Println("Role:", self.State)
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
	for _, peer := range self.Conf.Peers {
		tryEstablishConnection(peer)
	}

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
	fmt.Println("Term", self.StateParam.CurrentTerm)

	// Vote for self
	self.StateParam.VotedFor = self.Conf.ProcessID
	self.GotNumVotes = 1

	// TODO: Send RequestVote RPCs to all other servers

	done := make(chan bool)

	for _, peer := range self.Conf.Peers {
		if !tryEstablishConnection(peer) {
			continue
		}

		go func(peer *Peer) {
			// Asynchronous call
			requestVoteRequest := RequestVoteRequest{
				Term:         self.StateParam.CurrentTerm,
				CandidateId:  self.Conf.ProcessID,
				LastLogIndex: self.StateParam.LastApplied,
				LastLogTerm: 0,
			}

			lenOfLogs := len(self.StateParam.Logs)
			if lenOfLogs > 0 {
				requestVoteRequest.LastLogTerm = self.StateParam.Logs[lenOfLogs-1].Term
			}

			requestVoteReply := new(RequestVoteReply)

			if peer.Comm == nil {
				return
			}

			err := peer.Comm.Call("DataCenterComm.RequestVoteHandler", requestVoteRequest, &requestVoteReply)
			if err != nil {
				peer.Comm = nil
				peer.Connected = false
				return
			}

			if requestVoteReply.VoteGranted {
				self.GotNumVotes++
				log.Printf("Received %v votes", self.GotNumVotes)
			} else {
				self.ChangeState(FOLLOWER)
				self.StateParam.CurrentTerm = requestVoteReply.Term
				self.StateParam.VotedFor = peer.ProcessId
				done <- false
				return
			}

			if receivedMajorityVotes() {
				done <- true
			}
		}(peer)
	}

	go func() {
		leaderGranted := <-done
		if leaderGranted{
			self.ChangeState(LEADER)
		}

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
			log.Printf("Cannot reach peer, %s\n", peer.Address)
			return false
		}
	}
	return true
}

func receivedMajorityVotes() bool {
	if self.GotNumVotes >= self.Conf.NumMajority {
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

	// respond after entry applied to state machine (§5.3)
	// Handled in ClientRPC handler

	sendAppendEntriesToAll()

	// If there exists an N such that N > commitIndex, a majority
	// of matchIndex[i] ≥ N, and log[N].term == currentTerm:
	// set commitIndex = N (§5.3, §5.4).

	updateCommitIndex()

	// If commitIndex > lastApplied: increment lastApplied, apply
	// log[lastApplied] to state machine (§5.3)

	checkAndUpdateLogs()
}

func updateCommitIndex() {
	matchIndexMap := make(map[int]int) // <matchIndex, count>
	for _, peer := range self.Conf.Peers {
		matchIndex := peer.MatchedIndex
		if matchIndex < 0 {
			continue
		}
		if _, ok := matchIndexMap[matchIndex]; ok {
			matchIndexMap[matchIndex]++
		} else {
			matchIndexMap[matchIndex] = 1
		}
	}

	for k, v := range matchIndexMap {
		if v + 1 >= self.Conf.NumMajority {
			for i := self.StateParam.CommitIndex + 1; i <= k; i++ {
				if self.StateParam.Logs[i].Term == self.StateParam.CurrentTerm {
					self.StateParam.CommitIndex = k
				}
			}
		}
	}

}

func sendAppendEntriesToAll() {

	done := make(chan bool, len(self.Conf.Peers))

	for _, peer := range self.Conf.Peers {
		if !tryEstablishConnection(peer) {
			done <- true
			continue
		}

		// Send an append entry rpc to peer
		go sendAppendEntriesToPeer(peer, done)
	}
	for i := 0; i < len(self.Conf.Peers); i++ {
		<-done
	}

}

func sendAppendEntriesToPeer(peer *Peer, done chan<- bool) {
	// Asynchronous call

	appendEntriesRequest := AppendEntriesRequest{
		Term:         self.StateParam.CurrentTerm,
		LeaderId:     self.Conf.ProcessID,
		PrevLogIndex: peer.MatchedIndex,
		LeaderCommit: self.StateParam.CommitIndex,
	}

	if len(self.StateParam.Logs) > 0 {
		if peer.MatchedIndex >= 0 {
			appendEntriesRequest.PrevLogTerm =
				self.StateParam.Logs[peer.MatchedIndex].Term
		} else {
			appendEntriesRequest.PrevLogTerm = self.StateParam.CurrentTerm
		}

		for i := peer.NextIndex; i < len(self.StateParam.Logs); i++ {
			appendEntriesRequest.Entries = append(appendEntriesRequest.Entries,
				self.StateParam.Logs[i])
		}
	} else {
		appendEntriesRequest.PrevLogTerm = 0
	}

	appendEntriesReply := new(AppendEntriesReply)

	err := peer.Comm.Call("DataCenterComm.AppendEntriesHandler",
		appendEntriesRequest, &appendEntriesReply)
	if err != nil {
		log.Printf("Cannot reach peer, %s\n", peer.Address)
		peer.Comm = nil
		peer.Connected = false
		done <- true
		return
	}

	if appendEntriesReply.Success {
		// If successful: update nextIndex and matchIndex for follower (§5.3)
		peer.MatchedIndex = len(self.StateParam.Logs) - 1
		peer.NextIndex = peer.MatchedIndex + 1
	} else {
		if appendEntriesReply.Term > self.StateParam.CurrentTerm {
			// If RPC request or response contains term T > currentTerm:
			// set currentTerm = T, convert to follower (§5.1)

			self.StateParam.CurrentTerm = appendEntriesReply.Term
			self.ChangeState(FOLLOWER)
			self.SetLeaderID(peer.ProcessId)
			done <- true
			return
		} else {
			// If AppendEntries fails because of log inconsistency:
			// decrement nextIndex and retry (§5.3)
			peer.NextIndex--
			sendAppendEntriesToPeer(peer, done) //retry
		}
	}
	done <- true
}

func checkAndUpdateLogs() {
	// If commitIndex > lastApplied: increment lastApplied,


	if self.StateParam.CommitIndex > self.StateParam.LastApplied {
		// apply log[lastApplied] to state machine (§5.3)

		self.ApplyLogsToStateMachine()

		self.StateParam.LastApplied = self.StateParam.CommitIndex
	}

	self.PrintLogs()
	self.WriteToStorage()
	log.Println("AppliedIndex:", self.StateParam.LastApplied)
	log.Println("CommitIndex :", self.StateParam.CommitIndex)
	log.Println("Tickets     :", self.StateParam.RemainingTickets)
}
