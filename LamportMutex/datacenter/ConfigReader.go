package main

import (
	"io/ioutil"
	"log"
	"encoding/json"
	"fmt"
)

type Config struct{
	Self string         `json:"self"`
	Servers []Server    `json:"servers"`
}

type Server struct{
	Address string      `json:"address"`
}


func ReadConfig() *Config{
	var conf Config
	
	file, err := ioutil.ReadFile("./servers.conf")
	if err != nil {
		log.Fatal(err, "\r\n")
	}
	
	err = json.Unmarshal(file, &conf)
	if err != nil {
		log.Fatal(err, "\r\n")
	}
	for i, server := range conf.Servers{
		if server.Address == conf.Self{
			conf.Servers = append(conf.Servers[:i], conf.Servers[i+1:]...)
		}
	}
	
	fmt.Println(conf)
	
	return &conf
}


