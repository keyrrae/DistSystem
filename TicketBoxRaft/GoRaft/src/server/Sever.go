package main

import (
	"time"
	_ "fmt"
	"log"
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
		server.LeaderID = server.Conf.ProcessID
		self.ResetPeers()
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
		peer.NextIndex = 0//len(server.StateParam.Logs)
		peer.MatchedIndex = -1
	}
}

func (server *Server) WriteToStorage() {
	stateJson, err := json.Marshal(server.StateParam)
	
	err = ioutil.WriteFile("./saved_state.json", stateJson, 0755)
	check(err)
}

func (server *Server) ApplyLogsToStateMachine() {
	server.StateParam.RemainingTickets = server.Conf.InitialTktNum
	if server.StateParam.CommitIndex > len(server.StateParam.Logs) - 1{
		return
	}
	for i := 0; i <= server.StateParam.CommitIndex; i++ {
		if server.StateParam.Logs[i].IsConfigurationChange{
			server.UpdateConfigAndWriteToStorage(server.StateParam.Logs[i].NewConfig)
		}
		server.StateParam.RemainingTickets -= server.StateParam.Logs[i].Num
	}
}

func (server *Server) UpdateConfigAndWriteToStorage(confStr string){
	var newConfig []Peer

	err := json.Unmarshal([]byte(confStr), &newConfig)
	if err != nil{
		log.Println("parse log error")
	}

	oldConfigAddressMap := make(map[string]bool)
	for _, serverFromOldConfig := range server.Conf.Servers{
		oldConfigAddressMap[serverFromOldConfig.Address] = true
	}

	for _, peer := range newConfig{
		if _, ok := oldConfigAddressMap[peer.Address]; !ok {
			newServer := Peer{
				Address: peer.Address,
				ProcessId: peer.ProcessId,
			}
			server.Conf.Servers = append(server.Conf.Servers, &newServer)
		}
	}

	newConfigAddressMap := make(map[string]bool)
	for _, serverFromNewConfig := range newConfig{
		newConfigAddressMap[serverFromNewConfig.Address] = true
	}

	for i, serverFromSuperSet := range server.Conf.Servers {
		if _, ok := newConfigAddressMap[serverFromSuperSet.Address]; !ok{
			server.Conf.Servers = append(server.Conf.Servers[:i],
				server.Conf.Servers[i+1:]...)
		}
	}

	configJson, err := json.MarshalIndent(server.Conf, "", "    ")
	check(err)
	err = ioutil.WriteFile("./server_conf.json", configJson, 0755)
}

func check(err error){
	if err != nil {
		panic(err)
	}
}