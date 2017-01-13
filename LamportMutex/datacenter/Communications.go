package main

import (
	"net/rpc"
	"log"
	"time"
)

func (t *Mutex) BuyTicketRequest(args *Args, reply *int) error {
	conf.RemainingTickets -= args.BuyTickets
	lamClock.logicalClock++
	changeRequest := &Request{
		request: args.BuyTickets,
		clock:   LamportClock{lamClock.logicalClock, lamClock.procId},
	}
	
	waitQueue.Push(changeRequest)
	*reply = conf.RemainingTickets
	return nil
}

var connectionStatus map[*Server]bool
type Connections []*rpc.Client

var connections Connections

func EstablishConnections(){
	// Establishing connections to other data centers
	// TODO: test this part
	connectionCounter :=0
	
	connectionStatus = make(map[*Server]bool)
	for _, server := range conf.Servers{
		connectionStatus[&server] = false
	}
	
	var i int
	for i=0; i < conf.MaxAttempts; i++{
			
		for _, server := range conf.Servers{
			if connectionStatus[&server]{
				continue
			}
			
			client, err := rpc.DialHTTP("tcp", server.Address)
			if err != nil {
				log.Println("dialing:", err.Error()+", retrying...")
			} else {
				connectionStatus[&server] = true
				connectionCounter++
				connections = append(connections, client)
				break
			}
		}
		if connectionCounter == conf.NumOfServers(){
			break
		}
		time.Sleep(1000 * time.Millisecond)
	}
}