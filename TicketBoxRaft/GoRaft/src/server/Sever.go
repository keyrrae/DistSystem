package main

import (
	"time"
	_ "fmt"
	"encoding/json"
	"io/ioutil"
)

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

func (server *Server) ResetPeers(){
	for _, peer := range server.Conf.Peers {
		peer.NextIndex = len(server.StateParam.Logs)
		peer.MatchedIndex = -1
	}
}

func (server *Server) WriteToStorage() {
	stateJson, err := json.Marshal(server.StateParam)
	
	err = ioutil.WriteFile("./saved_state.json", stateJson, 0755)
	check(err)
}

func (server *Server) ApplyLogsToStateMachine() {
	self.Conf.RemainingTickets = self.Conf.InitialTktNum
	for i := 0; i <= self.StateParam.CommitIndex; i++ {
		self.Conf.RemainingTickets -= self.StateParam.Logs[i].Num
	}
}

func check(err error){
	if err != nil {
		panic(err)
	}
}