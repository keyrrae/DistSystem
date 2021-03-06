package main

import (
	"log"
)

type RequestVoteRequest struct {
	Term         int // candidate’s term
	CandidateId  int // candidate requesting vote
	LastLogIndex int // index of candidate’s last log entry
	LastLogTerm  int // term of candidate’s last log entry
}

type RequestVoteReply struct {
	Term        int  // current term, for candidate to update itself
	VoteGranted bool // true means candidate received vote
}

type ClientComm struct {
	value int
}

type DataCenterComm struct {
	value int
}

type DataCenterReply struct {
}

func (t *DataCenterComm) RequestVoteHandler(req *RequestVoteRequest,
	reply *RequestVoteReply) error {
	/*
		1. Reply false if term < currentTerm (§5.1)
		2. If votedFor is null or candidateId, and candidate’s log is at
		least as up-to-date as receiver’s log, grant vote (§5.2, §5.4)
	*/
	if req.Term < self.StateParam.CurrentTerm {
		reply.VoteGranted = false
		// TODO: check how to reply term
		reply.Term = self.StateParam.CurrentTerm
		return nil
	}

	if req.Term > self.StateParam.CurrentTerm{
		//reply.VoteGranted = true
		//reply.Term = self.StateParam.CurrentTerm
		self.StateParam.VotedFor = -1
		self.StateParam.CurrentTerm = req.Term
		if self.State == LEADER {
			self.ChangeState(FOLLOWER)
		} else if self.State == CANDIDATE {
			self.ChangeState(FOLLOWER)
		}
	}

	correctID := (self.StateParam.VotedFor == -1 ||
		self.StateParam.VotedFor == req.CandidateId)

	myLastIndex := len(self.StateParam.Logs) - 1
	satisfactoryLogIndex := (req.LastLogIndex >= myLastIndex)

	myLastTerm := 0
	if myLastIndex >= 0 {
		myLastTerm = self.StateParam.Logs[myLastIndex].Term
	}
	satisfactoryLogTerm := (req.LastLogTerm >= myLastTerm)

	log.Println("term", req.LastLogTerm, myLastTerm)

	log.Println("conditions", correctID, satisfactoryLogIndex, satisfactoryLogTerm)

	if correctID && satisfactoryLogIndex && satisfactoryLogTerm {
		reply.VoteGranted = true
		reply.Term = self.StateParam.CurrentTerm
		self.StateParam.VotedFor = req.CandidateId
	} else {
		reply.VoteGranted = false
		reply.Term = self.StateParam.CurrentTerm
	}

	return nil
}

type AppendEntriesRequest struct {
	Term         int // leader’s term
	LeaderId     int // so follower can redirect clients
	PrevLogIndex int // prevLogIndex index of log entry immediately preceding new ones

	PrevLogTerm int        // term of prevLogIndex entry
	Entries     []LogEntry // log entries to store
	// (empty for heartbeat; may send more than one for efficiency)
	LeaderCommit int //leader’s commitIndex
}

type AppendEntriesReply struct {
	Term    int  // currentTerm, for leader to update itself
	Success bool //true if follower contained entry matching PrevLogIndex and PrevLogTerm
}

type handlers interface {
}

func (stateParams StateParameters) GetLastLogEntryTerm() int {
	lenOfLogs := len(stateParams.Logs)
	if lenOfLogs == 0 {
		return 0
	}
	return stateParams.Logs[lenOfLogs-1].Term
}

func (t *DataCenterComm) AppendEntriesHandler(req *AppendEntriesRequest,
	reply *AppendEntriesReply) error {
	/*
		1. Reply false if term < currentTerm (§5.1)
		2. Reply false if log doesn’t contain an entry at prevLogIndex
		whose term matches prevLogTerm (§5.3)
		3. If an existing entry conflicts with a new one (same index
		but different terms), delete the existing entry and all that
		follow it (§5.3)
		4. Append any new entries not already in the log
		5. If leaderCommit > commitIndex, set commitIndex =
			min(leaderCommit, index of last new entry)
	*/

	// For a candidate:
	// if AppendEntries RPC received from new leader: convert to follower
	if self.State == CANDIDATE {
		self.ChangeState(FOLLOWER)
		self.ResetHeartbeat()
		self.SetLeaderID(req.LeaderId)
		self.StateParam.CurrentTerm = req.Term
	}

	if self.State == FOLLOWER {
		self.ResetHeartbeat()
		self.SetLeaderID(req.LeaderId)
		if req.Term > self.StateParam.CurrentTerm {
			self.StateParam.CurrentTerm = req.Term
		}
	}

	if self.State == LEADER {
		if req.Term > self.StateParam.CurrentTerm {
			self.ChangeState(FOLLOWER)
			self.ResetHeartbeat()
			self.SetLeaderID(req.LeaderId)
		}
	}

	if req.Term < self.StateParam.CurrentTerm {
		reply.Success = false
		reply.Term = self.StateParam.CurrentTerm
		return nil
	}

	// If an existing entry conflicts with a new one (same index
	// but different terms), delete the existing entry and all that
	// follow it (§5.3)

	if req.PrevLogIndex < 0 {
		reply.Success = true
	} else if req.PrevLogIndex > len(self.StateParam.Logs)-1 {
		reply.Success = false
	} else {
		reply.Success = (self.StateParam.Logs[req.PrevLogIndex].Term == req.PrevLogTerm)
	}

	for i := req.PrevLogIndex + 1; i < req.PrevLogIndex+1+len(req.Entries); i++ {
		logIndex := i - req.PrevLogIndex - 1
		if len(self.StateParam.Logs)-1 < i {
			// Append any new entries not already in the log
			self.StateParam.Logs = append(self.StateParam.Logs, req.Entries[logIndex])
		} else {
			// If an existing entry conflicts with a new one (same index
			// but different terms), delete the existing entry and all that
			// follow it (§5.3)
			if self.StateParam.Logs[i].Term != req.Entries[logIndex].Term {
				self.StateParam.Logs = self.StateParam.Logs[:i]
				//self.StateParam.Logs = append(self.StateParam.Logs[:i], req.Entries[logIndex])
			}
		}
	}

	// If leaderCommit > commitIndex,
	// set commitIndex = min(leaderCommit, index of last new entry)

	if req.LeaderCommit > self.StateParam.CommitIndex {
		self.StateParam.CommitIndex = min(req.LeaderCommit, len(self.StateParam.Logs)-1)
	}

	return nil
}
