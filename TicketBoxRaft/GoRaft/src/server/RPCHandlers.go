package main

import (
	"os"
	"log"
	"time"
	"errors"
	"encoding/json"
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

type BuyTicketRequest struct {
	NumTickets int
}

type BuyTicketReply struct {
	Success bool
	Remains int
}

type DataCenterComm struct {
	value int
}

type DataCenterReply struct {
}

type ChangeConfigRequest struct {
	Servers []byte
}

type ChangeConfigReply struct {
	Success bool
}

func (t *ClientComm) ChangeConfigHandler(req *ChangeConfigRequest, reply *ChangeConfigReply) error {
	log.Println("received configuration change request from a client")

	if self.LeaderID != self.Conf.ProcessID {
		leader := self.Conf.PeersMap[self.LeaderID]

		leaderReply := new(ChangeConfigReply)

		if leader == nil || leader.Comm == nil {
			reply.Success = false
			return errors.New("No leader in the cluster yet")
		}

		err := leader.Comm.Call("DataCenterComm.ChangeConfigHandler", req, &leaderReply)
		if err != nil {
			leader.Comm = nil
			leader.Connected = false
			reply.Success = false
			return err
		}
		reply.Success = leaderReply.Success
	} else {
		performConfigurationChange(req, reply)
	}

	return nil
}

func (t *DataCenterComm) ChangeConfigHandler(req *ChangeConfigRequest, reply *ChangeConfigReply) error {
	log.Println("received configuration change request from a follower redirection")

	if self.LeaderID == self.Conf.ProcessID {
		performConfigurationChange(req, reply)
	} else {
		reply.Success = false
	}
	return nil
}

func performConfigurationChange(req *ChangeConfigRequest, reply *ChangeConfigReply) {
	var newConfig []Peer

	err := json.Unmarshal(req.Servers, &newConfig)
	if err != nil{
		log.Println("parse log error")
	}

	// Form a joint consensus configuration
	addressMap := make(map[string]bool)

	addressMap[self.Conf.MyAddress] = true
	for _, oldPeer := range self.Conf.Peers{
		addressMap[oldPeer.Address] = true
	}

	var jointConfig []Peer

	jointConfig = append(jointConfig, Peer{Address:self.Conf.MyAddress, ProcessId:self.Conf.ProcessID})

	for _, peer := range self.Conf.Peers {
		newPeer := Peer{
			Address:peer.Address,
			ProcessId: peer.ProcessId,
			MatchedIndex: -1,
			NextIndex: 0,
		}
		jointConfig = append(jointConfig, newPeer)
	}


    shouldStay := false
	for _, peer := range newConfig{
		if peer.Address == self.Conf.MyAddress{
			shouldStay = true
		}
		if _, ok := addressMap[peer.Address]; !ok {
			newPeer := Peer{
				Address:peer.Address,
				ProcessId: peer.ProcessId,
				MatchedIndex: -1,
				NextIndex: 0,
			}
			self.Conf.Peers = append(self.Conf.Peers, &newPeer)
			jointConfig = append(jointConfig, newPeer)
			addressMap[peer.Address] = true
		}
	}

	self.Conf.NumMajority = len(self.Conf.Peers) / 2 + 1

	if err != nil{
		log.Println("convert to json failed")
	}

	jointConfigJson, err := json.Marshal(jointConfig)
	// append new config as an entry
	jointConfigLogEntry := LogEntry{
		Num:  0,
		Term: self.StateParam.CurrentTerm,
		IsConfigurationChange: true,
		NewConfig: string(jointConfigJson),

		//NewConfig: string(req.Servers),
	}

	self.StateParam.Logs = append(self.StateParam.Logs, jointConfigLogEntry)

	leaderBehavior()
	for{
		if self.StateParam.CommitIndex == self.StateParam.LastApplied{
			break
		}
		time.Sleep(100 * time.Millisecond)
	}

	// Change to new configuration
	newConfigAddressMap := make(map[string]bool)
	for _, serverFromNewConfig := range newConfig{
		newConfigAddressMap[serverFromNewConfig.Address] = true
	}

	for i, serverFromJointConfig := range self.Conf.Peers {
		if _, ok := newConfigAddressMap[serverFromJointConfig.Address]; !ok{
			self.Conf.Peers = append(self.Conf.Peers[:i],
				self.Conf.Peers[i+1:]...)
		}
	}

	self.Conf.NumMajority = ( len(self.Conf.Peers) + 1 ) / 2 + 1

	//newConf, err := json.Marshal(req.Servers)

	// append new config as an entry
	newConfigLogEntry := LogEntry{
		Num:  0,
		Term: self.StateParam.CurrentTerm,
		IsConfigurationChange: true,
		NewConfig: string(req.Servers),

		//NewConfig: string(req.Servers),
	}

	self.StateParam.Logs = append(self.StateParam.Logs, newConfigLogEntry)

	leaderBehavior()
	for{
		if self.StateParam.CommitIndex == self.StateParam.LastApplied{
			break
		}
		time.Sleep(100 * time.Millisecond)
	}

	reply.Success = true

	if !shouldStay{
		log.Println("Leader not in the new Configuration")
		log.Println("Stepping down....")
		os.Exit(0)
	}
}

func (t *ClientComm) BuyTicketHandler(req *BuyTicketRequest, reply *BuyTicketReply) error {
	log.Println("received buy ticket from a client")

	if self.LeaderID != self.Conf.ProcessID {

		leader := self.Conf.PeersMap[self.LeaderID]
		leaderReply := new(BuyTicketReply)

		if leader == nil || leader.Comm == nil {
			reply.Success = false
			return errors.New("No leader in the cluster yet")
		}

		err := leader.Comm.Call("DataCenterComm.BuyTicketHandler", req, &leaderReply)
		if err != nil {
			leader.Comm = nil
			leader.Connected = false
			reply.Success = false
			reply.Remains = self.StateParam.RemainingTickets
			return err
		}

		reply.Remains = leaderReply.Remains
		reply.Success = leaderReply.Success
	} else {
		sendReqToFollowers(req, reply)
	}

	return nil
}

func (t *DataCenterComm) BuyTicketHandler(req *BuyTicketRequest, reply *BuyTicketReply) error {
	log.Println("received buy ticket from a follower redirection")

	if self.LeaderID == self.Conf.ProcessID {
		sendReqToFollowers(req, reply)
	} else {
		reply.Success = false
		reply.Remains = self.StateParam.RemainingTickets
	}

	return nil
}

func sendReqToFollowers(req *BuyTicketRequest, reply *BuyTicketReply) {
	// I'm the leader
	logEntry := LogEntry{
		Num:  req.NumTickets,
		Term: self.StateParam.CurrentTerm,
	}

	if req.NumTickets > self.StateParam.RemainingTickets {
		reply.Success = false
		reply.Remains = self.StateParam.RemainingTickets
		return
	}
	self.StateParam.Logs = append(self.StateParam.Logs, logEntry)

	leaderBehavior()

	reply.Success = true
	reply.Remains = self.StateParam.RemainingTickets
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

	correctID := (self.StateParam.VotedFor == -1 ||
		self.StateParam.VotedFor == req.CandidateId)

	satisfactoryTerm := (req.Term >= self.StateParam.CurrentTerm)
	if correctID && satisfactoryTerm {
		reply.VoteGranted = true
		reply.Term = self.StateParam.CurrentTerm
		self.StateParam.VotedFor = req.CandidateId
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

	//log.Println("req.PrevLogIndex", req.PrevLogIndex)
	//log.Println("req.Entries", req.Entries)
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
		//log.Println(i, self.StateParam.Logs)
	}

	// If leaderCommit > commitIndex,
	// set commitIndex = min(leaderCommit, index of last new entry)

	if req.LeaderCommit > self.StateParam.CommitIndex {
		self.StateParam.CommitIndex = min(req.LeaderCommit, len(self.StateParam.Logs)-1)
	}

	return nil
}
