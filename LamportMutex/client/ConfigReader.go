package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

type Server struct {
	Address string `json:"address"`
}

func ReadConfig() string {
	var server Server

	file, err := ioutil.ReadFile("./client.conf")
	if err != nil {
		log.Fatal(err, "\r\n")
	}

	err = json.Unmarshal(file, &server)
	if err != nil {
		log.Fatal(err, "\r\n")
	}

	return server.Address
}
