package main

import (
	"fmt"
	"net/http"
	"net/rpc"
)

type Args struct {
	BuyTickets int
}

type Mutex struct {
	value int
}

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
	
	EstablishConnections()
	rpc.Register(arith)
	rpc.HandleHTTP()

	err := http.ListenAndServe(conf.Self, nil)
	if err != nil {
		fmt.Println(err.Error())
	}
}