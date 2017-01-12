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
}

type Server struct {
	Address string `json:"address"`
}

func ReadConfig() Config {
	var conf Config

	file, err := ioutil.ReadFile("./servers.conf")
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
	fmt.Println(conf)
	return conf
}
