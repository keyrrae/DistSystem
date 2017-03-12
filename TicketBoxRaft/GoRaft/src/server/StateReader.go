package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
)

type StateParameters struct {
	CurrentTerm int        `json:"self"`
	VotedFor    int        `json:"processid"`
	Logs        []LogEntry `json:"logs"`
	CommitIndex int
	LastApplied int
}

type LogEntry struct {
	Num                   int `json:"value"`
	Term                  int `json:"term"`
	IsConfigurationChange bool
	NewConfig             string
}

func readSavedState() StateParameters {
	var stateParam StateParameters

	file, err := ioutil.ReadFile("./saved_state.json")
	if err != nil {
		stateParam.CurrentTerm = 0
		stateParam.VotedFor = -1
		stateParam.CommitIndex = -1
		stateParam.LastApplied = -1
		fmt.Println(stateParam)
		return stateParam
	}

	err = json.Unmarshal(file, &stateParam)
	if err != nil {
		log.Fatal(err, "\r\n")
	}

	stateParam.LastApplied = -1

	fmt.Println(stateParam)
	return stateParam
}
