package main

import (
	"log"
	"net/rpc"
	"time"
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
	changeRequest := &Request{
		request: args.BuyTickets,
		clock:   LamportClock{lamClock.logicalClock, lamClock.procId},
	}

	waitQueue.Push(changeRequest)
	*reply = conf.RemainingTickets
	return nil
}

var connectionStatus map[string]bool

type Connections []*rpc.Client
var connections Connections
var allConnected bool

func EstablishConnections() {
	// Establishing connections to other data centers
	// TODO: test this part
	connectionCounter := 0
	allConnected = false

	connectionStatus = make(map[string]bool)
	log.Println("number of servers:", conf.Servers)
	for _, server := range conf.Servers {
		connectionStatus[server.Address] = false
	}

	var i int
	for i = 0; i < conf.MaxAttempts; i++ {
		log.Printf("%v  ", i)
		for _, server := range conf.Servers {
			log.Println(connectionStatus)
			if connectionStatus[server.Address] {
				continue
			}
			
			client, err := rpc.DialHTTP("tcp", server.Address)
			if err != nil {
				log.Println("dialing:", err.Error()+", retrying...")
			} else {
				connectionStatus[server.Address] = true
				connectionCounter++
				connections = append(connections, client)
				break
			}
		}
		if connectionCounter == conf.NumOfServers() {
			allConnected = true
			break
		}
		time.Sleep(1000 * time.Millisecond)
	}

	if i == conf.MaxAttempts {
		log.Fatal("Maximum attempts: cannot reach all the datacenters")
	}
}
