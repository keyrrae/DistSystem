package main

import (
	"os"
	"log"
	"errors"
	"time"
	"encoding/json"
)

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
