package main

import (
	"fmt"
	_ "log"
	"net/http"
	"net/rpc"
	"time"
)

var self Server

func init() {
	configuration := ReadConfig()

	self = Server{
		Conf:        configuration,
		State:       FOLLOWER,
		GotNumVotes: 0,
	}

	stateParams := readSavedState()

	self.StateParam = stateParams
	self.LastHeartbeat = time.Now()
}

func main() {

	clientComm := new(ClientComm)
	dataCenterComm := new(DataCenterComm)

	clientComm.value = 1000
	rpc.Register(clientComm)
	rpc.Register(dataCenterComm)
	rpc.HandleHTTP()

	//go EstablishConnections()
	go runStateMachine()
	go waitUserInput()

	err := http.ListenAndServe(self.Conf.MyAddress, nil)
	if err != nil {
		fmt.Println(err.Error())
	}
}
