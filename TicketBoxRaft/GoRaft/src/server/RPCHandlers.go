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

func (t *DataCenterComm) RequestVoteHandler(req *RequestVoteRequest,
	reply *RequestVoteReply) error {
	/*
		1. Reply false if term < currentTerm (§5.1)
		2. If votedFor is null or candidateId, and candidate’s log is at
		least as up-to-date as receiver’s log, grant vote (§5.2, §5.4)
	*/
	if req.Term < stateParams.CurrentTerm {
		reply.VoteGranted = false
		// TODO: check how to reply term
		reply.Term = stateParams.CurrentTerm
		return nil
	}

	if stateParams.VotedFor == -1 || stateParams.VotedFor == req.CandidateId {
		reply.VoteGranted = true
		// TODO: check how to reply term
		reply.Term = stateParams.CurrentTerm
	}

	return nil
}

type AppendEntriesRequest struct {
	Term         int        // leader’s term
	LeaderId     int        // so follower can redirect clients prevLogIndex index of log entry immediately preceding new ones
	PrevLogTerm  int        // term of prevLogIndex entry
	Entries      []LogEntry // log entries to store (empty for heartbeat; may send more than one for efficiency)
	LeaderCommit int        //leader’s commitIndex
}

type LogEntry struct {
	Num  int
	Term int
}

type AppendEntriesReply struct {
	Term    int  // currentTerm, for leader to update itself
	Success bool //true if follower contained entry matching PrevLogIndex and PrevLogTerm
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
	if req.Term < stateParams.CurrentTerm {
		reply.Success = false
	}
	/*
	hasEntry := false
	for i := 0; i < len(req.Entries); i++{
		if req.Entries[i].Term ==
	}
	*/
	
	if req.LeaderCommit > stateParams.CommitIndex {
		stateParams.CommitIndex = min(req.LeaderCommit, req.Entries[-1])
	}
	

	return nil
}