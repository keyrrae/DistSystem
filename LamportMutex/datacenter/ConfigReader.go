package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
)

type Config struct {
	Self             string   `json:"self"`
	ProcessID        int      `json:"processid"`
	Servers          []Server `json:"servers"`
	RemainingTickets int      `json:"tickets"`
	MaxAttempts      int      `json:"max_attempts"`
	Delay            int      `json:"delay_in_seconds"`
	InitialTktNum    int
}

type Server struct {
	Address string `json:"address"`
}

func (conf Config) NumOfServers() int {
	return len(conf.Servers)
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
	for i, server := range conf.Servers {
		if server.Address == conf.Self {
			conf.Servers = append(conf.Servers[:i], conf.Servers[i+1:]...)
		}
	}
	conf.InitialTktNum = conf.RemainingTickets
	fmt.Println(conf)
	return conf
}
