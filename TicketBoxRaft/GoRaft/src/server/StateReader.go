package main

import (
	"io/ioutil"
	"log"
	"fmt"
	"encoding/json"
)

type StateParams struct {
	CurrentTerm int `json:"self"`
	VotedFor    int    `json:"processid"`
	Logs        []int  `json:"logs"`
	CommitIndex int
	LastApplied int
	NextIndex   []int // for leader, reinitialized after election
	MatchIndex  []int // for leader, reinitialized after election
}

func ReadSavedState() StateParams {
	var stateParam StateParams
	
	file, err := ioutil.ReadFile("./saved_state.json")
	if err != nil {
		stateParam.CurrentTerm = 0
		stateParam.VotedFor = -1
		fmt.Println(stateParam)
		return stateParam
	}
	
	err = json.Unmarshal(file, &stateParam)
	if err != nil {
		log.Fatal(err, "\r\n")
	}
	
	fmt.Println(stateParam)
	return stateParam
}