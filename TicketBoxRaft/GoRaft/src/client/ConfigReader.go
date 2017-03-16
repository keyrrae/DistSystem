package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"strconv"
)

type Server struct {
	Address     string `json:"server_address"`
	MaxAttempts int    `json:"max_attempts"`
	NewConfig []Peer `json:"new_configuration"`
	AllServers []Peer `json:"all_servers"`
}

type Peer struct {
	Address   string `json:"address"`
	ProcessId int    `json:"id"`
}

var ConfigMap map[string]Peer

func ReadConfig() Server {
	var server Server

	file, err := ioutil.ReadFile("./client_conf.json")
	if err != nil {
		log.Fatal(err, "\r\n")
	}

	err = json.Unmarshal(file, &server)
	if err != nil {
		log.Fatal(err, "\r\n")
	}
	log.Println(server)

	ConfigMap = make(map[string]Peer)
	for i, peer := range server.AllServers{
		s := strconv.Itoa(i+1)
		ConfigMap["dc"+s] = peer
	}

	log.Println(ConfigMap)
	return server
}
