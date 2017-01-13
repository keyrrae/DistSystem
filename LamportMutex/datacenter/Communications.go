package main

import (
	"log"
	"net/rpc"
	"time"
	"fmt"
	"sync"
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
	//conf.RemainingTickets -= args.BuyTickets
	lamClock.LogicalClock++
	changeRequest := &Request{
		Request: args.BuyTickets,
		Clock:   LamportClock{lamClock.LogicalClock, lamClock.ProcId},
	}

	waitQueue.Push(changeRequest)
	
	// TODO: send and receive from other data centers
	req := DataCenterRequest{
		RequestType: ASK,
		RequestBody: *changeRequest,
	}
	
	lock := &sync.Mutex{}
	count := 0
	
	done := make(chan bool)
	
	go func(){
		for{
			lock.Lock()
			if count == conf.NumOfServers(){
				
				done <- true
			}
			lock.Unlock()
			time.Sleep(100 * time.Millisecond)
		}
	}()
	
	for _, server := range connections{
			
			dataCenterReply := new(DataCenterReply)
			divCall := server.Go("DataCenterComm.CriticalSectionRequest", req, dataCenterReply, nil)
		
		go func(){
			replyCall := <-divCall.Done
			
			gotReply := replyCall.Reply.(*DataCenterReply)
			if gotReply.Grant{
				lock.Lock()
				count++
				lock.Unlock()
			}
			fmt.Println(gotReply.Grant)
		}()
	}
	
	<- done
	conf.RemainingTickets -= changeRequest.Request
	*reply = conf.RemainingTickets
	
	for _, server := range connections{
		
		dataCenterReply := new(DataCenterReply)
		req = DataCenterRequest{
			RequestType: RELEASE,
			RequestBody: *changeRequest,
		}
		releaseCall := server.Go("DataCenterComm.CriticalSectionRequest", req, dataCenterReply, nil)
		_ = releaseCall
		/*
		go func(){
			replyCall := <-divCall.Done
			
			gotReply := replyCall.Reply.(*DataCenterReply)
			if gotReply.Grant{
				lock.Lock()
				count++
				lock.Unlock()
			}
			fmt.Println(gotReply.Grant)
		}()*/
	}
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
	Grant bool
	TimeStamp LamportClock
}

func max(a int64, b int64) int64{
	if a > b {
		return a
	}
	return b
}

func (t *DataCenterComm) CriticalSectionRequest(req *DataCenterRequest, reply *DataCenterReply) error {
	lamClock.LogicalClock = max(req.RequestBody.Clock.LogicalClock, lamClock.LogicalClock) + 1
	// TODO: identify request type
	fmt.Println(req.RequestType)
	
	switch req.RequestType{
	case ASK:
		waitQueue.Push(&(req.RequestBody))
		if req.RequestBody.equalsTo(*(waitQueue.Peek())){
			reply.Grant = true
			reply.TimeStamp = lamClock
			
		} else{
			reply.Grant = false
			reply.TimeStamp = lamClock
		}
	case RELEASE:
		conf.RemainingTickets -= req.RequestBody.Request
	default:
		
	}
	/*
	changeRequest := &Request{
		request: req.RequestBody.request,
		clock:   LamportClock{req.RequestBody.clock.logicalClock, req.RequestBody.clock.procId},
	}
	
	waitQueue.Push(changeRequest)
	
	if changeRequest.equalsTo(*(waitQueue.Peek())){
		reply = &DataCenterReply{
			Grant: true,
			TimeStamp: *(changeRequest.clock),
		}
	} else{
		
	}
	*/
	//TODO: reply
	
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
