package main

import "time"

type Server struct {
	Conf          Config
	State         ServerState
	LastHeartbeat time.Time
	LeaderID      int
	StateParam    StateParameters
	GotNumVotes   int
}

func (server *Server) ChangeState(state ServerState) {
	server.State = state
	server.ResetHeartbeat()
	if state == LEADER {
		self.LeaderID = self.Conf.ProcessID
	}
}

func (server *Server) ResetHeartbeat() {
	server.LastHeartbeat = time.Now()
}

func (server *Server) SetLeaderID(leaderId int) {
	server.LeaderID = leaderId
}
