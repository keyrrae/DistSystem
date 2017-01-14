package main

import (
	"log"
	"net/rpc"
	"time"
	"fmt"
)

type Args struct {
	BuyTickets int
}

type ClientComm struct {
	value int
}

const (
	ASK = "ASK"
	RELEASE = "RELEASE"
)

func (t *ClientComm) BuyTicketRequest(args *Args, reply *int) error {
	lamClock.LogicalClock++
	
	// piggybacking time stamp with request
	changeRequest := &Request{
		Request: args.BuyTickets,
		Clock:   LamportClock{lamClock.LogicalClock, lamClock.ProcId},
	}
	
	// places the request on waitQueue
	waitQueue.Push(changeRequest)
	
	// TODO: send and receive from other data centers
	// sends request message to all sites
	
	// asking for permission to enter critical section
	req := DataCenterRequest{
		RequestType: ASK,
		RequestBody: *changeRequest,
	}
	
	count := 0
	done := make(chan bool)
		
	go func(){
		for{
			// has received messages with larger timestamps from all other sites
			allOtherSitesAgree := count == conf.NumOfServers()
			
			// request is at the top of waitQueue
			requestAtTop := changeRequest.Clock.equalsTo(waitQueue.Peek().Clock)
			
			if allOtherSitesAgree && requestAtTop {
				done <- true
				break
			}
		
			time.Sleep(1000 * time.Millisecond)
		}
	}()

	for _, server := range connections{
			
		dataCenterReply := new(DataCenterReply)
		divCall := server.Go("DataCenterComm.CriticalSectionRequest", req, dataCenterReply, nil)
		
		go func(){
			replyCall := <-divCall.Done
			gotReply := replyCall.Reply.(*DataCenterReply)
			
			lamClock.LogicalClock = max(gotReply.TimeStamp.LogicalClock, lamClock.LogicalClock) + 1
			// increase the counter if got larger time stamp
			if gotReply.TimeStamp.largerThan(changeRequest.Clock){
				count++
			}

		}()
	}
	
	<- done  // block the main thread
	
	var releaseDecAmount int
	if conf.RemainingTickets < changeRequest.Request {
		releaseDecAmount = 0
	} else {
		releaseDecAmount = changeRequest.Request
	}
	
	waitQueue.Pop()
	releaseRequest := &Request{
		Request: releaseDecAmount,
		Clock:   LamportClock{lamClock.LogicalClock, lamClock.ProcId},
	}
	
	conf.RemainingTickets -= releaseDecAmount
	*reply = conf.RemainingTickets
	count = 0
	
	go func() {
		for {
			// has received messages with larger timestamps from all other sites
			allOtherSitesAgree := count == conf.NumOfServers()
						
			if allOtherSitesAgree {
				done <- true
				break
			}
			
			time.Sleep(100 * time.Millisecond)
		}
	}()
	
	for _, server := range connections{
		
		dataCenterReply := new(DataCenterReply)
		req = DataCenterRequest{
			RequestType: RELEASE,
			RequestBody: *releaseRequest,
		}
		releaseCall := server.Go("DataCenterComm.CriticalSectionRequest", req, dataCenterReply, nil)
		
		go func(){
			replyCall := <-releaseCall.Done
			
			gotReply := replyCall.Reply.(*DataCenterReply)
		
			lamClock.LogicalClock = max(gotReply.TimeStamp.LogicalClock, lamClock.LogicalClock) + 1
			count++
		}()
	}
	<- done
	return nil
}

type DataCenterComm struct {
	value int
}

type DataCenterRequest struct{
	RequestType string
	RequestBody Request
}

type DataCenterReply struct {
	TimeStamp LamportClock
}

func max(a int64, b int64) int64{
	if a > b {
		return a
	}
	return b
}

func (t *DataCenterComm) CriticalSectionRequest(req *DataCenterRequest, reply *DataCenterReply) error {
	// upon receives a request, increase the logic clock
	lamClock.LogicalClock = max(req.RequestBody.Clock.LogicalClock, lamClock.LogicalClock) + 1
	
	fmt.Println()
	fmt.Print(req.RequestType)
	
	switch req.RequestType{
	case ASK:
		// receives a request asking to enter critical section
		// reply with timestamp of this site
		reply.TimeStamp = lamClock
		
		// push request to waitQueue
		waitQueue.Push(&(req.RequestBody))
	case RELEASE:
		waitQueue.Pop()
		conf.RemainingTickets -= req.RequestBody.Request
		fmt.Println()
		fmt.Print("> ")
	default:
		
	}
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
			fmt.Println("Number of clients:", conf.NumOfServers())
			//fmt.Println(connections)
			allConnected = true
			break
		}
		time.Sleep(1000 * time.Millisecond)
	}

	if i == conf.MaxAttempts {
		log.Fatal("Maximum attempts: cannot reach all the datacenters")
	}
}
