package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

type Server struct {
	Address     string `json:"address"`
	MaxAttempts int    `json:"max_attempts"`
}

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

	return server
}
