package main

import (
	"os"
	"time"
	"fmt"
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

	for i := server.StateParam.CommitIndex + 1; i < len(self.StateParam.Logs); i++ {
		if server.StateParam.Logs[i].IsConfigurationChange{
			server.UpdateConfigAndWriteToStorage(server.StateParam.Logs[i].NewConfig)
		}
	}
	server.CheckIfShouldStepDown()
}

func (server *Server) CheckIfShouldStepDown(){
	shouldStay := false

	for _, peer := range server.Conf.Servers {
		if peer.Address == server.Conf.MyAddress {
			shouldStay = true
		}
	}

	if !shouldStay {
		log.Println("Follower not in the new Configuration")
		log.Println("Stepping down....")
		os.Exit(0)
	}
}

func (server *Server) UpdateConfigAndWriteToStorage(confStr string){
	var newConfig []Peer

	err := json.Unmarshal([]byte(confStr), &newConfig)
	if err != nil{
		log.Println("parse log error")
	}

	server.Conf.Servers = self.Conf.Servers[:0]
	for _, peer := range newConfig {
		newServer := Peer{
			Address: peer.Address,
			ProcessId: peer.ProcessId,
		}
		server.Conf.Servers = append(server.Conf.Servers, &newServer)
	}

	server.Conf.Peers = server.Conf.Peers[:0]
	for _, peer := range server.Conf.Servers {
		if peer.Address != server.Conf.MyAddress {
			newPeer := Peer{
				Address: peer.Address,
				ProcessId: peer.ProcessId,
				NextIndex    :0,
				MatchedIndex: -1,
			}
			server.Conf.Peers = append(server.Conf.Peers, &newPeer)
		}
	}

	server.Conf.PeersMap = make(map[int]*Peer)

	for _, peer := range server.Conf.Peers {
		server.Conf.PeersMap[peer.ProcessId] = peer
	}

	log.Println("peers", len(server.Conf.Peers))
	server.Conf.NumMajority = ( len(server.Conf.Peers) + 1) / 2 + 1

	log.Println("majority", server.Conf.NumMajority)

	configJson, err := json.MarshalIndent(server.Conf, "", "    ")
	check(err)
	err = ioutil.WriteFile("./server_conf.json", configJson, 0755)
}

func (server *Server) PrintLogs(){
	if len(server.StateParam.Logs) == 0{
		return
	}

	fmt.Println()
	fmt.Println("Logs:")
	for _, logEntry := range server.StateParam.Logs{
		fmt.Println(logEntry)
	}
	fmt.Println()
}

func check(err error){
	if err != nil {
		panic(err)
	}
}