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

func (t *Mutex) BuyTicketRequest(args *Args, reply *int) error {
	conf.RemainingTickets -= args.BuyTickets
	lamClock.logicalClock++
	lamClockCopy := &Request{
		request: args.BuyTickets,
		clock: LamportClock{lamClock.logicalClock, lamClock.procId},
	}
	
	waitQueue.Push(lamClockCopy)
	*reply = conf.RemainingTickets
	return nil
}

var waitQueue PriorityQueue
var conf Config
var lamClock LamportClock

func init(){
	conf = ReadConfig()
	lamClock = LamportClock{0, conf.ProcessID}
}

func main() {

	arith := new(Mutex)
	arith.value = 1000

	rpc.Register(arith)
	rpc.HandleHTTP()

	err := http.ListenAndServe(conf.Self, nil)
	if err != nil {
		fmt.Println(err.Error())
	}
}
