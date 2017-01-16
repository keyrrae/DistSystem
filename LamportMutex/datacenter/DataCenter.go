package main

import (
	"container/heap"
	"fmt"
	"net/http"
	"net/rpc"
	"sync"
)

var waitQueue PriorityQueue
var conf Config
var lamClock LamportClock

func init() {
	conf = ReadConfig()
	heap.Init(&waitQueue)

	lamClock = LamportClock{1, conf.ProcessID}
}

var lock *sync.Mutex

func main() {
	lock = &sync.Mutex{}
	clientComm := new(ClientComm)

	dataCenterComm := new(DataCenterComm)
	clientComm.value = 1000
	rpc.Register(clientComm)
	rpc.Register(dataCenterComm)
	rpc.HandleHTTP()

	go EstablishConnections()
	go waitUserInput()

	err := http.ListenAndServe(conf.Self, nil)
	if err != nil {
		fmt.Println(err.Error())
	}
}
