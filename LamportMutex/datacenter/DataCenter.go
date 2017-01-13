package main

import (
	"fmt"
	"net/http"
	"net/rpc"
)

var waitQueue PriorityQueue
var conf Config
var lamClock LamportClock

func init() {
	conf = ReadConfig()
	lamClock = LamportClock{0, conf.ProcessID}
}

func main() {

	arith := new(Mutex)
	arith.value = 1000
	rpc.Register(arith)
	rpc.HandleHTTP()
	
	go EstablishConnections()
	go waitUserInput()
	
	err := http.ListenAndServe(conf.Self, nil)
	if err != nil {
		fmt.Println(err.Error())
	}
}
