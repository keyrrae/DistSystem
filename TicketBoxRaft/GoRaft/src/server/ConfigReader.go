package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/rpc"
	"time"
	_ "math/rand"
	"math/rand"
)

type Config struct {
	MyAddress        string        `json:"self"`
	ProcessID        int           `json:"processid"`
	Servers          []*Peer       `json:"servers"`
	RemainingTickets int           `json:"tickets"`
	Timeout          time.Duration `json:"election_timeout"`
	MaxAttempts      int           `json:"max_attempts"`
	Delay            int           `json:"delay_in_seconds"`
	InitialTktNum    int
	NumMajority      int
	Peers			 []*Peer
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

	conf.NumMajority = len(conf.Servers)/2 + 1
	for _, peer := range conf.Servers {
		if peer.Address != conf.MyAddress {
			newPeer := Peer{
				Address: peer.Address,
				ProcessId: peer.ProcessId,
				NextIndex    :0,

				MatchedIndex: -1,
			}
			conf.Peers = append(conf.Peers, &newPeer)
		}
	}

	conf.PeersMap = make(map[int]*Peer)

	for _, peer := range conf.Peers {
		conf.PeersMap[peer.ProcessId] = peer
	}

	rand.Seed( time.Now().UTC().UnixNano())
	fmt.Println("peersmap", conf.PeersMap)
	
	rand := int64(conf.Timeout) + rand.Int63n(int64(conf.Timeout))
	fmt.Println("rand", rand)
	
	conf.Timeout = time.Duration(rand)
	
	fmt.Println("timeout", conf.Timeout)
	conf.InitialTktNum = conf.RemainingTickets
	fmt.Println(conf)
	return conf
}
