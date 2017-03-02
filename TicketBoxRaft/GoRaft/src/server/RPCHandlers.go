package main

/*
type handlers interface {
	RequestVoteHandler(req *DataCenterRequest,
		reply *DataCenterReply) error
	AppendEntriesHandler() float64
}*/

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

type BuyTicketReply struct{
	Success bool
	Remains int
}

type DataCenterComm struct {
	value int
}

type DataCenterReply struct {
}

func (t *ClientComm) BuyTicketHandler(req *BuyTicketRequest, reply *BuyTicketReply) error {
	return nil
	
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
	
	correctID := ( self.StateParam.VotedFor == -1 ||
		self.StateParam.VotedFor == req.CandidateId )
	
	satisfactoryTerm := (req.Term >= self.StateParam.CurrentTerm)
	if correctID && satisfactoryTerm {
		reply.VoteGranted = true
		// TODO: check how to reply term
		reply.Term = self.StateParam.CurrentTerm
	}

	return nil
}

type AppendEntriesRequest struct {
	Term         int        // leader’s term
	LeaderId     int        // so follower can redirect clients
                            // prevLogIndex index of log entry immediately preceding new ones
	PrevLogTerm  int        // term of prevLogIndex entry
	Entries      []LogEntry // log entries to store
                            // (empty for heartbeat; may send more than one for efficiency)
	LeaderCommit int        //leader’s commitIndex
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
	
	// TODO: For a candidate:
	// if AppendEntries RPC received from new leader: convert to follower
	if self.State == CANDIDATE {
		self.ChangeState(FOLLOWER)
	}
	
	if self.State == FOLLOWER {
		self.ResetHeartbeat()
	}
	if req.Term < self.StateParam.CurrentTerm {
		reply.Success = false
	}
	
	/*
	hasEntry := false
	for i := 0; i < len(req.Entries); i++{
		if req.Entries[i].Term ==
	}
	*/
	
	if req.LeaderCommit > self.StateParam.CommitIndex {
		self.StateParam.CommitIndex = min(req.LeaderCommit, len(req.Entries) - 1)
	}
	
	return nil
}