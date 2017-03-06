package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/rpc"
	"time"
)

type Config struct {
	MyAddress        string        `json:"self"`
	ProcessID        int           `json:"processid"`
	Peers            []Peer        `json:"servers"`
	RemainingTickets int           `json:"tickets"`
	Timeout          time.Duration `json:"election_timeout"`
	MaxAttempts      int           `json:"max_attempts"`
	Delay            int           `json:"delay_in_seconds"`
	InitialTktNum    int
	NumMajority      int
	PeersMap         map[int]*Peer
}

type Peer struct {
	Address      string `json:"address"`
	ProcessId    int    `json:"id"`
	NextIndex    int
	MatchedIndex int
	Connected    bool
	Comm         *rpc.Client
}

func (conf Config) NumOfServers() int {
	return len(conf.Peers)
}

func ReadConfig() Config {
	var conf Config

	file, err := ioutil.ReadFile("./server_conf.json")
	if err != nil {
		log.Fatal(err, "\r\n")
	}

	err = json.Unmarshal(file, &conf)
	if err != nil {
		log.Fatal(err, "\r\n")
	}

	conf.NumMajority = len(conf.Peers)/2 + 1
	for i, peer := range conf.Peers {
		if peer.Address == conf.MyAddress {
			conf.Peers = append(conf.Peers[:i], conf.Peers[i+1:]...)
			break
		}
	}

	self.Conf.PeersMap = make(map[int]*Peer)

	for _, peer := range conf.Peers {
		peer.NextIndex = 1
		self.Conf.PeersMap[peer.ProcessId] = &peer
	}

	conf.Timeout = conf.Timeout * time.Second
	conf.InitialTktNum = conf.RemainingTickets
	fmt.Println(conf)
	return conf
}
